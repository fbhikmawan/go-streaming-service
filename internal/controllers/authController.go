package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/unbot2313/go-streaming-service/internal/services"
)

type AuthControllerImp struct {
	authService services.AuthService
}

type AuthController interface {
	Login(c *gin.Context)
	Register(c *gin.Context)
}

// NewAuthController crea una nueva instancia del controlador de autenticación
func NewAuthController(authService services.AuthService) AuthController {
	return &AuthControllerImp{authService}
}

// Login es el controlador para el endpoint de login
func (controller *AuthControllerImp) Login(c *gin.Context) {
	// Implementar
}

// Register es el controlador para el endpoint de registro
func (controller *AuthControllerImp) Register(c *gin.Context) {
	// Implementar
}