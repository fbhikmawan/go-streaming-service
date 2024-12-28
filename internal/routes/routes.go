package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/unbot2313/go-streaming-service/internal/controllers"
	"github.com/unbot2313/go-streaming-service/internal/middlewares"
)

// SetupRoutes configures all routes
func SetupRoutes(router *gin.RouterGroup, userController controllers.UserController, authController controllers.AuthController, videoController controllers.VideoController) {
	// User routes
	userRoutes := router.Group("/users")
	{
		userRoutes.GET("/id/:id", userController.GetUserByID)
		userRoutes.GET("/username/:username", userController.GetUserByUserName)
		userRoutes.POST("/", userController.CreateUser)
		userRoutes.DELETE("/:id", userController.DeleteUserByID)
	}

	// Authentication paths
	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/login", authController.Login)
		authRoutes.POST("/register", authController.Register)
	}

    VideoRoutes := router.Group("/streaming")
    {
		ProtectedRoute := VideoRoutes.Group("")
		ProtectedRoute.Use(middlewares.AuthMiddleware)

		// Public roads
        VideoRoutes.GET("/latest", videoController.GetLatestVideos)
		VideoRoutes.GET("/id/:videoid", videoController.GetVideoByID)
		VideoRoutes.PATCH("/views/:videoid", videoController.IncrementViews)

		// Protected route
        ProtectedRoute.POST("/upload", videoController.CreateVideo)
    }
	
}
