package services

import (
	"context"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/unbot2313/go-streaming-service/config"
)

type S3Configuration struct {
	Region string
	BucketName string
	AccessKey string
	SecretKey string
	Uploader *manager.Uploader
}

func GetS3Configuration() S3Configuration {

	Config := config.GetConfig()

	return S3Configuration{
		Region: Config.AWSRegion,
		BucketName: Config.AWSBucketName,
		AccessKey: Config.AWSAccessKey,
		SecretKey: Config.AWSSecretKey,
		Uploader: config.GetS3Uploader(),
	}
}

type S3ServiceImp struct {
	configuration S3Configuration
}

type S3Service interface {
	UploadFilesFromFolderToS3(folder string) ([]string, error)
}

func NewS3Service(configuration S3Configuration) S3Service {
	return &S3ServiceImp{configuration: configuration}
}

func (s3Service *S3ServiceImp) UploadFilesFromFolderToS3(folder string) ([]string, error) {

	// Obtener el nombre de la carpeta actual
	baseFolder := filepath.Base(folder)

	var uploadedFiles []string

	files, err := os.ReadDir(folder)

	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filePath := filepath.Join(folder, file.Name())
		f, err := os.Open(filePath)
		if err != nil {
			return nil, err
		}
		defer f.Close()

		// Construir el Key como nombre de la carpeta + nombre del archivo
		key := filepath.Join(baseFolder, file.Name())

		// Subir el archivo a S3
		result, errS3 := s3Service.configuration.Uploader.Upload(context.TODO(), &s3.PutObjectInput{
			Bucket: aws.String(s3Service.configuration.BucketName),
			Key:    aws.String(key),
			Body:   f,
			// ACL:    "public-read",
		})

		if errS3 != nil {
			return nil, errS3
		}

		// Guardar la URL del archivo subido
		uploadedFiles = append(uploadedFiles, result.Location)
	}

	//borrar la carpeta tras subir los archivos
	os.RemoveAll(folder)
	
	return uploadedFiles, nil
}