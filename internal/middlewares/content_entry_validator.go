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
		// Parse content type ID from URL parameter
		contentTypeID, err := uuid.Parse(c.Params("contentTypeId"))
		if err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid content type ID",
			})
		}

		// Load content type with its fields
		var contentType models.ContentType
		if err := config.DB.Preload("Fields").First(&contentType, "id = ?", contentTypeID).Error; err != nil {
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
			if err := config.DB.Model(&models.ContentEntry{}).Where("content_type_id = ?", contentTypeID).Count(&count).Error; err != nil {
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

		jsonData, err := json.Marshal(data)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "Failed to process data",
			})
		}

		c.Locals("contentTypeID", contentTypeID)
		c.Locals("jsonData", datatypes.JSON(jsonData))
		return c.Next()
	}
}
