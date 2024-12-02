package app

import (
	"github.com/unbot2313/go-streaming-service/internal/controllers"
	"github.com/unbot2313/go-streaming-service/internal/services"
)

// InitializeComponents crea las instancias de los servicios y controladores
func InitializeComponents() (controllers.UserController, controllers.AuthController) {
	// Inicializa los servicios
	userService := services.NewUserService()
	authService := services.NewAuthService()

	// Inicializa los controladores
	userController := controllers.NewUserController(userService)
	authController := controllers.NewAuthController(authService)

	return userController, authController
}
