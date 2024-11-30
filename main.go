package main

import (
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	docs "github.com/unbot2313/go-streaming-service/docs"
	"github.com/unbot2313/go-streaming-service/internal/controllers"
	"github.com/unbot2313/go-streaming-service/internal/routes"
	"github.com/unbot2313/go-streaming-service/internal/services"
)

// @title Go Streaming Service API
// @version 1.0
// @description A streaming service API using Go and Gin framework, with Swagger documentation and ffmpeg integration.

// @host	localhost:3003
// @BasePath /api/v1

func main() {

	r := gin.Default()

	docs.SwaggerInfo.BasePath = "/api/v1"

	apiGroup := r.Group("/api")

	v1Group := apiGroup.Group("/v1")

	service := services.NewUserService()

	controller := controllers.NewUserController(service)

	routes.SetupRoutes(v1Group, controller)

	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))


	r.Run(":3003")

}