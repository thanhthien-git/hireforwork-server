package service

import (
	"fmt"
)

// ServiceContainer chứa tất cả các service dưới dạng map
type ServiceContainer struct {
	services map[string]interface{}
}

// NewServiceContainer tạo một instance của ServiceContainer
func NewServiceContainer() *ServiceContainer {
	return &ServiceContainer{
		services: make(map[string]interface{}),
	}
}

// Register đăng ký một service vào container
func (c *ServiceContainer) Register(serviceName string, service interface{}) {
	if service == nil {
		fmt.Printf("Warning: Attempted to register nil service for %s\n", serviceName)
		return
	}
	c.services[serviceName] = service
}

// Get lấy service từ container
func (c *ServiceContainer) Get(serviceName string) interface{} {
	service, exists := c.services[serviceName]
	if !exists {
		fmt.Printf("Service %s not found in container\n", serviceName)
		return nil
	}
	return service
}

// AppServices là container cuối cùng để sử dụng trong ứng dụng
type AppServices struct {
	container *ServiceContainer
}

// NewAppServices tạo một instance của AppServices
func NewAppServices(container *ServiceContainer) *AppServices {
	return &AppServices{container: container}
}

// GetService lấy service từ AppServices
func (a *AppServices) GetService(serviceName string) interface{} {
	return a.container.Get(serviceName)
}
