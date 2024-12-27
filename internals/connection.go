package internals

import (
	"fmt"
	"os"
	"shopping-site/pkg/loggers"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitiatePgConnection() *gorm.DB {
	dsn := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))

	client, err := gorm.Open(postgres.Open(dsn), &gorm.Config{TranslateError: true})
	if err != nil {
		loggers.InfoLog.Fatalf("Failed to open postgres client %v", err)
	}

	loggers.InfoLog.Print("Connected to postgress client")
	fmt.Println("Connected to postgress client")

	return client
}
