package models

import (
	"time"
)

// Esto es lo que deberia recibir el controlador al crear
// un nuevo usuario
type UserSwagger struct {
	Id         string `json:"id"`
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	Email      string `json:"email"`
}

type User struct {
	Id           string    `json:"id" gorm:"primaryKey;not null;uniqueIndex"`
	Username     string    `json:"username" gorm:"type:varchar(100);not null;uniqueIndex"`
	Password     string    `json:"password" gorm:"not null"`
	Email        string    `json:"email" gorm:"type:varchar(100)"`
	RefreshToken string    `json:"refresh_token"`
	Videos       []VideoModel
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time `gorm:"index"`
}