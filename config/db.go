package config

import (
	"fmt"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	dbInstance *gorm.DB
	once       sync.Once
)

// GetDsn genera la cadena de conexión para la base de datos.
func GetDsn() string {
	config := GetConfig()
	host := config.PostgresHost
	port := config.PostgresPort
	user := config.PostgresUser
	password := config.PostgresPassword
	dbname := config.PostgresDBName
	
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbname, port)
}

// GetDB devuelve una instancia única de la conexión a la base de datos.
func GetDB() (*gorm.DB, error) {
	var err error
	once.Do(func() {
		dsn := GetDsn()
		dbInstance, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	})
	if err != nil {
		return nil, err
	}
	return dbInstance, nil
}
