package handlers

import (
	"contentive/config"
	"contentive/internal/models"
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

func CreateContentEntry(c *fiber.Ctx) error {
	// Get validated data from context
	contentTypeID := c.Locals("contentTypeID").(uuid.UUID)
	jsonData := c.Locals("jsonData").(datatypes.JSON)
	slug := c.Locals("slug").(string)

	// Create new entry
	entry := models.ContentEntry{
		ContentTypeID: contentTypeID,
		Slug:          slug,
		Data:          jsonData,
	}

	if err := config.DB.Create(&entry).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to create content entry",
		})
	}

	return c.Status(201).JSON(entry)
}

func GetContentEntries(c *fiber.Ctx) error {
	identifier := c.Params("identifier")
	var contentType models.ContentType
	if err := config.DB.First(&contentType, "slug = ?", identifier).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Content type not found",
		})
	}

	var entries []models.ContentEntry
	if err := config.DB.Where("content_type_id = ?", contentType.ID).Find(&entries).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to retrieve content entries",
		})
	}

	return c.Status(200).JSON(entries)
}

func GetContentEntry(c *fiber.Ctx) error {
	identifier := c.Params("identifier")
	var contentType models.ContentType
	if err := config.DB.First(&contentType, "slug = ?", identifier).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Content type not found",
		})
	}

	slug := c.Params("slug")

	// Fetch the entry by slug
	var entry models.ContentEntry
	if err := config.DB.Where("content_type_id = ? AND slug = ?", contentType.ID, slug).First(&entry).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Content entry not found",
		})
	}

	return c.Status(200).JSON(entry)
}

func UpdateContentEntry(c *fiber.Ctx) error {
	// Get validated data from context
	contentTypeID := c.Locals("contentTypeID").(uuid.UUID)
	currentSlug := c.Params("slug")

	// Fetch existing entry by slug
	var entry models.ContentEntry
	if err := config.DB.Where("content_type_id = ? AND slug = ?", contentTypeID, currentSlug).First(&entry).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Content entry not found",
		})
	}

	if newSlug, ok := c.Locals("slug").(string); ok && newSlug != "" && newSlug != currentSlug {
		entry.Slug = newSlug
	}

	// Parse existing data
	var existingData map[string]interface{}
	if err := json.Unmarshal(entry.Data, &existingData); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to parse existing data",
		})
	}

	// Get update data
	newData := c.Locals("jsonData").(datatypes.JSON)
	var updateData map[string]interface{}
	if err := json.Unmarshal(newData, &updateData); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to parse update data",
		})
	}

	// Merge data: update only provided fields
	for k, v := range updateData {
		existingData[k] = v
	}

	// Convert merged data back to JSON
	mergedData, err := json.Marshal(existingData)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to process merged data",
		})
	}

	// Update entry with merged data
	entry.Data = datatypes.JSON(mergedData)
	if err := config.DB.Save(&entry).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to update content entry",
		})
	}

	return c.Status(200).JSON(entry)
}

func DeleteContentEntry(c *fiber.Ctx) error {
	identifier := c.Params("identifier")
	var contentType models.ContentType
	if err := config.DB.First(&contentType, "slug = ?", identifier).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Content type not found",
		})
	}

	slug := c.Params("slug")

	var entry models.ContentEntry
	if err := config.DB.Where("content_type_id = ? AND slug = ?", contentType.ID, slug).First(&entry).Error; err != nil {
		return c.Status(404).JSON(fiber.Map{
			"error": "Content entry not found",
		})
	}

	if err := config.DB.Delete(&entry).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to delete content entry",
		})
	}

	return c.Status(204).SendString("Content entry deleted successfully")
}
