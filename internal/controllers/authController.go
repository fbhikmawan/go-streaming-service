package controllers

import (
	"fmt"
    "errors"
    "net/http"
    "regexp"

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

    // Retrieve user details
    user, err := controller.userService.GetUserByUserName(userLogin.Username)
    if err != nil {
        c.JSON(404, gin.H{"error": "User not found"})
        return
    }

	c.JSON(200, gin.H{
		"token": token,
        "user": gin.H{
            "id": user.Id,
            "username": user.Username,
            "email": user.Email,
        },
	})
}

// Register is the controller for the register endpoint
func (controller *AuthControllerImp) Register(c *gin.Context) {
    var userInput models.UserCreate

    // Bind JSON input
    if err := c.ShouldBindJSON(&userInput); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Validate input
    if err := validateUserInput(userInput); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Create a new user
    newUser := models.User{
        Username: userInput.Username,
        Password: userInput.Password,
        Email:    userInput.Email,
    }

    // Use the userService to create the user
    createdUser, err := controller.userService.CreateUser(&newUser)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
        return
    }

    // Generate JWT token
    token, err := controller.authService.GenerateToken(createdUser)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
        return
    }

    // Return response
    c.JSON(http.StatusOK, gin.H{
        "message": "User registered successfully",
        "token":   token,
        "user":    createdUser.Username,
    })
}

type AuthControllerImp struct {
    authService services.AuthService
    userService services.UserService
}

// NewAuthController creates a new instance of the authentication controller
func NewAuthController(authService services.AuthService, userService services.UserService) AuthController {
    return &AuthControllerImp{
        authService: authService,
        userService: userService,
    }
}

// Helper functions
func validateUserInput(input models.UserCreate) error {
    if len(input.Username) < 3 || len(input.Username) > 50 {
        return errors.New("username must be between 3 and 50 characters")
    }
    if len(input.Password) < 8 {
        return errors.New("password must be at least 8 characters long")
    }
    if !isValidEmailFormat(input.Email) {
        return errors.New("invalid email format")
    }
    return nil
}
func isValidEmailFormat(email string) bool {
    emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
    return emailRegex.MatchString(email)
}