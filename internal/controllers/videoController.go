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
// @Description 	Upload a video file along with metadata (title and description) and save it to the AWS bucket.
// @Tags 			streaming
// @Produce 		json
// @Success 		200 {object} models.Video
// @Failure 		400 {object} map[string]string
// @Failure 		500 {object} map[string]string
// @Router 			/streaming/ [get]
func (vc *VideoControllerImpl) GetVideos(c *gin.Context) {
	vc.videoService.GetVideos()
	fmt.Println("GetVideos")
}

// SaveVideo		godoc
// @Summary 		Save a video
// @Description 	Upload a video file along with metadata (title and description) and save to the AWS bucket.
// @Tags 			streaming
// @Accept 			multipart/form-data
// @Produce 		json
// @Param 			title formData string true "Video Title"
// @Param 			description formData string false "Video Description"
// @Param 			video formData file true "Video File"
// @Success 		200 {object} models.Video
// @Failure 		400 {object} map[string]string
// @Failure 		500 {object} map[string]string
// @Router 			/streaming/upload [post]
func (vc *VideoControllerImpl) CreateVideo(c *gin.Context) {

	// verificar si el archivo es válido
	if !vc.videoService.IsValidVideoExtension(c) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El archivo no es un tipo de video válido."})
		return
	}

	fileSize := c.Request.ContentLength
	const maxFileSize = 100 * 1024 * 1024 // 100 MB
	if fileSize > maxFileSize {
		c.JSON(http.StatusBadRequest, gin.H{"error": "El archivo excede el límite de tamaño permitido."})
		return
	}

	// guardar archivo en local
	videoData, err := vc.videoService.SaveVideo(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// borrar el archivo original
	defer vc.videoService.GetFilesService().RemoveFile(videoData.LocalPath)

	// comprimir el video
	// pendiente

	//pasar a archivos .ts y .m3u8 con ffmpeg y guardarlo en local
	filesPath ,err := vc.videoService.FormatVideo(videoData.UniqueName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// borrar archivos locales .ts y .m3u8
	defer vc.videoService.GetFilesService().RemoveFolder(filesPath)



	// subir el video a s3
	savedDataInS3, err := vc.videoService.UploadFilesFromFolderToS3(filesPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	

	// finalmente, guardar la url del video en la base de datos
	// pendiente

	videoData.S3FilesPath = savedDataInS3

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