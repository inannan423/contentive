package middlewares

import (
	"contentive/config"
	"contentive/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func ValidateField() fiber.Handler {
	return func(c *fiber.Ctx) error {
		identifier := c.Params("identifier")
		var contentType models.ContentType
		if err := config.DB.First(&contentType, "slug = ?", identifier).Error; err != nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Content type not found",
			})
		}

		c.Locals("contentTypeID", contentType.ID)

		if fieldID := c.Params("id"); fieldID != "" {
			uid, err := uuid.Parse(fieldID)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Invalid field ID",
				})
			}
			c.Locals("fieldID", uid)
		}

		if c.Method() != "DELETE" {
			var field models.Field
			if err := c.BodyParser(&field); err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Invalid request body",
				})
			}

			if field.Label == "" {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Label cannot be empty",
				})
			}

			if field.Type == "" {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Type cannot be empty",
				})
			}

			if c.Method() == "POST" && !field.Required {
				field.Required = false
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
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Invalid field type",
				})
			}

			c.Locals("field", field)
		}

		return c.Next()
	}
}
