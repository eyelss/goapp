package registry

import (
	"context"
	"log"
)

type StubRegistry struct{}

func (r *StubRegistry) Register(ctx context.Context, instance ServiceInstance) error {
	log.Printf("register service instance: %#v", instance)
	return nil
}

func (r *StubRegistry) Unregister(ctx context.Context, instanceID ServiceID) error {
	log.Printf("unregister service instance: %#v", instanceID)
	return nil
}

func (r *StubRegistry) Close() error {
	return nil
}

func NewStubRegistry() *StubRegistry {
	return &StubRegistry{}
}
