package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/unbot2313/go-streaming-service/internal/services"
)

type VideoController interface {
	GetVideos(c *gin.Context)
	CreateVideo(c *gin.Context)
}

type VideoControllerImpl struct {
	videoService services.VideoService
}

func NewVideoController(videoService services.VideoService) VideoController {
	return &VideoControllerImpl{
		videoService: videoService,
	}
}

func (vc *VideoControllerImpl) GetVideos(c *gin.Context) {
	vc.videoService.GetVideos()
	fmt.Println("GetVideos")
}

func (vc *VideoControllerImpl) CreateVideo(c *gin.Context) {

	videoData, err := vc.videoService.SaveVideo(c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, videoData)
	

}
