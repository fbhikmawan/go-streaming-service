package main

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/unbot2313/go-streaming-service/config"
	_ "github.com/unbot2313/go-streaming-service/docs"
	"github.com/unbot2313/go-streaming-service/internal/app"
	"github.com/unbot2313/go-streaming-service/internal/routes"
)

// @title Go Streaming Service API
// @version 1.0
// @description A streaming service API using Go and Gin framework, with Swagger documentation and ffmpeg integration.

// @host	localhost:3003
// @BasePath /api/v1

func main() {

	r := gin.Default()

    // Enable debug mode
    gin.SetMode(gin.DebugMode)

    // Configure CORS middleware
    corsConfig := cors.Config{
        AllowOrigins:     []string{"*"},
        AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: false,
        MaxAge:           12 * time.Hour,
    }
    r.Use(cors.New(corsConfig))

	apiGroup := r.Group("/api")

	v1Group := apiGroup.Group("/v1")

	// conect to database
	_, err := config.GetDB()
	if err != nil {
		panic(err)
	}

	// Serving static files (STREAMING)
	// when accessing the /static path, files in the public folder are served,
	// e.g.: http://localhost:3003/static/index.html, /public/index.html is served.
	v1Group.Static("/static", "./static/temp")

	// Initialize the application components
	userController, authController, videoController := app.InitializeComponents()

	// Configure the routes
	routes.SetupRoutes(v1Group, userController, authController, videoController)
	// Configuring Swagger documentation
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))


	r.Run(":3003")

}