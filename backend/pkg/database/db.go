package database

import (
	"log"

	"webman/pkg/models"

	"github.com/google/uuid"

	_ "github.com/mattn/go-sqlite3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(dbPath string) error {
	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return err
	}

	// Auto migrate the schema
	err = DB.AutoMigrate(
		&models.Collection{},
		&models.Request{},
		&models.DefaultHeader{},
	)
	if err != nil {
		return err
	}

	// Initialize default headers if they don't exist
	var count int64
	DB.Model(&models.DefaultHeader{}).Count(&count)
	if count == 0 {
		headers := []models.DefaultHeader{
			{ID: uuid.New().String(), Name: "Accept", Value: "application/json", Description: "Indicates that the client expects JSON response"},
			{ID: uuid.New().String(), Name: "Content-Type", Value: "application/json", Description: "Indicates that the request body is in JSON format"},
			{ID: uuid.New().String(), Name: "Authorization", Value: "Bearer ", Description: "Bearer token authentication"},
			{ID: uuid.New().String(), Name: "Cache-Control", Value: "no-cache", Description: "Controls caching behavior"},
			{ID: uuid.New().String(), Name: "User-Agent", Value: "Webman/1.0", Description: "Identifies the client application"},
			{ID: uuid.New().String(), Name: "Accept-Language", Value: "en-US,en;q=0.9", Description: "Preferred language for response"},
			{ID: uuid.New().String(), Name: "X-Requested-With", Value: "XMLHttpRequest", Description: "Indicates an AJAX request"},
		}
		if err := DB.Create(&headers).Error; err != nil {
			log.Printf("Failed to create default headers: %v", err)
		}
	}

	return nil
}
