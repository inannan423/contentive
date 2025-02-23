package middlewares

import (
	"contentive/config"
	"contentive/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func ValidateContentType() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var contentType models.ContentType
		if err := c.BodyParser(&contentType); err != nil {
			return c.Status(400).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		// For PUT requests, handle partial updates
		if c.Method() == "PUT" {
			// Check if request body is empty
			if contentType.Name == "" && contentType.Type == "" {
				return c.Status(400).JSON(fiber.Map{
					"error": "At least one field must be provided for update",
				})
			}

			// Parse and validate content type ID
			contentTypeID, err := uuid.Parse(c.Params("contentTypeId"))
			if err != nil {
				return c.Status(400).JSON(fiber.Map{
					"error": "Invalid content type ID",
				})
			}

			// Fetch existing content type
			var existingType models.ContentType
			if err := config.DB.First(&existingType, "id = ?", contentTypeID).Error; err != nil {
				return c.Status(404).JSON(fiber.Map{
					"error": "Content type not found",
				})
			}

			// Prevent type modification
			if contentType.Type != "" && contentType.Type != existingType.Type {
				return c.Status(400).JSON(fiber.Map{
					"error": "Content type cannot be changed after creation",
				})
			}

			// Check name uniqueness if name is being updated
			if contentType.Name != "" && contentType.Name != existingType.Name {
				var duplicateName models.ContentType
				if err := config.DB.Where("name = ? AND id != ?", contentType.Name, contentTypeID).First(&duplicateName).Error; err == nil {
					return c.Status(400).JSON(fiber.Map{
						"error": "Content type with this name already exists",
					})
				}
			}

			// Set type to existing type for consistency
			contentType.Type = existingType.Type
		} else {
			// For POST requests, validate required fields
			if contentType.Name == "" {
				return c.Status(400).JSON(fiber.Map{
					"error": "Name is required",
				})
			}

			// Validate type for new content types
			if contentType.Type != string(models.Single) && contentType.Type != string(models.Collection) {
				return c.Status(400).JSON(fiber.Map{
					"error": "Type must be 'single' or 'collection'",
				})
			}

			// Check name uniqueness for new content types
			var existingType models.ContentType
			if err := config.DB.Where("name = ?", contentType.Name).First(&existingType).Error; err == nil {
				return c.Status(400).JSON(fiber.Map{
					"error": "Content type with this name already exists",
				})
			}
		}

		c.Locals("contentType", contentType)
		return c.Next()
	}
}
