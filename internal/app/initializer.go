package app

import (
	"github.com/unbot2313/go-streaming-service/internal/controllers"
	"github.com/unbot2313/go-streaming-service/internal/services"
)

// InitializeComponents crea las instancias de los servicios y controladores
func InitializeComponents() (controllers.UserController, controllers.AuthController, controllers.VideoController, controllers.S3Controller) {
	// Inicializa los servicios
	userService := services.NewUserService()
	authService := services.NewAuthService()

	// Inicializa los controladores
	userController := controllers.NewUserController(userService)
	authController := controllers.NewAuthController(authService)

	// Inicializa el controlador de videos
	videoService := services.NewVideoService()
	videoController := controllers.NewVideoController(videoService)

	// Inicializa el controlador de S3
	s3Configuration := services.GetS3Configuration()
	s3Service := services.NewS3Service(s3Configuration)
	s3Controller := controllers.NewS3Controller(s3Service)


	return userController, authController, videoController, s3Controller
}
