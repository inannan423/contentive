package middlewares

import (
	"contentive/config"
	"contentive/internal/models"

	"github.com/gofiber/fiber/v2"
)

func ValidateContentType() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var contentType models.ContentType
		if err := c.BodyParser(&contentType); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		if contentType.Name == "" {
			return c.Status(400).JSON(fiber.Map{
				"error": "Name is required",
			})
		}

		if contentType.Type != string(models.Single) && contentType.Type != string(models.Collection) {
			return c.Status(400).JSON(fiber.Map{
				"error": "Type must be 'single' or 'collection'",
			})
		}

		var existingType models.ContentType
		if err := config.DB.Where("name = ?", contentType.Name).First(&existingType).Error; err == nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Content type with this name already exists",
			})
		}

		c.Locals("contentType", contentType)
		return c.Next()
	}
}
