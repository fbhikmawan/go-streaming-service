package models

import "gorm.io/gorm"

// Esto es lo que deberia recibir el controlador al crear
// un nuevo usuario
type UserSwagger struct {
	ID         string `json:"id"`
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	Email      string `json:"email"`
}

type User struct {
	gorm.Model
	ID		   string    `json:"id" gorm:"primaryKey;not null;unique_index"`
	Username   string    `json:"username" gorm:"type:varchar(100);not null"`
	Password   string    `json:"password" gorm:"not null"`
	Email      string    `json:"email" gorm:"type:varchar(100)"`
	RefreshToken string	 `json:"refresh_token"`
	Videos	 []VideoModel
}
