package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/unbot2313/go-streaming-service/internal/models"
	"github.com/unbot2313/go-streaming-service/internal/services"
)

type VideoController interface {
	GetLatestVideos(c *gin.Context)
	CreateVideo(c *gin.Context)
	GetVideoByID(c *gin.Context)
	IncrementViews(c *gin.Context)
}

// SaveVideo		godoc
// @Summary 		Save a video
// @Description 	Upload a video file along with metadata (title and description) and save it to the AWS bucket.
// @Tags 			streaming
// @Produce 		json
// @Success 		200 {object} models.VideoSwagger{}
// @Failure 		400 {object} map[string]string
// @Failure 		500 {object} map[string]string
// @Router 			/streaming/ [get]
func (vc *VideoControllerImpl) GetLatestVideos(c *gin.Context) {
	videos, err := vc.databaseVideoService.FindLatestVideos()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, videos)
}


// GetVideoByID		godoc
// @Summary 		Get a video by ID
// @Description 	Get a video by its ID
// @Tags 			streaming
// @Produce 		json
// @Param 			videoid path string true "Video ID"
// @Success 		200 {object} models.VideoSwagger{}
// @Failure 		400 {object} map[string]string
// @Failure 		500 {object} map[string]string
// @Router 			/streaming/id/{videoid} [get]
func (vc *VideoControllerImpl) GetVideoByID(c *gin.Context) {
	videoId := c.Param("videoid")

	video, err := vc.databaseVideoService.FindVideoByID(videoId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, video)
}

// IncrementViews		godoc
// @Summary 		Increment the views of a video
// @Description 	Increment the views of a video by 1
// @Tags 			streaming
// @Produce 		json
// @Param 			videoid path string true "Video ID"
// @Success 		200 {object} models.VideoSwagger{}	
// @Failure 		400 {object} map[string]string
// @Failure 		500 {object} map[string]string
// @Router 			/streaming/views/{videoid} [patch]
func (vc *VideoControllerImpl) IncrementViews(c *gin.Context) {
	videoId := c.Param("videoid")

	video, err := vc.databaseVideoService.IncrementViews(videoId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, video)
	
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
// @Success 		200 {object} models.VideoSwagger{}
// @Failure 		400 {object} map[string]string
// @Failure 		500 {object} map[string]string
// @Router 			/streaming/upload [post]
func (vc *VideoControllerImpl) CreateVideo(c *gin.Context) {

	// Retrieve the user from the context
	user, exists := c.Get("user")
	if !exists {
		c.JSON(500, gin.H{"error": "User not found in context"})
		return
	}

	// Convert to User type
	authenticatedUser, ok := user.(*models.User)
	if !ok {
		c.JSON(500, gin.H{"error": "Failed to parse user data"})
		return
	}

	// check if the file is valid
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

	// save file locally
	videoData, err := vc.videoService.SaveVideo(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// delete the original file
	defer vc.videoService.GetFilesService().RemoveFile(videoData.LocalPath)

	// compress the video
	// pending

	// transfer to .ts and .m3u8 files with ffmpeg and save locally
	filesPath, err := vc.videoService.FormatVideo(videoData.UniqueName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// delete local .ts and .m3u8 files
	defer vc.videoService.GetFilesService().RemoveFolder(filesPath)

	// generate thumbnail of second 1 of the video
	_, err = services.SaveThumbnail(videoData.LocalPath, filesPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// upload the video to s3
	savedDataInS3, baseFolder, err := vc.videoService.UploadFilesFromFolderToS3(filesPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	videoData.M3u8FileURL = savedDataInS3.M3u8FileURL
	videoData.ThumbnailURL = savedDataInS3.ThumbnailURL
	
	// finally, save the url of the video in the database
	Video, err := vc.databaseVideoService.CreateVideo(videoData, authenticatedUser.Id)
	if err != nil {

		// as the video was not saved in the database, it must be deleted from s3
		// as folder/
		defer vc.videoService.DeleteS3Folder(baseFolder + "/")

		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, Video)
}

type VideoControllerImpl struct {
	videoService services.VideoService;
	databaseVideoService services.DatabaseVideoService
}

func NewVideoController(videoService services.VideoService, databaseVideoService services.DatabaseVideoService) VideoController {
	return &VideoControllerImpl{
		videoService: videoService,
		databaseVideoService: databaseVideoService,
	}
}