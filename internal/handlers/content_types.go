package handlers

import (
	"contentive/config"
	"contentive/internal/models"

	"github.com/gofiber/fiber/v2"
)

func CreateContentType(c *fiber.Ctx) error {
	contentType := c.Locals("contentType").(models.ContentType)

	if err := config.DB.Create(&contentType).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to create content type",
		})
	}

	return c.Status(201).JSON(contentType)
}

func GetAllContentTypes(c *fiber.Ctx) error {
	var contentTypes []models.ContentType
	if err := config.DB.Find(&contentTypes).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to retrieve content types",
		})
	}

	// Preload fields
	for i, contentType := range contentTypes {
		if err := config.DB.Preload("Fields").First(&contentTypes[i], contentType.ID).Error; err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to retrieve content types",
			})
		}
	}

	return c.JSON(contentTypes)
}
