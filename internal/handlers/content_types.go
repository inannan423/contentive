package handlers

import (
	"contentive/config"
	"contentive/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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

func GetContentType(c *fiber.Ctx) error {
	contentTypeID, err := uuid.Parse(c.Params("contentTypeId"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid content type ID",
		})
	}

	var contentType models.ContentType
	if err := config.DB.Preload("Fields").First(&contentType, "id = ?", contentTypeID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Content type not found",
		})
	}
	return c.JSON(contentType)
}

func UpdateContentType(c *fiber.Ctx) error {
	// Get validated content type from middleware
	updatedType := c.Locals("contentType").(models.ContentType)

	// Parse content type ID from URL
	contentTypeID, err := uuid.Parse(c.Params("contentTypeId"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid content type ID",
		})
	}

	// Fetch existing content type
	var existingType models.ContentType
	if err := config.DB.First(&existingType, "id = ?", contentTypeID).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Content type not found",
		})
	}

	// Update only name field
	if updatedType.Name != "" {
		existingType.Name = updatedType.Name
	}

	// Save updates
	if err := config.DB.Save(&existingType).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to update content type",
		})
	}

	// Reload content type with fields
	if err := config.DB.Preload("Fields").First(&existingType, "id = ?", contentTypeID).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to retrieve updated content type",
		})
	}

	return c.Status(200).JSON(existingType)
}
