package framework

type ServiceInstance struct {
}

func (inst *ServiceInstance) request() string {}

func load() (instance ServiceInstance, err error) {
	loadConfig()
	instance = ServiceInstance{}

	return
}
