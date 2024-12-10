package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/unbot2313/go-streaming-service/internal/services"
)

type S3Controller interface {
	UploadVideos (c *gin.Context)
}

type S3ControllerImpl struct {
	s3Service services.S3Service
}

func NewS3Controller(s3Service services.S3Service) S3Controller {
	return &S3ControllerImpl{s3Service}
}

// UploadVideos godoc
// @Summary Upload a video to S3
// @Description Upload a video file to an S3 bucket.
// @Tags s3
// @Accept multipart/form-data
// @Produce json
// @Param video formData file true "Video File"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /s3/upload [post]
func (sc *S3ControllerImpl) UploadVideos(c *gin.Context) {
	
	videoData, err := sc.s3Service.UploadFilesFromFolderToS3("static/temp")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, videoData)
}