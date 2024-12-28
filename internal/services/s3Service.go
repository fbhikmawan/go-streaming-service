package services

// videoService extension focused on S3 file handling

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/unbot2313/go-streaming-service/config"
)

type importantFiles struct {
	M3u8FileURL string
	ThumbnailURL string
}

func (s3Service *videoServiceImp) UploadFilesFromFolderToS3(folder string) (
	importantFiles,
	string,
	error,
) {

	// Get the name of the current folder
	baseFolder := filepath.Base(folder)

	var m3u8FileURL string

	var thumbnailURL string

	files, err := os.ReadDir(folder)

	if err != nil {
		return importantFiles{}, baseFolder, err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filePath := filepath.Join(folder, file.Name())
		f, err := os.Open(filePath)
		if err != nil {
			return importantFiles{}, baseFolder, err
		}
		defer f.Close()

		// Construct the Key as folder name + file name
		key := filepath.Join(baseFolder, file.Name())

		// Upload the file to S3
		result, errS3 := s3Service.S3configuration.Uploader.Upload(context.TODO(), &s3.PutObjectInput{
			Bucket: aws.String(s3Service.S3configuration.BucketName),
			Key:    aws.String(key),
			Body:   f,
			// ACL:    "public-read",
		})

		if errS3 != nil {
			return importantFiles{}, baseFolder, errS3
		}

		// If it is an m3u8 file, save its URL to the database
        if strings.HasSuffix(file.Name(), ".m3u8") {
            m3u8FileURL = result.Location
        }

		if strings.HasSuffix(file.Name(), ".webp") {
			thumbnailURL = result.Location
		}

	}

	if m3u8FileURL == "" {
		return importantFiles{}, baseFolder, fmt.Errorf("file .m3u8 not found")
    }
	
	return importantFiles{
		M3u8FileURL: m3u8FileURL,
		ThumbnailURL: thumbnailURL,
	}, baseFolder, nil
}

// DeleteFolder will delete all objects inside the specified “folder”.
func (s3Service *videoServiceImp) DeleteS3Folder(folderName string) error {
    ctx := context.Background() // Define the context

	log.Println("Deleting objects in the folder: ", folderName)


    // List the objects inside the “folder”.
    input := &s3.ListObjectsV2Input{
        Bucket: aws.String(s3Service.S3configuration.BucketName),
        Prefix: aws.String(folderName), // Prefijo de la "carpeta"
    }
    
    objectPaginator := s3.NewListObjectsV2Paginator(s3Service.S3configuration.Client, input)
    var objectsToDelete []types.ObjectIdentifier

    // Go through all objects within the prefix (folder) and add them to the deletion list
    for objectPaginator.HasMorePages() {
        output, err := objectPaginator.NextPage(ctx)
        if err != nil {
            log.Printf("Error when listing objects: %v\n", err)
            return err
        }

        for _, object := range output.Contents {
            objectsToDelete = append(objectsToDelete, types.ObjectIdentifier{
                Key: object.Key,
            })
        }
    }

    if len(objectsToDelete) == 0 {
        log.Println("No objects were found to remove.")
        return nil
    }

    // Delete listed objects
    deleteInput := &s3.DeleteObjectsInput{
        Bucket: aws.String(s3Service.S3configuration.BucketName),
        Delete: &types.Delete{
            Objects: objectsToDelete,
        },
    }

    _, err := s3Service.S3configuration.Client.DeleteObjects(ctx, deleteInput)
    if err != nil {
        log.Printf("Error when deleting objects: %v\n", err)
        return err
    }

    log.Printf("Objects in the folder %v.\n", folderName)
    return nil
}

// ListObjects lists the objects in a bucket.
func (s3Service *videoServiceImp) ListObjects(ctx context.Context, bucketName string, folder string) ([]types.Object, error) {
	var err error
	var output *s3.ListObjectsV2Output
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
		Prefix: aws.String(folder),
	}
	var objects []types.Object
	objectPaginator := s3.NewListObjectsV2Paginator(s3Service.S3configuration.Client, input)
	for objectPaginator.HasMorePages() {
		output, err = objectPaginator.NextPage(ctx)
		if err != nil {
			var noBucket *types.NoSuchBucket
			if errors.As(err, &noBucket) {
				log.Printf("Bucket %s does not exist.\n", bucketName)
				err = noBucket
			}
			break
		} else {
			objects = append(objects, output.Contents...)
		}
	}
	return objects, err
}

// Configuration
type S3Configuration struct {
	Region string
	BucketName string
	AccessKey string
	SecretKey string
	Client *s3.Client
	Uploader *manager.Uploader
}

func GetS3Configuration() S3Configuration {

	Config := config.GetConfig()

	return S3Configuration{
		Region: Config.AWSRegion,
		BucketName: Config.AWSBucketName,
		AccessKey: Config.AWSAccessKey,
		SecretKey: Config.AWSSecretKey,
		Client: config.GetS3Client(),
		Uploader: config.GetS3Uploader(),
	}
}
