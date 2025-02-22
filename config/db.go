package config

import (
	"contentive/models"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dsn := "host=" + AppConfig.DBHost + " user=" + AppConfig.DBUser + " dbname=" + AppConfig.DBName + " password=" + AppConfig.DBPassword + " port=" + AppConfig.DBPort + " sslmode=disable"

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to the database:", err)
	}
	log.Println("Connected to the database!")

	if err := DB.AutoMigrate(&models.ContentType{}, &models.Field{}, &models.Content{}, &models.ContentItem{}); err != nil {
		log.Fatal("failed to migrate database:", err)
	}
	log.Println("Database migration completed!")
}
