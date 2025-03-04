package handlers

import (
	"contentive/config"
	"contentive/internal/logger"
	"contentive/internal/models"
	"encoding/json"
	"fmt"
	"math"

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
		logger.Error("Error creating content entry: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to create content entry",
		})
	}

	return c.Status(201).JSON(entry)
}

func GetContentEntries(c *fiber.Ctx) error {
	// Get pagination parameters from query
	page := c.QueryInt("page", 1)
	pageSize := c.QueryInt("pageSize", 10)
	sortBy := c.Query("sortBy", "created_at") // Default sort by created_at
	sortOrder := c.Query("sortOrder", "desc") // Default sort order is descending
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	// Validate sort order
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "desc"
	}

	// Get content type by identifier with fields
	identifier := c.Params("identifier")
	var contentType models.ContentType
	if err := config.DB.Preload("Fields").First(&contentType, "slug = ?", identifier).Error; err != nil {
		logger.Error("Error fetching content type: %v", err)
		return c.Status(404).JSON(fiber.Map{
			"error": "Content type not found",
		})
	}

	// Count total entries
	var total int64
	if err := config.DB.Model(&models.ContentEntry{}).Where("content_type_id = ?", contentType.ID).Count(&total).Error; err != nil {
		logger.Error("Error counting entries: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to count entries",
		})
	}

	// Build order clause
	orderClause := fmt.Sprintf("%s %s", sortBy, sortOrder)

	// Get paginated entries
	var entries []models.ContentEntry
	if err := config.DB.Where("content_type_id = ?", contentType.ID).
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Order(orderClause).
		Find(&entries).Error; err != nil {
		logger.Error("Error fetching content entries: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to retrieve content entries",
		})
	}

	// Calculate total pages
	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	logger.Info("Fetched %d entries for content type %s", len(entries), contentType.Slug)

	// Return paginated response with sort info and content type
	return c.Status(200).JSON(fiber.Map{
		"contentType": contentType,
		"data":        entries,
		"pagination": fiber.Map{
			"current":  page,
			"pageSize": pageSize,
			"total":    total,
			"pages":    totalPages,
		},
		"sort": fiber.Map{
			"field": sortBy,
			"order": sortOrder,
		},
	})
}

func GetContentEntry(c *fiber.Ctx) error {
	identifier := c.Params("identifier")
	var contentType models.ContentType
	if err := config.DB.First(&contentType, "slug = ?", identifier).Error; err != nil {
		logger.Error("Error fetching content type: %v", err)
		return c.Status(404).JSON(fiber.Map{
			"error": "Content type not found",
		})
	}

	slug := c.Params("slug")

	// Fetch the entry by slug
	var entry models.ContentEntry
	if err := config.DB.Where("content_type_id = ? AND slug = ?", contentType.ID, slug).First(&entry).Error; err != nil {
		logger.Error("Error fetching content entry: %v", err)
		return c.Status(404).JSON(fiber.Map{
			"error": "Content entry not found",
		})
	}

	logger.Info("Fetched content entry for content type %s", contentType.Slug)

	return c.Status(200).JSON(entry)
}

func UpdateContentEntry(c *fiber.Ctx) error {
	// Get validated data from context
	contentTypeID := c.Locals("contentTypeID").(uuid.UUID)
	currentSlug := c.Params("slug")

	// Fetch existing entry by slug
	var entry models.ContentEntry
	if err := config.DB.Where("content_type_id = ? AND slug = ?", contentTypeID, currentSlug).First(&entry).Error; err != nil {
		logger.Error("Error fetching content entry: %v", err)
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
		logger.Error("Error parsing existing data: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to parse existing data",
		})
	}

	// Get update data
	newData := c.Locals("jsonData").(datatypes.JSON)
	var updateData map[string]interface{}
	if err := json.Unmarshal(newData, &updateData); err != nil {
		logger.Error("Error parsing update data: %v", err)
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
		logger.Error("Error processing merged data: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to process merged data",
		})
	}

	// Update entry with merged data
	entry.Data = datatypes.JSON(mergedData)
	if err := config.DB.Save(&entry).Error; err != nil {
		logger.Error("Error updating content entry: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to update content entry",
		})
	}

	logger.Info("Updated content entry for content type %s", contentTypeID)

	return c.Status(200).JSON(entry)
}

func DeleteContentEntry(c *fiber.Ctx) error {
	identifier := c.Params("identifier")
	var contentType models.ContentType
	if err := config.DB.First(&contentType, "slug = ?", identifier).Error; err != nil {
		logger.Error("Error fetching content type: %v", err)
		return c.Status(404).JSON(fiber.Map{
			"error": "Content type not found",
		})
	}

	slug := c.Params("slug")

	var entry models.ContentEntry
	if err := config.DB.Where("content_type_id = ? AND slug = ?", contentType.ID, slug).First(&entry).Error; err != nil {
		logger.Error("Error fetching content entry: %v", err)
		return c.Status(404).JSON(fiber.Map{
			"error": "Content entry not found",
		})
	}

	if err := config.DB.Delete(&entry).Error; err != nil {
		logger.Error("Error deleting content entry: %v", err)
		return c.Status(500).JSON(fiber.Map{
			"error": "Failed to delete content entry",
		})
	}

	logger.Info("Deleted content entry for content type %s", contentType.Slug)

	return c.Status(204).SendString("Content entry deleted successfully")
}
