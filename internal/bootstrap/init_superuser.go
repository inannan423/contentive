package bootstrap

import (
	"contentive/internal/config"
	"contentive/internal/database"
	"contentive/internal/logger"
	"contentive/internal/models"
	"log"
	"time"

	"github.com/lib/pq"
)

// InitSuperUser initializes the super admin user from the environment variables.
func InitSuperUser() {
	username := config.AppConfig.SUPER_USER_NAME
	email := config.AppConfig.SUPER_USER_EMAIL
	password := config.AppConfig.SUPER_USER_PASSWORD

	// Check if super admin user already exists
	var count int64
	if err := database.DB.Model(&models.AdminUser{}).Where("Name = ?", username).Count(&count).Error; err != nil {
		logger.Error("Error checking super admin user: %v", err)
		log.Fatal("Error checking super admin user: ", err)
	}

	if count > 0 {
		logger.GeneralAction("Super admin user already exists")
		return
	}

	// Create super admin user
	superAdmin := models.AdminUser{
		Name:        username,
		Email:       email,
		Password:    password, // will be hashed by BeforeCreate hook
		Role:        models.AdminUserRoleSuperAdmin,
		Status:      models.AdminUserStatusActive,
		LastLoginAt: time.Now(),
	}

	if err := database.DB.Create(&superAdmin).Error; err != nil {
		logger.Error("Error creating super admin user, error: %v", err)
		log.Fatal("Error creating super admin user: ", err)
	}

	// Update super admin permissions
	if err := database.DB.Model(&superAdmin).Update("permissions", pq.Array([]string{"all"})).Error; err != nil {
		logger.Error("Error updating super admin permissions, error: %v", err)
		log.Fatal("Error updating super admin permissions: ", err)
	}

	logger.GeneralAction("Super admin user created")
}
