package middlewares

import (
	"contentive/config"
	"contentive/internal/logger"
	"contentive/internal/models"
	"regexp"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func ValidateContentType() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var contentType models.ContentType
		if err := c.BodyParser(&contentType); err != nil {
			logger.Error("Error parsing request body", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid request body",
			})
		}

		// For PUT requests, handle partial updates
		if c.Method() == "PUT" {
			// Check if request body is empty
			if contentType.Name == "" && contentType.Type == "" && contentType.Slug == "" {
				logger.Error("Request body is empty")
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "At least one field must be provided for update",
				})
			}

			// Get identifier from URL
			identifier := c.Params("identifier")
			var existingType models.ContentType

			// Try to find by UUID first
			if uid, err := uuid.Parse(identifier); err == nil {
				if err := config.DB.First(&existingType, "id = ?", uid).Error; err != nil {
					logger.Error("Error finding content type by UUID", err)
					return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
						"error": "Content type not found",
					})
				}
			} else {
				// If not UUID, try to find by slug
				if err := config.DB.First(&existingType, "slug = ?", identifier).Error; err != nil {
					logger.Error("Error finding content type by slug", err)
					return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
						"error": "Content type not found",
					})
				}
			}

			// Prevent type modification
			if contentType.Type != "" && contentType.Type != existingType.Type {
				logger.Error("Content type cannot be changed after creation")
				return c.Status(400).JSON(fiber.Map{
					"error": "Content type cannot be changed after creation",
				})
			}

			// Check name uniqueness if name is being updated
			if contentType.Name != "" && contentType.Name != existingType.Name {
				var duplicateName models.ContentType
				if err := config.DB.Where("name = ? AND id != ?", contentType.Name, existingType.ID).First(&duplicateName).Error; err == nil {
					logger.Error("Content type with this name already exists")
					return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
						"error": "Content type with this name already exists",
					})
				}
			}

			// Check slug uniqueness if slug is being updated
			if contentType.Slug != "" && contentType.Slug != existingType.Slug {
				// Validate slug format
				if !isValidSlug(contentType.Slug) {
					logger.Error("Invalid slug format. Use only lowercase letters, numbers, and hyphens")
					return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
						"error": "Invalid slug format. Use only lowercase letters, numbers, and hyphens",
					})
				}

				var duplicateSlug models.ContentType
				if err := config.DB.Where("slug = ? AND id != ?", contentType.Slug, existingType.ID).First(&duplicateSlug).Error; err == nil {
					logger.Error("Content type with this slug already exists")
					return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
						"error": "Content type with this slug already exists",
					})
				}
			}

			// Set type to existing type for consistency
			contentType.Type = existingType.Type
		} else {
			// For POST requests, validate required fields
			if contentType.Name == "" {
				logger.Error("Name is required")
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Name is required",
				})
			}

			// Validate slug presence and format
			if contentType.Slug == "" {
				logger.Error("Slug is required")
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Slug is required",
				})
			}

			if !isValidSlug(contentType.Slug) {
				logger.Error("Invalid slug format. Use only lowercase letters, numbers, and hyphens")
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Invalid slug format. Use only lowercase letters, numbers, and hyphens",
				})
			}

			// Check slug uniqueness for new content types
			var existingSlug models.ContentType
			if err := config.DB.Where("slug = ?", contentType.Slug).First(&existingSlug).Error; err == nil {
				logger.Error("Content type with this slug already exists")
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Content type with this slug already exists",
				})
			}

			// Validate type for new content types
			if contentType.Type != string(models.Single) && contentType.Type != string(models.Collection) {
				logger.Error("Type must be'single' or 'collection'")
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Type must be 'single' or 'collection'",
				})
			}

			// Check name uniqueness for new content types
			var existingType models.ContentType
			if err := config.DB.Where("name = ?", contentType.Name).First(&existingType).Error; err == nil {
				logger.Error("Content type with this name already exists")
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Content type with this name already exists",
				})
			}
		}

		logger.Info("Content type validation successful")

		c.Locals("contentType", contentType)
		return c.Next()
	}
}

// Helper function to validate slug format
func isValidSlug(slug string) bool {
	// Only allow lowercase letters, numbers, and hyphens
	// Must start and end with a letter or number
	pattern := regexp.MustCompile(`^[a-z0-9][a-z0-9-]*[a-z0-9]$`)
	return pattern.MatchString(slug)
}
