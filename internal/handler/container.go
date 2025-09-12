package handler

import (
	"gin-mysql-api/internal/service"
)

// Container 处理器容器
type Container struct {
	HealthHandler *HealthHandler
	AuthHandler   *AuthHandler
	UserHandler   *UserHandler
	DramaHandler  *DramaHandler
	AdminHandler  *AdminHandler
	FileHandler   *FileHandler
}

// NewContainer 创建处理器容器
func NewContainer(services *service.Container) *Container {
	return &Container{
		HealthHandler: NewHealthHandler(),
		AuthHandler:   NewAuthHandler(services.AuthService),
		UserHandler:   NewUserHandler(services.UserService),
		DramaHandler:  NewDramaHandler(services.DramaService),
		AdminHandler:  NewAdminHandler(services.AdminService, services.UserService),
		FileHandler:   NewFileHandler(services.FileService),
	}
}