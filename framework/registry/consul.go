package registry

import (
	"context"

	"github.com/hashicorp/consul/api"
)

type ConsulRegistry struct {
	client *api.Client
}

func NewConsulRegistry(consulAddr string) (*ConsulRegistry, error) {
	config := api.DefaultConfig()
	config.Address = consulAddr
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &ConsulRegistry{client: client}, nil
}

func (c *ConsulRegistry) Register(ctx context.Context, instance ServiceInstance) error {
	registration := &api.AgentServiceRegistration{
		ID:      instance.ID,
		Name:    instance.Name,
		Address: instance.Address,
		Port:    parsePort(instance.Address),
		Meta:    instance.Meta,
		Check: &api.AgentServiceCheck{
			GRPC:     instance.Address,
			Interval: "10s",
			Timeout:  "5s",
		},
	}

	return c.client.Agent().ServiceRegister(registration)
}

func (c *ConsulRegistry) Unregister(ctx context.Context, instanceID ServiceID) error {
	return c.client.Agent().ServiceDeregister(instanceID)
}

func (c *ConsulRegistry) Close() error { return nil }

func parsePort(address string) int {
	return 50051
}
