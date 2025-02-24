package middlewares

import (
	"contentive/config"
	"contentive/internal/models"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

func ValidateContentEntry() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 通过 slug 获取 content type
		identifier := c.Params("identifier")
		var contentType models.ContentType
		if err := config.DB.Preload("Fields").First(&contentType, "slug = ?", identifier).Error; err != nil {
			return c.Status(404).JSON(fiber.Map{
				"error": "Content type not found",
			})
		}

		// Parse request body
		var data map[string]interface{}
		if err := c.BodyParser(&data); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		// Check Single type constraint
		if contentType.Type == string(models.Single) {
			var count int64
			if err := config.DB.Model(&models.ContentEntry{}).Where("content_type_id = ?", contentType.ID).Count(&count).Error; err != nil {
				return c.Status(500).JSON(fiber.Map{
					"error": "Failed to check existing entries",
				})
			}

			if c.Method() == "POST" && count > 0 {
				return c.Status(400).JSON(fiber.Map{
					"error": "Single type content can only have one entry",
				})
			}
		}

		// Create a map for field validation
		fieldTypes := make(map[string]models.FieldTypeEnum)
		for _, field := range contentType.Fields {
			fieldTypes[field.Label] = field.Type
		}

		// Validate provided fields
		for fieldName, value := range data {
			if fieldName == "slug" {
				continue // Skip validation for slug field
			}

			fieldType, exists := fieldTypes[fieldName]
			if !exists {
				return c.Status(400).JSON(fiber.Map{
					"error": fmt.Sprintf("Field '%s' is not defined in content type", fieldName),
				})
			}

			// Skip validation for null values in PUT requests
			if c.Method() == "PUT" && value == nil {
				continue
			}

			// Validate field type
			switch fieldType {
			case models.Text, models.RichText:
				if _, ok := value.(string); !ok {
					return c.Status(400).JSON(fiber.Map{
						"error": fmt.Sprintf("Field '%s' must be a string", fieldName),
					})
				}
			case models.Number:
				switch value.(type) {
				case float64:
					// JSON numbers are decoded as float64
					break
				default:
					return c.Status(400).JSON(fiber.Map{
						"error": fmt.Sprintf("Field '%s' must be a number", fieldName),
					})
				}
			case models.Boolean:
				if _, ok := value.(bool); !ok {
					return c.Status(400).JSON(fiber.Map{
						"error": fmt.Sprintf("Field '%s' must be a boolean", fieldName),
					})
				}
			case models.Date:
				dateStr, ok := value.(string)
				if !ok {
					return c.Status(400).JSON(fiber.Map{
						"error": fmt.Sprintf("Field '%s' must be a date string", fieldName),
					})
				}
				if _, err := time.Parse(time.RFC3339, dateStr); err != nil {
					return c.Status(400).JSON(fiber.Map{
						"error": fmt.Sprintf("Field '%s' must be a valid ISO 8601 date", fieldName),
					})
				}
			case models.Relation:
				idStr, ok := value.(string)
				if !ok {
					return c.Status(400).JSON(fiber.Map{
						"error": fmt.Sprintf("Field '%s' must be a UUID string", fieldName),
					})
				}
				if _, err := uuid.Parse(idStr); err != nil {
					return c.Status(400).JSON(fiber.Map{
						"error": fmt.Sprintf("Field '%s' must be a valid UUID", fieldName),
					})
				}
			case models.Enum:
				if _, ok := value.(string); !ok {
					return c.Status(400).JSON(fiber.Map{
						"error": fmt.Sprintf("Field '%s' must be a string", fieldName),
					})
				}
			}
		}

		// For POST requests, ensure all required fields are present
		if c.Method() == "POST" {
			for _, field := range contentType.Fields {
				if _, exists := data[field.Label]; !exists {
					return c.Status(400).JSON(fiber.Map{
						"error": fmt.Sprintf("Field '%s' is required", field.Label),
					})
				}
			}
		}

		// Validate slug for POST requests
		if c.Method() == "POST" {
			slug, ok := data["slug"].(string)
			if !ok || slug == "" {
				return c.Status(400).JSON(fiber.Map{
					"error": "Slug is required",
				})
			}

			// Validate slug format
			if !isValidSlug(slug) {
				return c.Status(400).JSON(fiber.Map{
					"error": "Invalid slug format. Use only lowercase letters, numbers, and hyphens",
				})
			}

			// Check slug uniqueness within content type
			var existingEntry models.ContentEntry
			if err := config.DB.Where("content_type_id = ? AND slug = ?", contentType.ID, slug).First(&existingEntry).Error; err == nil {
				return c.Status(400).JSON(fiber.Map{
					"error": "An entry with this slug already exists",
				})
			}

			c.Locals("slug", slug)
			delete(data, "slug") // Remove slug from data to avoid storing it twice
		}

		// For PUT requests, validate slug if provided
		if c.Method() == "PUT" {
			if slug, ok := data["slug"].(string); ok {
				if slug == "" {
					return c.Status(400).JSON(fiber.Map{
						"error": "Slug cannot be empty",
					})
				}

				if !isValidSlug(slug) {
					return c.Status(400).JSON(fiber.Map{
						"error": "Invalid slug format. Use only lowercase letters, numbers, and hyphens",
					})
				}

				// Check slug uniqueness, excluding current entry
				currentSlug := c.Params("slug")
				if slug != currentSlug {
					var existingEntry models.ContentEntry
					if err := config.DB.Where("content_type_id = ? AND slug = ? AND slug != ?", contentType.ID, slug, currentSlug).First(&existingEntry).Error; err == nil {
						return c.Status(400).JSON(fiber.Map{
							"error": "An entry with this slug already exists",
						})
					}
				}

				c.Locals("slug", slug)
				delete(data, "slug") // Remove slug from data to avoid storing it twice
			}
		}

		jsonData, err := json.Marshal(data)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to process data",
			})
		}

		c.Locals("contentTypeID", contentType.ID)
		c.Locals("jsonData", datatypes.JSON(jsonData))
		return c.Next()
	}
}
