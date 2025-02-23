package handlers

import (
	"contentive/config"
	"contentive/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

func CreateContentEntry(c *fiber.Ctx) error {
	contentTypeID := c.Locals("contentTypeID").(uuid.UUID)
	jsonData := c.Locals("jsonData").(datatypes.JSON)

	entry := models.ContentEntry{
		ContentTypeID: contentTypeID,
		Data:          jsonData,
	}

	if err := config.DB.Create(&entry).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to create content entry",
		})
	}

	return c.Status(201).JSON(entry)
}
