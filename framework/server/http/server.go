package http

import (
	"context"
	"fmt"
	config "goapp/framework/lib"
	"goapp/framework/registry"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/labstack/echo/v5"
)

type Server struct {
	http       *http.Server
	Echo       *echo.Echo
	opts       Options
	instanceID registry.ServiceID
	mutex      sync.Mutex
	running    bool
	regCancel  context.CancelFunc
}

func New(opts ...Option) (*Server, error) {
	config.Load()

	e := echo.New()

	o := &Options{
		ServiceName: config.Get[string]("app-name"),
		ServiceAddr: config.Get[string]("addr"),
	}

	for _, opt := range opts {
		opt(o)
	}

	return &Server{
		Echo: e,
		opts: *o,
	}, nil
}

func (s *Server) Start() error {
	s.mutex.Lock()

	if s.running {
		s.mutex.Unlock()

		return fmt.Errorf("server already running")
	}

	s.running = true

	s.mutex.Unlock()

	httpSrv := http.Server{Addr: s.opts.ServiceAddr, Handler: s.Echo}

	go func() {
		if err := httpSrv.ListenAndServe(); err != nil {
			fmt.Printf("could not start http server: %v\n", err)
		}
	}()

	log.Printf("server listening at %v", httpSrv.Addr)

	s.http = &httpSrv

	go s.waitForShutdown()

	return nil
}

func (s *Server) Stop() {
	s.mutex.Lock()

	if !s.running {
		s.mutex.Unlock()
		return
	}

	s.running = false
	s.mutex.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.http.Shutdown(ctx); err != nil {
		log.Printf("could not shutdown http server: %v\n", err)
	}

	os.Exit(0)
}

func (s *Server) waitForShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	log.Println("shutting down server...")

	s.Stop()
}
