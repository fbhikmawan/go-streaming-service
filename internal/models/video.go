package models

import (
	"gorm.io/gorm"
)

type Video struct {
	Id          	string
	Video       	string
	Title       	string
	Description 	string
	LocalPath       string
	UniqueName  	string
	S3FilesPath  	[]string

}

type VideoModel struct {
	gorm.Model

	Id			string		`json:"id" gorm:"primaryKey;not null;unique_index"`
	Video		string		`json:"video" gorm:"not null"`
	Title		string		`json:"title" gorm:"type:varchar(100);not null"`
	Description	string		`json:"description"`
	UserID		uint		`json:"user_id" gorm:"not null"`
	
}