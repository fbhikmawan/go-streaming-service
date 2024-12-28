package controllers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/unbot2313/go-streaming-service/internal/models"
	"github.com/unbot2313/go-streaming-service/internal/services"
)

type AuthController interface {
	Login(c *gin.Context)
	Register(c *gin.Context)
}

// GetUserByUserName		godoc
// @Summary 				Log in user
// @Tags 					Auth
// @Produce 				json
// @Accept 					json
// @Param 					user body models.UserLogin{} true "User object containing all user details"
// @Success 				200 {object} map[string]string
// @Failure 				404 {object} map[string]string
// @Failure 				500 {object} map[string]string
// @Router 					/auth/login [post]
func (controller *AuthControllerImp) Login(c *gin.Context) {
	
	var userLogin models.UserLogin

	if err := c.ShouldBindJSON(&userLogin); err != nil {
		err = fmt.Errorf("a user name and password are required")
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	token, err := controller.authService.Login(userLogin.Username, userLogin.Password)

	if err != nil {
		c.JSON(404, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{
		"token": token,
		"user": userLogin.Username,
	})
}

// Register is the controller for the register endpoint
func (controller *AuthControllerImp) Register(c *gin.Context) {
	// Implementing
}


type AuthControllerImp struct {
	authService services.AuthService
}

// NewAuthController creates a new instance of the authentication controller
func NewAuthController(authService services.AuthService) AuthController {
	return &AuthControllerImp{authService}
}
