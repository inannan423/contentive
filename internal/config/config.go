package config

import (
	"contentive/internal/logger"
	"contentive/internal/models"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUser                string
	DBPassword            string
	DBName                string
	DBHost                string
	DBPort                string
	JWTSecret             string
	SUPER_USER_NAME       string
	SUPER_USER_PASSWORD   string
	SUPER_USER_EMAIL      string
	MEDIA_STORAGE_TYPE    string
	MEDIA_STORAGE_PATH    string // for local storage
	MEDIA_STORAGE_URL     string // for aliyun oss storage
	OSS_REGION_ID         string
	OSS_ACCESS_KEY_ID     string
	OSS_ACCESS_KEY_SECRET string
	OSS_BUCKET_NAME       string
	LLM_PROVIDER          string
	LLM_BASE_URL          string
	LLM_API_KEY           string
	LLM_MODEL             string
	LLM_MAX_TOKENS        int
	LLM_TEMPERATURE       float64
	LLM_TOP_P             float64
}

var AppConfig Config

func InitConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file, please check if .env file exists")
	}

	AppConfig = Config{
		DBUser:                os.Getenv("DB_USER"),
		DBPassword:            os.Getenv("DB_PASSWORD"),
		DBName:                os.Getenv("DB_NAME"),
		DBHost:                os.Getenv("DB_HOST"),
		DBPort:                os.Getenv("DB_PORT"),
		JWTSecret:             os.Getenv("JWT_SECRET"),
		SUPER_USER_NAME:       os.Getenv("SUPER_USER_NAME"),
		SUPER_USER_PASSWORD:   os.Getenv("SUPER_USER_PASSWORD"),
		SUPER_USER_EMAIL:      os.Getenv("SUPER_USER_EMAIL"),
		MEDIA_STORAGE_TYPE:    os.Getenv("MEDIA_STORAGE_TYPE"),
		MEDIA_STORAGE_PATH:    os.Getenv("MEDIA_STORAGE_PATH"),
		MEDIA_STORAGE_URL:     os.Getenv("MEDIA_STORAGE_URL"),
		OSS_REGION_ID:         os.Getenv("OSS_REGION_ID"),
		OSS_ACCESS_KEY_ID:     os.Getenv("OSS_ACCESS_KEY_ID"),
		OSS_ACCESS_KEY_SECRET: os.Getenv("OSS_ACCESS_KEY_SECRET"),
		OSS_BUCKET_NAME:       os.Getenv("OSS_BUCKET_NAME"),
		LLM_PROVIDER:          os.Getenv("LLM_PROVIDER"),
		LLM_BASE_URL:          os.Getenv("LLM_BASE_URL"),
		LLM_API_KEY:           os.Getenv("LLM_API_KEY"),
		LLM_MODEL:             os.Getenv("LLM_MODEL"),
		LLM_MAX_TOKENS:        getEnvAsInt("LLM_MAX_TOKENS", 2048),   // default value for max_tokens is 2048, you can change it to your own requiremen
		LLM_TEMPERATURE:       getEnvAsFloat("LLM_TEMPERATURE", 0.7), // default value for temperature is 0.7
		LLM_TOP_P:             getEnvAsFloat("LLM_TOP_P", 1),         // default value for top_p is 1
	}

	models.SetSecret(AppConfig.JWTSecret)

	logger.Info("Configuration loaded successfully!")
}

func getEnvAsInt(name string, defaultVal int) int {
	valueStr := os.Getenv(name)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultVal
}

func getEnvAsFloat(name string, defaultVal float64) float64 {
	valueStr := os.Getenv(name)
	if value, err := strconv.ParseFloat(valueStr, 64); err == nil {
		return value
	}
	return defaultVal
}
