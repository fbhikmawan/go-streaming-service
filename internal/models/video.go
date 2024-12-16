package models

import (
	"time"

	"gorm.io/gorm"
)

type Video struct {
	Id          	string
	Video       	string
	Title       	string
	Description 	string
	LocalPath       string
	UniqueName  	string
	M3u8FileURL  	string

}


// Tipo para usar en la documentacion con Swaggo
// ya que no reconoce los tags de gorm
type VideoSwagger struct {
	Id          	string    `json:"id" gorm:"primaryKey;not null;uniqueIndex"`
	VideoUrl       	string    `json:"video" gorm:"not null"`
	Title       	string    `json:"title" gorm:"type:varchar(100);not null"`
	Description 	string    `json:"description"`
	UserID			string		`json:"user_id" gorm:"not null"`
}

type VideoModel struct {
	Id			string		`json:"id" gorm:"primaryKey;not null;unique_index"`
	VideoUrl	string		`json:"video" gorm:"not null"`
	Title		string		`json:"title" gorm:"type:varchar(100);not null"`
	Description	string		`json:"description"`
	UserID		string		`json:"user_id" gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}