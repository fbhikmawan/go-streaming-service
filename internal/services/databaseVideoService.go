package services

import (
	"fmt"

	"github.com/unbot2313/go-streaming-service/config"
	"github.com/unbot2313/go-streaming-service/internal/models"
)

func (service *databaseVideoService) FindLatestVideos() (*[]models.VideoModel, error) {
	db, err := config.GetDB()

	if err != nil {
		return nil, err
	}

	var videos []models.VideoModel

	dbCtx := db.Find(&videos)

	fmt.Println(dbCtx.RowsAffected)


	if dbCtx.Error != nil {
		return nil, dbCtx.Error
	}

	return &videos, nil
}

func (service *databaseVideoService) FindVideoByID(videoId string) (*models.VideoModel, error) {
	db, err := config.GetDB()

	if err != nil {
		return nil, err
	}

	var video models.VideoModel

	dbCtx := db.Where("id = ?", videoId).First(&video)

	fmt.Println(dbCtx.RowsAffected)

	if dbCtx.Error != nil {
		return nil, dbCtx.Error
	}

	return &video, nil
}

func (service *databaseVideoService) FindUserVideos(userId string) ([]*models.VideoModel, error) {
	db, err := config.GetDB()

	if err != nil {
		return nil, err
	}

	var videos []*models.VideoModel

	dbCtx := db.Where(&models.VideoModel{UserID: userId}).Find(&videos)

	if dbCtx.Error != nil {
		return nil, dbCtx.Error
	}

	return videos, nil
}

func (service *databaseVideoService) CreateVideo(video *models.Video, userId string) (*models.VideoModel, error) {

	Video := models.VideoModel{
		Id: video.Id,
		Title: video.Title,
		Description: video.Description,
		UserID: userId,
		VideoUrl: video.M3u8FileURL,
	}
	
	db, err := config.GetDB()

	if err != nil {
		return nil, err
	}

	dbCtx := db.Create(&video)

	fmt.Println(dbCtx.RowsAffected)

	fmt.Println(dbCtx.Statement)

	if dbCtx.Error != nil {
		return nil, dbCtx.Error
	}
	
	return &Video, nil
}

func (service *databaseVideoService) UpdateVideo(video *models.VideoModel) (*models.VideoModel, error) {
	return &models.VideoModel{}, nil
}

func (service *databaseVideoService) DeleteVideo(videoId string) error {
	return nil
}


type databaseVideoService struct {}

type DatabaseVideoService interface {
	FindLatestVideos() (*[]models.VideoModel, error)
	FindVideoByID(videoId string) (*models.VideoModel, error) 
	FindUserVideos(userId string) ([]*models.VideoModel, error)
	CreateVideo(video *models.Video, userId string) (*models.VideoModel, error)
	UpdateVideo(video *models.VideoModel) (*models.VideoModel, error)
	DeleteVideo(videoId string) error
}

func NewDatabaseVideoService() DatabaseVideoService {
	return &databaseVideoService{}
}
