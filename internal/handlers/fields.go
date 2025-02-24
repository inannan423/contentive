package handlers

import (
	"contentive/config"
	"contentive/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func AddField(c *fiber.Ctx) error {
	field := c.Locals("field").(models.Field)
	contentTypeID := c.Locals("contentTypeID").(uuid.UUID)

	field.ContentTypeID = contentTypeID

	if err := config.DB.Create(&field).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to create field",
		})
	}

	return c.Status(201).JSON(field)
}

func UpdateField(c *fiber.Ctx) error {
	field := c.Locals("field").(models.Field)
	fieldID := c.Locals("fieldID").(uuid.UUID)
	contentTypeID := c.Locals("contentTypeID").(uuid.UUID)

	var existingField models.Field
	if err := config.DB.Where("id = ? AND content_type_id = ?", fieldID, contentTypeID).First(&existingField).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Field not found or does not belong to this content type",
		})
	}

	if field.Type != "" && field.Type != existingField.Type {
		return c.Status(400).JSON(fiber.Map{
			"error": "Field type cannot be changed after creation",
		})
	}

	if err := config.DB.Model(&models.Field{}).Where("id = ?", fieldID).Update("label", field.Label).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to update field",
		})
	}

	var updatedField models.Field
	if err := config.DB.Where("id = ?", fieldID).First(&updatedField).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to retrieve updated field",
		})
	}

	return c.Status(200).JSON(updatedField)
}

func DeleteField(c *fiber.Ctx) error {
	identifier := c.Params("identifier")
	var contentType models.ContentType
	if err := config.DB.First(&contentType, "slug = ?", identifier).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Content type not found",
		})
	}

	fieldID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid field ID",
		})
	}

	var field models.Field
	if err := config.DB.Where("id = ? AND content_type_id = ?", fieldID, contentType.ID).First(&field).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Field not found or does not belong to this content type",
		})
	}

	if err := config.DB.Delete(&field).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to delete field",
		})
	}

	return c.SendStatus(204)
}
