package config

import (
	"contentive/models"
	"log"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var DB *gorm.DB

func InitDB() {
	dbURL := "host=" + AppConfig.DBHost + " user=" + AppConfig.DBUser + " dbname=" + AppConfig.DBName + " password=" + AppConfig.DBPassword + " port=" + AppConfig.DBPort + " sslmode=disable"

	var err error
	DB, err = gorm.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("failed to connect to the database:", err)
	}
	log.Println("Connected to the database!")

	DB.AutoMigrate(&models.ContentType{})
}
