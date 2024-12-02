package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/unbot2313/go-streaming-service/internal/controllers"
)

// SetupRoutes configura todas las rutas
func SetupRoutes(router *gin.RouterGroup, userController controllers.UserController, authController controllers.AuthController) {
	// Rutas de usuarios
	userRoutes := router.Group("/users")
	{
		userRoutes.GET("/:id", userController.GetUserByID)
		userRoutes.POST("/", userController.CreateUser)
		userRoutes.DELETE("/:id", userController.DeleteUserByID)
	}

	// Rutas de autenticación
	authRoutes := router.Group("/auth")
	{
		authRoutes.POST("/login", authController.Login)
		authRoutes.POST("/register", authController.Register)
	}
}
