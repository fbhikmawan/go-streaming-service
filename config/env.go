package config

import (
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	Port         string
	DatabaseURL  string
	JWTSecretKey string
	AWSRegion	 string
	AWSBucketName string
	AWSAccessKey string
	AWSSecretKey string
	LocalStoragePath string

	DOCKER_MODE 	bool

	PostgresHost	 string
	PostgresPort	 string
	PostgresUser	 string
	PostgresPassword string
	PostgresDBName	 string
}


// el singleton de configuracion 
var (
	config     *Config
	configOnce sync.Once
)


func GetConfig() *Config {

	// use Sync.Once to ensure that the configuration is loaded only once and avoid performance issues
	// and use singleton to ensure that there is only one instance of the configuration in the entire application.

	configOnce.Do(func() {
		err := loadEnv()
		if err != nil {
			panic(fmt.Sprintf("Error loading .env file: %v", err))
		}

		config = &Config{
			Port:         getEnv("PORT", "8080"),
			JWTSecretKey: getEnv("JWT_SECRET_KEY", "secretJwtKey"),
			LocalStoragePath: getEnv("LOCAL_STORAGE_PATH", "videos"),
			AWSRegion:    getEnv("AWS_REGION", ""),
			AWSBucketName: getEnv("AWS_BUCKET_NAME", ""),
			AWSAccessKey: getEnv("AWS_ACCESS_KEY_ID", ""),
			AWSSecretKey: getEnv("AWS_SECRET_ACCESS_KEY", ""),

			DOCKER_MODE: getEnvAsBool("DOCKER_MODE", false),

			PostgresHost: getEnv("POSTGRES_HOST", "localhost"),
			PostgresPort: getEnv("POSTGRES_PORT", "5432"),
			PostgresUser: getEnv("POSTGRES_USER", "postgres"),
			PostgresPassword: getEnv("POSTGRES_PASSWORD", "postgres"),
			PostgresDBName: getEnv("POSTGRES_DBNAME", "golang"),
		}
	})

	return config
}

func loadEnv() error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("Error loading .env file: %v", err)
	}

	return nil

}

// getEnvAsBool gets an environment variable as a boolean or returns a default value.
func getEnvAsBool(key string, defaultValue bool) bool {
	valStr := getEnv(key, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}
	return defaultValue
}

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

