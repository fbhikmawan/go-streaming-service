package config

import (
	"fmt"
	"os"
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
}


// el singleton de configuracion 
var (
	config     *Config
	configOnce sync.Once
)


func GetConfig() *Config {

	// usar Sync.Once para garantizar que la configuración se cargue solo una vez y evitar problemas de rendimiento
	// y usa singleton para garantizar que solo haya una instancia de la configuración en toda la aplicación.

	configOnce.Do(func() {
		err := loadEnv()
		if err != nil {
			panic(fmt.Sprintf("Error al cargar el archivo .env: %v", err))
		}

		config = &Config{
			Port:         getEnv("PORT", "8080"),
			DatabaseURL:  getEnv("DATABASE_URL", "postgres://localhost:5432/mydb"),
			JWTSecretKey: getEnv("JWT_SECRET_KEY", "secretJwtKey"),
			AWSRegion:    getEnv("AWS_REGION", ""),
			AWSBucketName: getEnv("AWS_BUCKET_NAME", ""),
			AWSAccessKey: getEnv("AWS_ACCESS_KEY_ID", ""),
			AWSSecretKey: getEnv("AWS_SECRET_ACCESS_KEY", ""),
		}
	})

	return config
}

func loadEnv() error {
	err := godotenv.Load()
	if err != nil {
		return fmt.Errorf("error loading .env: %v", err)
	}

	return nil

}


func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

