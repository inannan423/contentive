package config

import (
	"contentive/internal/logger"
	"contentive/internal/models"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUser              string
	DBPassword          string
	DBName              string
	DBHost              string
	DBPort              string
	JWTSecret           string
	SUPER_USER_NAME     string
	SUPER_USER_PASSWORD string
	SUPER_USER_EMAIL    string
}

var AppConfig Config

func InitConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file, please check if .env file exists")
	}

	AppConfig = Config{
		DBUser:              os.Getenv("DB_USER"),
		DBPassword:          os.Getenv("DB_PASSWORD"),
		DBName:              os.Getenv("DB_NAME"),
		DBHost:              os.Getenv("DB_HOST"),
		DBPort:              os.Getenv("DB_PORT"),
		JWTSecret:           os.Getenv("JWT_SECRET"),
		SUPER_USER_NAME:     os.Getenv("SUPER_USER_NAME"),
		SUPER_USER_PASSWORD: os.Getenv("SUPER_USER_PASSWORD"),
		SUPER_USER_EMAIL:    os.Getenv("SUPER_USER_EMAIL"),
	}

	models.SetSecret(AppConfig.JWTSecret)

	logger.Info("Configuration loaded successfully!")
}
