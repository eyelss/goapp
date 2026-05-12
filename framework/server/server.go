package server

import (
	"context"
	"fmt"
	config "goapp/framework/lib"
	"goapp/framework/registry"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

type Server struct {
	opts       Options
	grpc       *grpc.Server
	health     *health.Server
	instanceID registry.ServiceID
	mutex      sync.Mutex
	running    bool
	regCancel  context.CancelFunc
}

func GetServiceID() registry.ServiceID {
	host, _ := os.Hostname()
	pid := os.Getpid()
	return fmt.Sprintf("%s-%d-%d", host, pid, time.Now().UnixNano())
}

func New(opts ...Option) (*Server, error) {
	config.Load()

	addr := config.Get[string]("addr")
	name := config.Get[string]("app-name")

	o := &Options{
		RegistryInterval: 60 * time.Second,
		GracefulTimeout:  15 * time.Second,
		Registry:         registry.NewStubRegistry(),
		ServiceAddr:      addr,
		ServiceName:      name,
	}

	for _, opt := range opts {
		opt(o)
	}

	serviceID := GetServiceID()

	serverOpts := append(o.ServiceOptions,
		grpc.ChainUnaryInterceptor(o.UnaryInterceptors...),
		grpc.ChainStreamInterceptor(o.StreamInterceptors...),
	)

	grpcServer := grpc.NewServer(serverOpts...)

	healthServer := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthServer)

	server := &Server{
		instanceID: serviceID,
		opts:       *o,
		grpc:       grpcServer,
		health:     healthServer,
	}

	return server, nil
}

func (s *Server) RegisterService(sd *grpc.ServiceDesc, impl interface{}) {
	s.grpc.RegisterService(sd, impl)
}

func getIP() (string, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return "", err
	}
	ips, err := net.LookupHost(hostname)
	if err != nil {
		return "", err
	}
	if len(ips) == 0 {
		return "", fmt.Errorf("no IP found for hostname %s", hostname)
	}
	for _, ip := range ips {
		if net.ParseIP(ip).To4() != nil {
			return ip, nil
		}
	}
	return "", fmt.Errorf("IPv4 not found")
}

func (s *Server) Start(ctx context.Context) error {
	s.mutex.Lock()

	if s.running {
		s.mutex.Unlock()

		return fmt.Errorf("server already running")
	}

	s.running = true
	s.mutex.Unlock()

	listener, err := net.Listen("tcp", s.opts.ServiceAddr)
	if err != nil {
		return fmt.Errorf("could not listen on socket: %w", err)
	}

	ip, _ := getIP()

	if s.opts.Registry != nil {
		instance := registry.ServiceInstance{
			ID:      s.instanceID,
			Name:    s.opts.ServiceName,
			Address: ip,
			Meta: map[string]string{
				"started_at": time.Now().UTC().String(),
			},
		}

		if err := s.opts.Registry.Register(ctx, instance); err != nil {
			listener.Close()
			fmt.Errorf("could not register service: %w", err)
		}

		// periodic heartbeat
		//var registerContext context.Context
		//registerContext, s.regCancel = context.WithCancel(context.Background())
		//go s.syncRegisterProcess(registerContext, instance)
	}

	s.health.SetServingStatus(s.opts.ServiceName, healthpb.HealthCheckResponse_SERVING)
	s.health.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	go func() {
		if err := s.grpc.Serve(listener); err != nil {
			fmt.Printf("could not start grpc server: %v\n", err)
		}
	}()

	log.Printf("server listening at %v", listener.Addr())

	//go s.waitForShutdown()

	return nil
}

func (s *Server) waitForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	s.Stop()
}

func (s *Server) Stop() {
	s.mutex.Lock()

	if !s.running {
		s.mutex.Unlock()
		return
	}

	s.running = false
	s.mutex.Unlock()

	s.health.SetServingStatus(s.opts.ServiceName, healthpb.HealthCheckResponse_NOT_SERVING)
	s.health.SetServingStatus("", healthpb.HealthCheckResponse_NOT_SERVING)

	if s.opts.Registry != nil && s.regCancel != nil {
		s.regCancel()

		ctx, cancel := context.WithTimeout(context.Background(), s.opts.GracefulTimeout)

		defer cancel()

		log.Printf("Stopping service: %s", s.opts.ServiceName)
		_ = s.opts.Registry.Unregister(ctx, s.instanceID)
		_ = s.opts.Registry.Close()

		stopChannel := make(chan struct{})
		go func() {
			s.grpc.GracefulStop()
			close(stopChannel)
		}()

		select {
		case <-stopChannel:
		// expected done
		case <-time.After(s.opts.GracefulTimeout):
			s.grpc.Stop() // force
		}

		os.Exit(0)
	}
}

func (s *Server) syncRegisterProcess(ctx context.Context, instance registry.ServiceInstance) {
	ticker := time.NewTicker(s.opts.RegistryInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			_ = s.opts.Registry.Register(ctx, instance)
		}
	}
}
