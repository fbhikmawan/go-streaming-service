package app

import (
	"github.com/unbot2313/go-streaming-service/internal/controllers"
	"github.com/unbot2313/go-streaming-service/internal/services"
)

// InitializeComponents creates the instances of the services and drivers
func InitializeComponents() (controllers.UserController, controllers.AuthController, controllers.VideoController) {
	// Initialize services
	userService := services.NewUserService()
	authService := services.NewAuthService()

	// Initializes the controllers
	userController := controllers.NewUserController(userService)
	authController := controllers.NewAuthController(authService, userService)

	// Initializes the video controller
	S3configuration := services.GetS3Configuration()
	filesService := services.NewFilesService()
	videoService := services.NewVideoService(S3configuration, filesService)
	databaseVideoService := services.NewDatabaseVideoService()
	videoController := controllers.NewVideoController(videoService, databaseVideoService)


	return userController, authController, videoController
}
