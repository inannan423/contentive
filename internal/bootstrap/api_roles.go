package bootstrap

import (
	"contentive/config"
	"contentive/internal/models"
	"contentive/internal/utils"
	"log"

	"github.com/google/uuid"
)

func InitAPIRoles() {
	roles := []models.APIRole{
		{
			Name:        "Public User",
			Type:        models.PublicUser,
			Description: "Public Access",
			IsSystem:    true,
		},
		{
			Name:        "Authenticated User",
			Type:        models.AuthenticatedUser,
			Description: "Authenticated Access",
			IsSystem:    true,
		},
	}

	for _, r := range roles {
		var existingRole models.APIRole
		if err := config.DB.Where(models.APIRole{Type: r.Type}).First(&existingRole).Error; err != nil {
			if r.Type == models.AuthenticatedUser {
				apiKey, err := utils.GenerateAPIKey()
				if err != nil {
					log.Printf("Error generating API key for %s: %v", r.Name, err)
					continue
				}
				r.APIKey = apiKey
			}

			if err := config.DB.Create(&r).Error; err != nil {
				log.Printf("Error creating API role %s: %v", r.Name, err)
				continue
			}
			log.Printf("Created API role: %s", r.Name)
		} else {
			existingRole.Name = r.Name
			existingRole.Description = r.Description
			if err := config.DB.Save(&existingRole).Error; err != nil {
				log.Printf("Error updating API role %s: %v", r.Name, err)
			}
		}
	}
}

func InitDefaultAPIPermissions() {
	var contentTypes []models.ContentType
	if err := config.DB.Find(&contentTypes).Error; err != nil {
		log.Printf("Error fetching content types: %v", err)
		return
	}

	var publicRole, authRole models.APIRole
	if err := config.DB.Where("type = ?", models.PublicUser).First(&publicRole).Error; err != nil {
		log.Printf("Error fetching public role: %v", err)
		return
	}
	if err := config.DB.Where("type = ?", models.AuthenticatedUser).First(&authRole).Error; err != nil {
		log.Printf("Error fetching authenticated role: %v", err)
		return
	}

	for _, ct := range contentTypes {
		createPublicPermission(publicRole.ID, ct.ID, models.ReadOperation, true)

		operations := []models.OperationType{models.CreateOperation, models.ReadOperation, models.UpdateOperation, models.DeleteOperation}
		for _, op := range operations {
			createPublicPermission(authRole.ID, ct.ID, op, true)
		}
	}
}

func createPublicPermission(apiRoleID uuid.UUID, contentTypeID uuid.UUID, operation models.OperationType, enabled bool) {
	var permission models.APIPermission
	result := config.DB.Where("api_role_id = ? AND content_type_id = ? AND operation = ?",
		apiRoleID, contentTypeID, operation).First(&permission)

	if result.Error != nil {
		permission = models.APIPermission{
			APIRoleID:     apiRoleID,
			ContentTypeID: contentTypeID,
			Operation:     operation,
			Enabled:       enabled,
		}
		if err := config.DB.Create(&permission).Error; err != nil {
			log.Printf("Error creating API permission: %v", err)
		}
	} else {
		permission.Enabled = enabled
		if err := config.DB.Save(&permission).Error; err != nil {
			log.Printf("Error updating API permission: %v", err)
		}
	}
}
