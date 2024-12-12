package app

import (
	"github.com/unbot2313/go-streaming-service/internal/controllers"
	"github.com/unbot2313/go-streaming-service/internal/services"
)

// InitializeComponents crea las instancias de los servicios y controladores
func InitializeComponents() (controllers.UserController, controllers.AuthController, controllers.VideoController) {
	// Inicializa los servicios
	userService := services.NewUserService()
	authService := services.NewAuthService()

	// Inicializa los controladores
	userController := controllers.NewUserController(userService)
	authController := controllers.NewAuthController(authService)

	// Inicializa el controlador de videos
	S3configuration := services.GetS3Configuration()
	filesService := services.NewFilesService()
	videoService := services.NewVideoService(S3configuration, filesService)
	videoController := controllers.NewVideoController(videoService)


	return userController, authController, videoController
}
