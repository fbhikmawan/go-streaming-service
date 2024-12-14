package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/unbot2313/go-streaming-service/internal/models"
	"github.com/unbot2313/go-streaming-service/internal/services"
)

type UserControllerImp struct {
	service	 services.UserService
}

type UserController interface {
	CreateUser(c *gin.Context)
	GetUserByID(c *gin.Context)
	DeleteUserByID(c *gin.Context)
}


// GetUserByID		godoc
// @Summary 		Get user by ID
// @Description 	Search user by ID in Db
// @Tags 			users
// @Param 			Id path string true "User ID"
// @Produce 		json
// @Success 		200 {object} models.User{}
// @Router 			/users/{UserId} [get]
func (controller *UserControllerImp) GetUserByID(c *gin.Context) {
	Id := c.Param("Id")
	users, err := controller.service.GetUserByID(Id)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, users)
}


// CreateUser		godoc
// @Summary 		Create a new user
// @Description 	Save user in Db
// @Tags 			users
// @Accept 			json
// @Param 			user body models.UserSwagger{} true "User object containing all user details"
// @Produce 		json
// @Success 		200 {object} models.User{}
// @Router 			/users/ [post]
func (controller *UserControllerImp) CreateUser(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	newUser, err := controller.service.CreateUser(&user)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "User created", "user": newUser})
}

// GetUserByID		godoc
// @Summary 		Delete user by ID
// @Description 	Delete user by ID ni Db
// @Tags 			users
// @Param 			Id path string true "User ID"
// @Produce 		json
// @Success 		200 {object} models.User{}
// @Router 			/users/{UserId} [delete]
func (controller *UserControllerImp) DeleteUserByID(c *gin.Context) {

	id := c.Param("Id")

	err := controller.service.DeleteUserByID(id)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "User deleted"})
}




func NewUserController(service services.UserService) *UserControllerImp {
	return &UserControllerImp{service: service}
}