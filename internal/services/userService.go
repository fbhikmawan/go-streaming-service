package services

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/unbot2313/go-streaming-service/config"
	"github.com/unbot2313/go-streaming-service/internal/models"
	"gorm.io/gorm"
)

type UserServiceImp struct{}

type UserService interface {
	GetUserByID(Id string) (*models.User, error)
	GetUserByUserName(userName string) (*models.User, error)
	CreateUser(user *models.User) (*models.User, error)
	DeleteUserByID(Id string) error
	// Pendiente
	UpdateUserByID(Id string, user *models.User) (*models.User, error)
}

func (service *UserServiceImp) GetUserByID(Id string) (*models.User, error) {
	var user models.User

	// Get the connection to the database
	db, err := config.GetDB()
	if err != nil {
		return nil, err
	}

	// Searches for the user by ID and includes the associated videos
	err = db.Preload("Videos").First(&user, "id = ?", Id).Error

	// Handles the user not found case
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("user with ID %s not found", Id)
	}

	// Handles any other errors
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (service *UserServiceImp) GetUserByUserName(userName string) (*models.User, error) {
	var user models.User

	// Get the connection to the database
	db, err := config.GetDB()
	if err != nil {
		return nil, err
	}

	// Search for the user by username and include the associated videos.
	err = db.Preload("Videos").First(&user, "username = ?", userName).Error

	// Handles the user not found case
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("user with username %s not found", userName)
	}

	// Handles any other errors
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (service *UserServiceImp) CreateUser(user *models.User) (*models.User, error) {
	db, err := config.GetDB()
	if err != nil {
		return nil, err
	}

	user.Id = uuid.New().String()

	hashedPassword, err := HashPassword(user.Password)

	if err != nil {
		return nil, err
	}

	user.Password = hashedPassword

	dbCtx := db.Create(user)

	if errors.Is(dbCtx.Error, gorm.ErrDuplicatedKey) {
		return nil, fmt.Errorf("user with username %s already exists", user.Username)
	}

	if dbCtx.Error != nil {
		return nil, dbCtx.Error
	}

	return user, nil
}

func (service *UserServiceImp) DeleteUserByID(Id string) error {

	db, err := config.GetDB()
	if err != nil {
		return err
	}

	// Deletes the user using the custom field `Id`.
    dbCtx := db.Where("id = ?", Id).Delete(&models.User{})

	if errors.Is(dbCtx.Error, gorm.ErrRecordNotFound) {
		return fmt.Errorf("user with ID %s not found", Id)
	}

    if dbCtx.Error != nil {
        return dbCtx.Error
    }

	return nil
}

// Pending
func (service *UserServiceImp) UpdateUserByID(Id string, user *models.User) (*models.User, error) {
	return nil, nil
}


func NewUserService() UserService {
	return &UserServiceImp{}
}