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

// SaveVideo		godoc
// @Summary 		Save a video
// @Description 	Upload a video file along with metadata (title and description) and save it to the server.
// @Tags 			videos
// @Produce 		json
// @Success 		200 {object} models.Video
// @Failure 		400 {object} map[string]string
// @Failure 		500 {object} map[string]string
// @Router 			/videos/ [get]
func (vc *VideoControllerImpl) GetVideos(c *gin.Context) {
	vc.videoService.GetVideos()
	fmt.Println("GetVideos")
}

// SaveVideo		godoc
// @Summary 		Save a video
// @Description 	Upload a video file along with metadata (title and description) and save it to the server.
// @Tags 			videos
// @Accept 			multipart/form-data
// @Produce 		json
// @Param 			title formData string true "Video Title"
// @Param 			description formData string false "Video Description"
// @Param 			video formData file true "Video File"
// @Success 		200 {object} models.Video
// @Failure 		400 {object} map[string]string
// @Failure 		500 {object} map[string]string
// @Router 			/videos/ [post]
func (vc *VideoControllerImpl) CreateVideo(c *gin.Context) {

	videoData, err := vc.videoService.SaveVideo(c)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, videoData)
	

}

type VideoControllerImpl struct {
	videoService services.VideoService
}

func NewVideoController(videoService services.VideoService) VideoController {
	return &VideoControllerImpl{
		videoService: videoService,
	}
}