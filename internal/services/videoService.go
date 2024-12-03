package services

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/unbot2313/go-streaming-service/internal/models"
)

type videoServiceImp struct{}

type VideoService interface {
	GetVideos()
	SaveVideo(c *gin.Context) (*models.Video, error)
	ensureDir(dirName string) error
}

func NewVideoService() VideoService {
	return &videoServiceImp{}
}

func (vs *videoServiceImp) GetVideos() {
	fmt.Println("GetVideos")
}

func (vs *videoServiceImp) SaveVideo(c *gin.Context) (*models.Video, error) {
	if err := vs.ensureDir("videos"); err != nil {
		return nil, err
	}

	// 1. Obtener los campos de texto del formulario
	title := c.PostForm("title")
	description := c.PostForm("description")

	// 2. Obtener el archivo del formulario
	header, err := c.FormFile("video")

	if err != nil {
		return nil, fmt.Errorf("error al obtener el archivo: %w", err)
	}

	storagePath := "videos"
	uniqueName := fmt.Sprintf("%s_%s", uuid.New().String(), header.Filename)

	// Guardar el archivo directamente con Gin
	savePath := filepath.Join(storagePath, uniqueName)
	if err := c.SaveUploadedFile(header, savePath); err != nil {
		return nil, fmt.Errorf("error al guardar el archivo: %w", err)
	}

	videoData := &models.Video{
		Title:       title,
		Description: description,
		Video: 	 header.Filename,
		Path: 	 savePath,
	}

	return videoData, nil

}

func (vs *videoServiceImp) ensureDir(dirName string) error {
	err := os.MkdirAll(dirName, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error al crear directorio: %w", err)
	}

	return nil
}