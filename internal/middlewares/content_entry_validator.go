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

		var data map[string]interface{}
		if err := c.BodyParser(&data); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		// Validate fields
		for _, field := range contentType.Fields {
			value, exists := data[field.Label]
			if !exists {
				return c.Status(400).JSON(fiber.Map{
					"error": fmt.Sprintf("Field '%s' is required", field.Label),
				})
			}

			// Validate field type according to its type
			switch field.Type {
			case models.Text, models.RichText:
				if _, ok := value.(string); !ok {
					return c.Status(400).JSON(fiber.Map{
						"error": fmt.Sprintf("Field '%s' must be a string", field.Label),
					})
				}
			case models.Number:
				switch value.(type) {
				case float64:
					// JSON numbers are decoded as float64
					break
				default:
					return c.Status(400).JSON(fiber.Map{
						"error": fmt.Sprintf("Field '%s' must be a number", field.Label),
					})
				}
			case models.Boolean:
				if _, ok := value.(bool); !ok {
					return c.Status(400).JSON(fiber.Map{
						"error": fmt.Sprintf("Field '%s' must be a boolean", field.Label),
					})
				}
			case models.Date:
				dateStr, ok := value.(string)
				if !ok {
					return c.Status(400).JSON(fiber.Map{
						"error": fmt.Sprintf("Field '%s' must be a date string", field.Label),
					})
				}
				if _, err := time.Parse(time.RFC3339, dateStr); err != nil {
					return c.Status(400).JSON(fiber.Map{
						"error": fmt.Sprintf("Field '%s' must be a valid ISO 8601 date", field.Label),
					})
				}
			case models.Relation:
				idStr, ok := value.(string)
				if !ok {
					return c.Status(400).JSON(fiber.Map{
						"error": fmt.Sprintf("Field '%s' must be a UUID string", field.Label),
					})
				}
				if _, err := uuid.Parse(idStr); err != nil {
					return c.Status(400).JSON(fiber.Map{
						"error": fmt.Sprintf("Field '%s' must be a valid UUID", field.Label),
					})
				}
			case models.Enum:
				// TODO: Implement enum validation
				if _, ok := value.(string); !ok {
					return c.Status(400).JSON(fiber.Map{
						"error": fmt.Sprintf("Field '%s' must be a string", field.Label),
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
