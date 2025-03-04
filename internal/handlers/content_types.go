package handlers

import (
	"contentive/config"
	"contentive/internal/logger"
	"contentive/internal/models"

	"github.com/gofiber/fiber/v2"
)

func CreateContentType(c *fiber.Ctx) error {
	contentType := c.Locals("contentType").(models.ContentType)

	if err := config.DB.Create(&contentType).Error; err != nil {
		logger.Error("Error creating content type: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to create content type",
		})
	}

	logger.Info("Created content type: %v", contentType)

	return c.Status(201).JSON(contentType)
}

func GetAllContentTypes(c *fiber.Ctx) error {
	var contentTypes []models.ContentType
	if err := config.DB.Find(&contentTypes).Error; err != nil {
		logger.Error("Error fetching content types: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to retrieve content types",
		})
	}

	// Preload fields
	for i, contentType := range contentTypes {
		if err := config.DB.Preload("Fields").First(&contentTypes[i], contentType.ID).Error; err != nil {
			logger.Error("Error fetching content types: %v", err)
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to retrieve content types",
			})
		}
	}

	logger.Info("Fetched content types: %v", contentTypes)

	return c.JSON(contentTypes)
}

func GetContentType(c *fiber.Ctx) error {
	identifier := c.Params("identifier")

	var contentType models.ContentType
	query := config.DB.Preload("Fields")

	// Find by slug
	if err := query.First(&contentType, "slug = ?", identifier).Error; err != nil {
		logger.Error("Error fetching content type: %v", err)
		return c.Status(404).JSON(fiber.Map{
			"error": "Content type not found",
		})
	}

	logger.Info("Fetched content type: %v", contentType)

	return c.JSON(contentType)
}

func UpdateContentType(c *fiber.Ctx) error {
	// Get validated content type from middleware
	updatedType := c.Locals("contentType").(models.ContentType)

	// Parse content type slug from URL
	identifier := c.Params("identifier")
	var existingType models.ContentType
	query := config.DB

	// Find by slug
	if err := query.First(&existingType, "slug = ?", identifier).Error; err != nil {
		logger.Error("Error fetching content type: %v", err)
		return c.Status(404).JSON(fiber.Map{
			"error": "Content type not found",
		})
	}

	// Update fields
	if updatedType.Name != "" {
		existingType.Name = updatedType.Name
	}
	if updatedType.Slug != "" {
		existingType.Slug = updatedType.Slug
	}

	// Save updates
	if err := config.DB.Save(&existingType).Error; err != nil {
		logger.Error("Error updating content type: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to update content type",
		})
	}

	// Reload content type with fields
	if err := config.DB.Preload("Fields").First(&existingType, "id = ?", existingType.ID).Error; err != nil {
		logger.Error("Error fetching updated content type: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to retrieve updated content type",
		})
	}

	logger.Info("Updated content type: %v", existingType)

	return c.Status(200).JSON(existingType)
}
