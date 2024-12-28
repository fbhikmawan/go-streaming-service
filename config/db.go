package config

import (
	"fmt"
	"sync"

	"github.com/unbot2313/go-streaming-service/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	dbInstance *gorm.DB
	once       sync.Once
)

// GetDsn generates the connection string for the database.
func getDsn() string {

	config := GetConfig()

	dockerMode := config.DOCKER_MODE
	// If the application is running in a Docker container, the database host must be changed.
	host := config.PostgresHost
	// if it is not in a container it is searched for outside the docker container connection
	if !dockerMode {
		host = "localhost"
	}
	port := config.PostgresPort
	user := config.PostgresUser
	password := config.PostgresPassword
	dbname := config.PostgresDBName
	
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbname, port)
}

// GetDB returns a single instance of the database connection.
func GetDB() (*gorm.DB, error) {
	var err error
	once.Do(func() {
		dsn := getDsn()
		dbInstance, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	})
	if err != nil {
		return nil, err
	}

	// Migrates the tables to the database.
	err = migrations(dbInstance)
	if err != nil {
		return nil, err
	}
	return dbInstance, nil
}

func migrations(db *gorm.DB) error {
	err := db.AutoMigrate(&models.User{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&models.VideoModel{})
	if err != nil {
		return err
	}

	return nil
}
