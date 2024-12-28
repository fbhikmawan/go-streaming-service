package services

import (
	"encoding/json"
	"fmt"
	"math"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/unbot2313/go-streaming-service/config"
	"github.com/unbot2313/go-streaming-service/internal/models"
)

var (
	rawVideoPathFromWSL = "./static/videos/"
	saveFormatedVideoPath = "./static/temp/"
)

var validVideoExtensions = []string{
	".mp4", ".webm", ".avi", ".mkv", ".mov", ".wmv", ".flv", ".3gp",
}

type VideoService interface {
	SaveVideo(c *gin.Context) (*models.Video, error)
	FormatVideo(videoName string) (string, error) 
	UploadFilesFromFolderToS3(folder string) (importantFiles, string, error)
	DeleteS3Folder(folderName string) error
	GetFilesService() FilesService // New method to access FilesService
	IsValidVideoExtension(c *gin.Context) bool
}


func (vs *videoServiceImp) IsValidVideoExtension(c *gin.Context) bool {

	// Attempt to get the file from the request
	file, err := c.FormFile("video")
	if err != nil {
		return false // The file does not exist or there was an error
	}

	// Get file extension in lowercase letters
	extension := strings.ToLower(filepath.Ext(file.Filename))

	// Check if the extension is valid
	for _, validExtension := range validVideoExtensions {
		if validExtension == extension {
			return true
		}
	}
	return false
}

func (vs *videoServiceImp) GetFilesService() FilesService {
	return vs.FilesService
}

func (vs *videoServiceImp) SaveVideo(c *gin.Context) (*models.Video, error) {
	if err := vs.FilesService.EnsureDir("static/videos"); err != nil {
		return nil, err
	}

	config := config.GetConfig()

	// 1. Get the text fields of the form
	title := c.PostForm("title")
	description := c.PostForm("description")

	// 2. Get the form file
	header, err := c.FormFile("video")

	if err != nil {
		return nil, fmt.Errorf("error getting the file: %w", err)
	}

	storagePath := config.LocalStoragePath

	uuid := uuid.New().String()

	uniqueName := fmt.Sprintf("%s_%s", uuid, header.Filename)

	// Save the file directly with Gin
	savePath := filepath.Join(storagePath, uniqueName)
	if err := c.SaveUploadedFile(header, savePath); err != nil {
		return nil, fmt.Errorf("error saving the file: %w", err)
	}

	// Get video duration
	duration, err := getVideoDuration(savePath)
	if err != nil {
		return nil, fmt.Errorf("error getting video duration: %w", err)
	}

	videoData := &models.Video{
		Id: 			uuid,
		Title:    		title,
		Description:	description,
		Video: 	 		header.Filename,
		LocalPath: 	 	savePath,
		UniqueName: 	uniqueName,
		Duration: 		duration,
	}

	return videoData, nil
}

func (vs *videoServiceImp) FormatVideo(VideoName string) (string, error) {

	//get the name of the video without the extension
	stringName := strings.Split(VideoName, ".")

	//create the folder where the formatted video will be stored
	err := vs.FilesService.CreateFolder("static/temp/" + stringName[0])

	if err != nil {
		return "", fmt.Errorf("error when creating the folder: %w", err)
	}

	saveFormatedPath := saveFormatedVideoPath + stringName[0] + "/output.m3u8"

	videoPath := rawVideoPathFromWSL + VideoName

	// run the ffmpeg command to fragment the video and save it in the folder already created for later uploading to s3
	cmd := exec.Command("ffmpeg", "-i", videoPath, "-c", "copy", "-start_number", "0", "-hls_time", "10", "-hls_list_size", "0", "-f", "hls", saveFormatedPath)

	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf("error when executing the ffmpeg command: %w", err)
	}

	ffmpegFilesPath := saveFormatedVideoPath + stringName[0]
	
	return ffmpegFilesPath, nil

}

func NewVideoService(S3Configuration S3Configuration, filesService FilesService) VideoService {
	return &videoServiceImp{
		S3configuration: S3Configuration,
		FilesService: filesService,
	}
}

type videoServiceImp struct{
	S3configuration S3Configuration
	FilesService FilesService
}

// Function to get the video duration using go-ffprobe
type FFProbeOutput struct {
    Format struct {
        Duration string `json:"duration"`
    } `json:"format"`
}

func formatDuration(seconds float64) string {
    // If less than 60 seconds, return seconds only.
    if seconds < 60 {
        return fmt.Sprintf("%.0fs", seconds)
    }
    
    // Calculate minutes and seconds
    minutes := math.Floor(seconds / 60)
    remainingSeconds := math.Round(seconds - (minutes * 60))
    
    return fmt.Sprintf("%.0f:%.0f", minutes, remainingSeconds)
}

func getVideoDuration(videoPath string) (string, error) {
    // Build the ffprobe command
    cmd := exec.Command("ffprobe",
        "-v", "quiet",
        "-print_format", "json",
        "-show_format",
        videoPath)

    // Ejecutar el comando y obtener la salida
    output, err := cmd.Output()
    if err != nil {
        return "", fmt.Errorf("error running ffprobe: %v", err)
    }

    // Parse JSON output
    var ffprobeOutput FFProbeOutput
    if err := json.Unmarshal(output, &ffprobeOutput); err != nil {
        return "", fmt.Errorf("error parsing the output of ffprobe: %v", err)
    }

    // Convert duration to float64
    seconds, err := strconv.ParseFloat(ffprobeOutput.Format.Duration, 64)
    if err != nil {
        return "", fmt.Errorf("error converting duration to number: %v", err)
    }

    return formatDuration(seconds), nil
}

func SaveThumbnail(videoPath string, folderPath string) (string, error) {
    thumbnailName := fmt.Sprintf("%s.webp", "thumbnail")
    thumbnailPath := filepath.Join(folderPath, thumbnailName)

    // The main change is to move -ss before -i
    cmd := exec.Command("ffmpeg",
        "-ss", "00:00:08",      // MOVED: place before -i
        "-i", videoPath,         // input file
        "-frames:v", "1",        // use frames:v instead of vframes
        "-vf", "scale=480:-1",   
        "-y",                   
        thumbnailPath,          
    )

    if output, err := cmd.CombinedOutput(); err != nil {
        return "", fmt.Errorf("error generating thumbnail: %w, output: %s", err, string(output))
    }

    return thumbnailPath, nil
}