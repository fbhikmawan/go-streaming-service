package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	r.Use(cors.Default()) // Habilita CORS para todos los orígenes

	apiGroup := r.Group("/api")

	v1Group := apiGroup.Group("/v1")

	// Servir archivos staticos (STREAMING)
	// cuando se accede a la ruta /static, se sirven los archivos que estan en la carpeta public,
	// ejm: http://localhost:3003/static/index.html, se sirve /public/index.html
	v1Group.Static("/static", "./public")


	// Inicializar los componentes de la aplicación
	userController, authController, videoController := app.InitializeComponents()

	// Configurar las rutas
	routes.SetupRoutes(v1Group, userController, authController, videoController)
	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))


	r.Run(":3003")

}