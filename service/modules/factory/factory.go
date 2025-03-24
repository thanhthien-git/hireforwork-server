package service

import (
	"fmt"
	"hireforwork-server/db"
	"hireforwork-server/service"
	modules "hireforwork-server/service/modules"
	auth "hireforwork-server/service/modules/auth"
	job "hireforwork-server/service/modules/jobs"
	observe "hireforwork-server/service/observe"
)

/*
1. Factory pattern
2. ServiceFactory is a factory for creating services
3. ServiceDependencies holds all dependencies that services might need
4. ServiceCreator is a function type that creates a service with dependencies
5. serviceCreators maps service types to their creation functions
*/

type ServiceDependencies struct {
	DB *db.DB
}

type ServiceFactory struct {
	deps *ServiceDependencies
}

// Factory method pattern
// Định nghĩa interface cho factory method
type ServiceCreator func(*ServiceDependencies) interface{}

// Map các hàm tạo đối tượng
/*
1. Sử dụng interface ServiceCreator để tạo đối tượng
2. Mỗi loại service có 1 hàm tạo riêng
3. Cho phép mở rộng = cách thêm hàm tạo mới
4. Tuân thủ nguyên tắt Open/Closed
*/
var serviceCreators = map[string]ServiceCreator{
	"auth":     func(deps *ServiceDependencies) interface{} { return auth.NewAuthService(deps.DB) },
	"job":      func(deps *ServiceDependencies) interface{} { return job.NewJobService(deps.DB) },
	"company":  func(deps *ServiceDependencies) interface{} { return modules.NewCompanyService(deps.DB) },
	"career":   func(deps *ServiceDependencies) interface{} { return modules.NewUserService(deps.DB) },
	"tech":     func(deps *ServiceDependencies) interface{} { return modules.NewTechService(deps.DB) },
	"category": func(deps *ServiceDependencies) interface{} { return modules.NewCategoryService(deps.DB) },
	"field":    func(deps *ServiceDependencies) interface{} { return modules.NewFieldService(deps.DB) },
	"observe": func(deps *ServiceDependencies) interface{} {
		return observe.NewJobEventManager()
	},
}

func NewServiceFactory(deps *ServiceDependencies) *ServiceFactory {
	return &ServiceFactory{deps: deps}
}

//Simple Factory Pattern
/*
1. Chỉ có 1 phương thức duy nhất
2. Dùng 1 string parameter để tạo đối tượng
3. Tạo đối tượng dựa trên kiểu đối tượng cần tạo
4. Không có đối tượng/interface phức tạp
*/
func (f *ServiceFactory) CreateService(serviceType string) interface{} {
	if creator, exists := serviceCreators[serviceType]; exists {
		return creator(f.deps)
	}
	return nil
}

func (f *ServiceFactory) RegisterAllServices(container *service.ServiceContainer) {
	// Only log once at the start
	fmt.Println("Initializing service container...")

	// Register all services
	for serviceType := range serviceCreators {
		container.Register(serviceType, f.CreateService(serviceType))
	}

	fmt.Println("Service container initialized successfully")
}
