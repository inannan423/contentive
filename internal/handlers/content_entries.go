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

	// Create new entry
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

func GetContentEntries(c *fiber.Ctx) error {
	// Parse content type ID from URL
	contentTypeID, err := uuid.Parse(c.Params("contentTypeId"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid content type ID",
		})
	}

	// Fetch all entries for the content type
	var entries []models.ContentEntry
	if err := config.DB.Where("content_type_id = ?", contentTypeID).Find(&entries).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to retrieve content entries",
		})
	}

	return c.Status(200).JSON(entries)
}

func GetContentEntry(c *fiber.Ctx) error {
	// Parse entry ID from URL
	entryID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid content entry ID",
		})
	}

	// Fetch the entry
	var entry models.ContentEntry
	if err := config.DB.Where("id = ?", entryID).First(&entry).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to retrieve content entry",
		})
	}

	return c.Status(200).JSON(entry)
}

func UpdateContentEntry(c *fiber.Ctx) error {
	// Get validated data from context
	contentTypeID := c.Locals("contentTypeID").(uuid.UUID)
	entryID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid content entry ID",
		})
	}

	// Fetch existing entry
	var entry models.ContentEntry
	if err := config.DB.Where("id = ?", entryID).First(&entry).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to retrieve content entry",
		})
	}

	// Verify entry belongs to the specified content type
	if entry.ContentTypeID != contentTypeID {
		return c.Status(400).JSON(fiber.Map{
			"error": "Content entry does not belong to the specified content type",
		})
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
	// Parse entry ID from URL
	entryID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid content entry ID",
		})
	}

	// Fetch the entry
	var entry models.ContentEntry
	if err := config.DB.Where("id = ?", entryID).First(&entry).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to retrieve content entry",
		})
	}

	// Delete the entry
	if err := config.DB.Delete(&entry).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to delete content entry",
		})
	}

	return c.Status(204).SendString("Content entry deleted successfully")
}
