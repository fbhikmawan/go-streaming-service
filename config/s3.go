package config

import (
	"context"

	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Config struct {
	uploader *manager.Uploader
}

var s3Config S3Config

func GetS3Uploader() *manager.Uploader {

	configOnce.Do(func() {
		cfg, err := awsConfig.LoadDefaultConfig(context.TODO())
		if err != nil {
			panic("unable to load SDK config, " + err.Error())
		}

		client := s3.NewFromConfig(cfg)
		uploader := manager.NewUploader(client)

		s3Config = S3Config{uploader: uploader}
	})

	return s3Config.uploader
}
