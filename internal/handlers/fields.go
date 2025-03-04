package handlers

import (
	"contentive/config"
	"contentive/internal/logger"
	"contentive/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func AddField(c *fiber.Ctx) error {
	field := c.Locals("field").(models.Field)
	contentTypeID := c.Locals("contentTypeID").(uuid.UUID)

	field.ContentTypeID = contentTypeID

	if err := config.DB.Create(&field).Error; err != nil {
		logger.Error("Error creating field %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to create field",
		})
	}

	logger.Info("Field created successfully")

	return c.Status(201).JSON(field)
}

func UpdateField(c *fiber.Ctx) error {
	field := c.Locals("field").(models.Field)
	fieldID := c.Locals("fieldID").(uuid.UUID)
	contentTypeID := c.Locals("contentTypeID").(uuid.UUID)

	var existingField models.Field
	if err := config.DB.Where("id = ? AND content_type_id = ?", fieldID, contentTypeID).First(&existingField).Error; err != nil {
		logger.Error("Error fetching field %v", err)
		return c.Status(404).JSON(fiber.Map{
			"error": "Field not found or does not belong to this content type",
		})
	}

	if field.Type != "" && field.Type != existingField.Type {
		logger.Error("Field type cannot be changed after creation")
		return c.Status(400).JSON(fiber.Map{
			"error": "Field type cannot be changed after creation",
		})
	}

	updates := map[string]interface{}{
		"label":    field.Label,
		"required": field.Required,
	}

	if err := config.DB.Model(&models.Field{}).Where("id = ?", fieldID).Updates(updates).Error; err != nil {
		logger.Error("Error updating field %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to update field",
		})
	}

	var updatedField models.Field
	if err := config.DB.Where("id = ?", fieldID).First(&updatedField).Error; err != nil {
		logger.Error("Error fetching updated field %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to retrieve updated field",
		})
	}

	logger.Info("Field updated successfully")

	return c.Status(200).JSON(updatedField)
}

func DeleteField(c *fiber.Ctx) error {
	identifier := c.Params("identifier")
	var contentType models.ContentType
	if err := config.DB.First(&contentType, "slug = ?", identifier).Error; err != nil {
		logger.Error("Error fetching content type %v", err)
		return c.Status(404).JSON(fiber.Map{
			"error": "Content type not found",
		})
	}

	fieldID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		logger.Error("Error parsing field ID %v", err)
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid field ID",
		})
	}

	var field models.Field
	if err := config.DB.Where("id = ? AND content_type_id = ?", fieldID, contentType.ID).First(&field).Error; err != nil {
		logger.Error("Error fetching field %v", err)
		return c.Status(404).JSON(fiber.Map{
			"error": "Field not found or does not belong to this content type",
		})
	}

	if err := config.DB.Delete(&field).Error; err != nil {
		logger.Error("Error deleting field %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to delete field",
		})
	}

	logger.Info("Field deleted successfully")

	return c.SendStatus(204)
}
