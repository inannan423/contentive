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

	if err := config.DB.Model(&models.Field{}).Where("id = ?", fieldID).Updates(field).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to update field",
		})
	}

	return c.Status(200).JSON(field)
}

func DeleteField(c *fiber.Ctx) error {
	contentTypeID := c.Locals("contentTypeID").(uuid.UUID)
	fieldID := c.Locals("fieldID").(uuid.UUID)

	// 验证字段是否属于指定的内容类型
	var field models.Field
	if err := config.DB.Where("id = ? AND content_type_id = ?", fieldID, contentTypeID).First(&field).Error; err != nil {
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
