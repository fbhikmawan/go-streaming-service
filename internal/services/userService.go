package services

import (
	"github.com/google/uuid"
	"github.com/unbot2313/go-streaming-service/config"
	"github.com/unbot2313/go-streaming-service/internal/models"
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
	var user *models.User

	db, err := config.GetDB()
	if err != nil {
		return nil, err
	}

	db.First(&user, Id)

	// Por alguna razon no devuelve el id del usuario
	return user, nil
}

func (service *UserServiceImp) GetUserByUserName(userName string) (*models.User, error) {
	var user *models.User
	db, err := config.GetDB()
	if err != nil {
		return nil, err
	}

	dbCtx := db.Where("user_name = ?", userName).First(&user)

	if dbCtx.Error != nil {
		return nil, dbCtx.Error
	}

	return user, nil
}

func (service *UserServiceImp) CreateUser(user *models.User) (*models.User, error) {
	db, err := config.GetDB()
	if err != nil {
		return nil, err
	}

	user.Id = uuid.New().String()

	newUser := db.Create(user)

	if newUser.Error != nil {
		return nil, newUser.Error
	}

	return user, nil
}

func (service *UserServiceImp) DeleteUserByID(Id string) error {

	db, err := config.GetDB()
	if err != nil {
		return err
	}

	 // Elimina el usuario usando el campo personalizado `Id`
    userDeleted := db.Where("id = ?", Id).Delete(&models.User{})
    if userDeleted.Error != nil {
        return userDeleted.Error
    }

	return nil
}

// Pendiente
func (service *UserServiceImp) UpdateUserByID(Id string, user *models.User) (*models.User, error) {
	return nil, nil
}


func NewUserService() UserService {
	return &UserServiceImp{}
}