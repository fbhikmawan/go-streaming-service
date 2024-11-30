package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/unbot2313/go-streaming-service/internal/controllers"
	"github.com/unbot2313/go-streaming-service/internal/middlewares"
)


func SetupRoutes(router *gin.RouterGroup, userController *controllers.UserControllerImp) {
    userRoutes := router.Group("/users")
    {
        userRoutes.GET("/:id", userController.GetUserByID)
        userRoutes.POST("/", userController.CreateUser)
        userRoutes.DELETE("/:id", userController.DeleteUserByID)
    }

    authRoutes := router.Group("/auth")
    {
        authRoutes.POST("/login", userController.DeleteUserByID)
        authRoutes.POST("/register", userController.CreateUser)
    }

    videoRoutes := router.Group("/videos")
    videoRoutes.Use(middlewares.AuthMiddleware)
    {
        videoRoutes.GET("/:id", userController.GetUserByID)
        videoRoutes.POST("/", userController.CreateUser)
        videoRoutes.DELETE("/:id", userController.DeleteUserByID)
    }
}
