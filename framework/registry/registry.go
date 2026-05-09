package registry

import "context"

type ServiceID = string

type ServiceInstance struct {
	ID      ServiceID
	Name    string
	Address string
	Meta    map[string]string
}

type IRegistry interface {
	Register(ctx context.Context, instance ServiceInstance) error
	Unregister(ctx context.Context, instanceID ServiceID) error
	Close() error
}
