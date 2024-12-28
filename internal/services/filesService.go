package services

import (
	"fmt"
	"log"
	"os"
)

type filesService struct{}

type FilesService interface {
	EnsureDir(dirName string) error
	CreateFolder(path string) error
	RemoveFolder(folder string) error
	RemoveFile(filePath string) error
}

func NewFilesService() FilesService {
	return &filesService{}
}

func (fs *filesService) EnsureDir(dirName string) error {
	err := os.MkdirAll(dirName, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error creating directory: %w", err)
	}

	return nil
}

func (fs *filesService) RemoveFile(filePath string) error {
	err := os.Remove(filePath)

	if err != nil {
		return fmt.Errorf("error when deleting file: %w", err)
	}

	return nil
}


func (fs *filesService) CreateFolder(path string) error {
	// Creates the folder and its parent folders if they do not exist
	err := os.MkdirAll(path, os.ModePerm) // os.ModePerm grants read, write and execute permissions
	if err != nil {
		return fmt.Errorf("error when creating the folder: %w", err)
	}
	return nil
}

func (fs *filesService) RemoveFolder(folder string) error {
	//delete the folder after uploading files
	err := os.RemoveAll(folder)
	if err != nil {
		return err
	}

	log.Println("error when deleting the folder: ", err)
	return err
}