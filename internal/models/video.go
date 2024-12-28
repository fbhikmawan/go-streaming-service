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
	Duration   		string	
	ThumbnailURL 	string
}


// Type to use in the documentation with Swaggo
// since it does not recognize gorm tags
type VideoSwagger struct {
	Id          	string    	`json:"id" gorm:"primaryKey;not null;uniqueIndex"`
	VideoUrl       	string    	`json:"video" gorm:"not null"`
	Title       	string    	`json:"title" gorm:"type:varchar(100);not null"`
	Description 	string    	`json:"description"`
	UserID			string		`json:"user_id" gorm:"not null"`
	Duration   		string	 	`json:"duration"`
	ThumbnailURL 	string   	`json:"thumbnail"`
	Views 			uint		`json:"views" gorm:"default:0"`

}

// the one used in the db
type VideoModel struct {
	Id				string			`json:"id" gorm:"primaryKey;not null;uniqueIndex"`
	VideoUrl		string			`json:"video" gorm:"not null"`
	Title			string			`json:"title" gorm:"type:varchar(100);not null"`
	Description		string			`json:"description"`
	UserID			string			`json:"user_id" gorm:"not null"`
	Duration   		string	 		`json:"duration"`
	ThumbnailURL 	string   		`json:"thumbnail"`
	Views 			uint			`json:"views" gorm:"default:0"`
	CreatedAt 		time.Time
	UpdatedAt		time.Time
	DeletedAt 		gorm.DeletedAt 	`gorm:"index"`
}

// videomodel table name
func (VideoModel) TableName() string {
    return "videos"
}
