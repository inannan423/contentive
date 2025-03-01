package handlers

import (
	"contentive/config"
	"contentive/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func GetAPIRolePermissions(c *fiber.Ctx) error {
	roleID := c.Params("id")

	var role models.APIRole
	if err := config.DB.First(&role, "id = ?", roleID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Role not found"})
	}

	var permissions []models.APIPermission
	if err := config.DB.Preload("ContentType").Where("api_role_id = ?", roleID).Find(&permissions).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to fetch permissions"})
	}

	return c.Status(fiber.StatusOK).JSON(permissions)
}

func UpdateAPIRolePermission(c *fiber.Ctx) error {
	roleID := c.Params("id")

	var role models.APIRole
	if err := config.DB.First(&role, "id = ?", roleID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "API Role not found",
		})
	}

	var input struct {
		ContentTypeID string `json:"content_type_id"`
		Operation     string `json:"operation"`
		Enabled       bool   `json:"enabled"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	contentTypeID, err := uuid.Parse(input.ContentTypeID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid content type ID",
		})
	}

	var contentType models.ContentType
	if err := config.DB.First(&contentType, "id = ?", contentTypeID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Content type not found",
		})
	}

	operation := models.OperationType(input.Operation)
	if operation != models.CreateOperation && operation != models.ReadOperation && operation != models.UpdateOperation && operation != models.DeleteOperation {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid operation",
		})
	}

	var permission models.APIPermission
	result := config.DB.Where("api_role_id = ? AND content_type_id = ? AND operation = ?",
		role.ID, contentTypeID, operation).First(&permission)

	if result.Error != nil {
		permission = models.APIPermission{
			APIRoleID:     role.ID,
			ContentTypeID: contentTypeID,
			Operation:     operation,
			Enabled:       input.Enabled,
		}
		if err := config.DB.Create(&permission).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to create permission",
			})
		}
	} else {
		permission.Enabled = input.Enabled
		if err := config.DB.Save(&permission).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update permission",
			})
		}
	}

	if err := config.DB.Preload("ContentType").First(&permission, permission.ID).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch permission",
		})
	}

	return c.Status(fiber.StatusOK).JSON(permission)
}

func BatchUpdateAPIRolePermissions(c *fiber.Ctx) error {
	roleID := c.Params("id")

	var role models.APIRole
	if err := config.DB.First(&role, "id = ?", roleID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "API Role not found",
		})
	}

	var input []struct {
		ContentTypeID string `json:"content_type_id"`
		Operation     string `json:"operation"`
		Enabled       bool   `json:"enabled"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	updatedPermissions := []models.APIPermission{}

	for _, item := range input {
		contentTypeID, err := uuid.Parse(item.ContentTypeID)
		if err != nil {
			continue
		}

		var contentType models.ContentType
		if err := config.DB.First(&contentType, "id = ?", contentTypeID).Error; err != nil {
			continue
		}

		operation := models.OperationType(item.Operation)
		if operation != models.CreateOperation && operation != models.ReadOperation && operation != models.UpdateOperation && operation != models.DeleteOperation {
			continue
		}

		var permission models.APIPermission
		result := config.DB.Where("api_role_id = ? AND content_type_id = ? AND operation = ?",
			role.ID, contentTypeID, operation).First(&permission)

		if result.Error != nil {
			permission = models.APIPermission{
				APIRoleID:     role.ID,
				ContentTypeID: contentTypeID,
				Operation:     operation,
				Enabled:       item.Enabled,
			}
			if err := config.DB.Create(&permission).Error; err != nil {
				continue
			}
		} else {
			permission.Enabled = item.Enabled
			if err := config.DB.Save(&permission).Error; err != nil {
				continue
			}
		}

		if err := config.DB.Preload("ContentType").First(&permission, permission.ID).Error; err == nil {
			updatedPermissions = append(updatedPermissions, permission)
		}
	}

	return c.Status(fiber.StatusOK).JSON(updatedPermissions)
}
