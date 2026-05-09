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
	RequestTimeout     time.Duration
	ServiceOptions     []grpc.ServerOption
	UnaryInterceptors  []grpc.UnaryServerInterceptor
	StreamInterceptors []grpc.StreamServerInterceptor
}

type Option func(*Options)

func WithRegistry(reg registry.IRegistry) Option {
	return func(o *Options) { o.Registry = reg }
}

func WithTimeout(timeout time.Duration) Option {
	return func(o *Options) { o.RegistryInterval = timeout }
}
