package config

import (
	"contentive/internal/models"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dsn := "host=" + AppConfig.DBHost + " user=" + AppConfig.DBUser + " dbname=" + AppConfig.DBName + " password=" + AppConfig.DBPassword + " port=" + AppConfig.DBPort + " sslmode=disable"

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		log.Fatal("failed to connect to the database:", err)
	}
	log.Println("Connected to the database!")

	if err := DB.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`).Error; err != nil {
		log.Fatal("failed to create uuid-ossp extension:", err)
	}
	log.Println("UUID extension enabled!")

	// if err := DB.Migrator().DropTable(&models.ContentEntry{}, &models.Field{}, &models.ContentType{}); err != nil {
	// 	log.Fatal("failed to drop tables:", err)
	// }

	if err := DB.AutoMigrate(&models.ContentType{}, &models.Field{}, &models.ContentEntry{}); err != nil {
		log.Fatal("failed to migrate database:", err)
	}
	log.Println("Database migration completed!")
}
