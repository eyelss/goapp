package server

import (
	"goapp/framework/registry"
	"time"

	"google.golang.org/grpc"
)

type Options struct {
	ServiceName        string
	ServiceAddr        string
	Registry           registry.IRegistry
	RegistryInterval   time.Duration
	GracefulTimeout    time.Duration
	ServiceOptions     []grpc.ServerOption
	UnaryInterceptors  []grpc.UnaryServerInterceptor
	StreamInterceptors []grpc.StreamServerInterceptor
}

type Option func(*Options)

func WithName(name string) Option {
	return func(o *Options) { o.ServiceName = name }
}

func WithAddr(addr string) Option {
	return func(o *Options) { o.ServiceAddr = addr }
}

func WithRegistry(reg registry.IRegistry) Option {
	return func(o *Options) { o.Registry = reg }
}

func WithTimeout(timeout time.Duration) Option {
	return func(o *Options) { o.GracefulTimeout = timeout }
}
