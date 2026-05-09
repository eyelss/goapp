package server

import (
	"context"
	"fmt"
	"goapp/framework/registry"
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
	o := &Options{
		RegistryInterval: 60 * time.Second,
		RequestTimeout:   15 * time.Second,
	}

	for _, opt := range opts {
		opt(o)
	}

	serviceID := GetServiceID()

	if o.ServiceName == "" {
		return nil, fmt.Errorf("service name is required")
	}

	if o.ServiceAddr == "" {
		return nil, fmt.Errorf("service address is required")
	}

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

	if s.opts.Registry != nil {
		instance := registry.ServiceInstancee{
			ID:      s.instanceID,
			Name:    s.opts.ServiceName,
			Address: s.opts.ServiceAddr,
			Meta: map[string]string{
				"started_at": time.Now().UTC().String(),
			},
		}

		if err := s.opts.Registry.Register(ctx, instance); err != nil {
			listener.Close()
			return fmt.Errorf("could not register service: %w", err)
		}

		// periodic heartbeat
		var registerContext context.Context
		registerContext, s.regCancel = context.WithCancel(context.Background())
		go s.syncRegisterProcess(registerContext, instance)
	}

	s.health.SetServingStatus(s.opts.ServiceName, healthpb.HealthCheckResponse_SERVING)
	s.health.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	go func() {
		if err := s.grpc.Serve(listener); err != nil {

		}
	}()

	return nil
}

func (s *Server) waitForShutdown(ctx context.Context) error {
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

		ctx, cancel := context.WithTimeout(context.Background(), s.opts.RequestTimeout)

		defer cancel()

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
		case <-time.After(s.opts.RequestTimeout):
			s.grpc.Stop() // force
		}
	}
}

func (s *Server) syncRegisterProcess(ctx context.Context, instance registry.ServiceInstancee) {
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
