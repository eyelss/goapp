package registry

import "context"

type ServiceID = string

type ServiceInstancee struct {
	ID      ServiceID
	Name    string
	Address string
	Meta    map[string]string
}

type IRegistry interface {
	Register(ctx context.Context, instance ServiceInstancee) error
	Unregister(ctx context.Context, instanceID ServiceID) error
	Close() error
}
