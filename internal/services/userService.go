package services

import "github.com/unbot2313/go-streaming-service/internal/models"

type UserServiceImp struct{}

type UserService interface {
	GetUserByID(Id string) (*models.User, error)
}

func (service *UserServiceImp) GetUserByID(Id string) (*models.User, error) {
	return &models.User{}, nil
}

func NewUserService() UserService {
	return &UserServiceImp{}
}