package middlewares

import (
	"contentive/config"
	"contentive/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func ValidateField() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if contentTypeID := c.Params("contentTypeId"); contentTypeID != "" {
			uid, err := uuid.Parse(contentTypeID)
			if err != nil {
				return c.Status(400).JSON(fiber.Map{
					"error": "Invalid content type ID",
				})
			}

			var contentType models.ContentType
			if err := config.DB.First(&contentType, uid).Error; err != nil {
				return c.Status(404).JSON(fiber.Map{
					"error": "Content type not found",
				})
			}

			c.Locals("contentTypeID", uid)
		}

		if fieldID := c.Params("id"); fieldID != "" {
			uid, err := uuid.Parse(fieldID)
			if err != nil {
				return c.Status(400).JSON(fiber.Map{
					"error": "Invalid field ID",
				})
			}
			c.Locals("fieldID", uid)
		}

		var field models.Field
		if err := c.BodyParser(&field); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		if field.Label == "" {
			return c.Status(400).JSON(fiber.Map{
				"error": "Label cannot be empty",
			})
		}

		if field.Type == "" {
			return c.Status(400).JSON(fiber.Map{
				"error": "Type cannot be empty",
			})
		}

		validTypes := map[models.FieldTypeEnum]bool{
			models.Text:     true,
			models.RichText: true,
			models.Number:   true,
			models.Date:     true,
			models.Boolean:  true,
			models.Enum:     true,
			models.Relation: true,
		}

		if !validTypes[field.Type] {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid field type",
			})
		}

		c.Locals("field", field)
		return c.Next()
	}
}
