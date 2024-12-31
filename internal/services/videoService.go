package services

import (
	"encoding/json"
	"fmt"
	"math"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
    "log"
	"os"
	"bytes"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/unbot2313/go-streaming-service/config"
	"github.com/unbot2313/go-streaming-service/internal/models"
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
	log.Println("Starting video save process")

	config := config.GetConfig()

	// 1. Get the text fields of the form
	title := c.PostForm("title")
	description := c.PostForm("description")
    // Validate required fields
    if title == "" {
        log.Println("Missing required field: title")
		return nil, fmt.Errorf("error missing required field: title")
    }
    if description == "" {
        log.Println("Missing required field: description")
		return nil, fmt.Errorf("error missing required field: description")
    }
	log.Printf("Title: %s, Description: %s", title, description)

	// 2. Get the form file
	header, err := c.FormFile("video")
	if err != nil {
		log.Printf("Error getting form file: %v", err)
		return nil, fmt.Errorf("error getting the file: %w", err)
	}
	log.Printf("File name: %s, Size: %d bytes", header.Filename, header.Size)

	storagePath := config.LocalStoragePath
	log.Printf("Storage path: %s", storagePath)

	uuid := uuid.New().String()
	uniqueName := fmt.Sprintf("%s_%s", uuid, header.Filename)
	log.Printf("Unique file name: %s", uniqueName)

	// Save the file directly with Gin
	savePath := filepath.Join(storagePath, uniqueName)
	log.Printf("Save path: %s", savePath)
	if err := c.SaveUploadedFile(header, savePath); err != nil {
		log.Printf("Error saving uploaded file: %v", err)
		return nil, fmt.Errorf("error saving the file: %w", err)
	}
	log.Printf("File saved successfully at: %s", savePath)

	// Get video duration
	duration, err := getVideoDuration(savePath)
	if err != nil {
		log.Printf("Error getting video duration: %v", err)
		return nil, fmt.Errorf("error getting video duration: %w", err)
	}
	log.Printf("Video duration: %s", duration)

	videoData := &models.Video{
		Id: 			uuid,
		Title:    		title,
		Description:	description,
		Video: 	 		header.Filename,
		LocalPath: 	 	savePath,
		UniqueName: 	uniqueName,
		Duration: 		duration,
	}

	log.Println("Video saved successfully")
	return videoData, nil
}

func (vs *videoServiceImp) FormatVideo(VideoName string) (string, error) {
	log.Printf("Starting video formatting process for: %s", VideoName)

	config := config.GetConfig()

	//get the name of the video without the extension
	stringName := strings.Split(VideoName, ".")
	log.Printf("Extracted video name without extension: %s", stringName[0])

	//create the folder where the formatted video will be stored
	storagePath := config.LocalStoragePath
	log.Printf("Storage path: %s", storagePath)

	// Ensure the output directory exists
    outputPath := filepath.Join(storagePath, "videos")
    if err := os.MkdirAll(outputPath, 0755); err != nil {
        return "", fmt.Errorf("failed to create output directory: %w", err)
    }

    // Use absolute paths
    absOutputPath, err := filepath.Abs(filepath.Join(outputPath, "output.m3u8"))
    if err != nil {
        return "", fmt.Errorf("failed to get absolute output path: %w", err)
    }

    absVideoPath, err := filepath.Abs(filepath.Join(storagePath, VideoName))
    if err != nil {
        return "", fmt.Errorf("failed to get absolute video path: %w", err)
    }

    log.Printf("Formatted video path: %s", absOutputPath)
    log.Printf("Original video path: %s", absVideoPath)

    // Check if the file exists and is accessible
    if err := checkFile(absVideoPath); err != nil {
        return "", fmt.Errorf("error checking input file: %w", err)
    }

	// Run the ffmpeg command
    cmd := exec.Command("ffmpeg", "-i", absVideoPath, "-c", "copy", "-start_number", "0", "-hls_time", "10", "-hls_list_size", "0", "-f", "hls", absOutputPath)
    log.Printf("Executing ffmpeg command: %v", cmd.Args)

    var out bytes.Buffer
    var stderr bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &stderr

    err = cmd.Run()
    if err != nil {
        log.Printf("Error executing ffmpeg command: %v", err)
        log.Printf("FFmpeg stdout:\n%s", out.String())
        log.Printf("FFmpeg stderr:\n%s", stderr.String())
        return "", fmt.Errorf("error when executing the ffmpeg command: %w\nstdout: %s\nstderr: %s", err, out.String(), stderr.String())
    }
    
    log.Printf("FFmpeg execution successful. stdout:\n%s", out.String())
    log.Printf("FFmpeg execution successful. stderr:\n%s", stderr.String())

    return outputPath, nil

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

// Helper function
func checkFile(filePath string) error {
    _, err := os.Stat(filePath)
    if err != nil {
        if os.IsNotExist(err) {
            return fmt.Errorf("file does not exist: %s", filePath)
        }
        return fmt.Errorf("unable to access file: %s", filePath)
    }
    return nil
}
