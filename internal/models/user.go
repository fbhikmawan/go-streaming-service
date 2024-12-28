package models

import (
	"time"

	"gorm.io/gorm"
)

// This is what the controller should receive when creating a new user
// a new user
type UserCreate struct {
	Id         string `json:"id"`
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	Email      string `json:"email"`
}

type UserLogin struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// UserSwagger is used in the documentation because Swaggo doesn't recognize
// gorms tags
type UserSwagger struct {
	Id           string    `json:"id" gorm:"primaryKey;not null;uniqueIndex"`
	Username     string    `json:"username" gorm:"type:varchar(100);not null;uniqueIndex"`
	Password     string    `json:"password" gorm:"not null"`
	Email        string    `json:"email" gorm:"type:varchar(100)"`
	RefreshToken string    `json:"refresh_token"`
	Videos []VideoSwagger 	`json:"videos" gorm:"foreignKey:UserID"`
}

// the one used in the db
type User struct {
	Id           string    `json:"id" gorm:"primaryKey;not null;uniqueIndex"`
	Username     string    `json:"username" gorm:"type:varchar(100);not null;uniqueIndex"`
	Password     string    `json:"password" gorm:"not null"`
	Email        string    `json:"email" gorm:"type:varchar(100)"`
	RefreshToken string    `json:"refresh_token"`
	Videos 		 []VideoModel 	`json:"videos" gorm:"foreignKey:UserID"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}