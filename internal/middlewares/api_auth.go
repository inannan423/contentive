package middlewares

import (
	"contentive/config"
	"contentive/internal/models"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func APIAuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		apiKey := c.Get("X-API-Key")
		if apiKey == "" {
			auth := c.Get("Authorization")
			if strings.HasPrefix(auth, "Bearer ") {
				apiKey = strings.TrimPrefix(auth, "Bearer ")
			}
		}

		if apiKey == "" {
			var publicRole models.APIRole
			if err := config.DB.Where("type = ?", models.PublicUser).First(&publicRole).Error; err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Cannot find public role",
				})
			}
			c.Locals("apiRole", &publicRole)
			return c.Next()
		}

		var apiRole models.APIRole
		if err := config.DB.Where("api_key = ?", apiKey).First(&apiRole).Error; err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid API key",
			})
		}

		c.Locals("apiRole", &apiRole)
		return c.Next()
	}
}

func APIPermissionMiddleware(operation models.OperationType) fiber.Handler {
	return func(c *fiber.Ctx) error {
		apiRole, ok := c.Locals("apiRole").(*models.APIRole)
		if !ok {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Cannot find API role",
			})
		}

		var contentTypeID uuid.UUID

		identifier := c.Params("identifier")
		if identifier == "" {
			identifier = c.Params("contentTypeId")
		}

		if identifier != "" {
			var contentType models.ContentType
			if uuid, err := uuid.Parse(identifier); err == nil {
				if err := config.DB.First(&contentType, "id = ?", uuid).Error; err != nil {
					return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
						"error": "Content type not found",
					})
				}
			} else {
				if err := config.DB.First(&contentType, "slug = ?", identifier).Error; err != nil {
					return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
						"error": "Content type not found",
					})
				}
			}
			contentTypeID = contentType.ID
		} else {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid content type identifier",
			})
		}

		var permission models.APIPermission
		if err := config.DB.Where(
			"api_role_id = ? AND content_type_id = ? AND operation = ? AND enabled = true",
			apiRole.ID, contentTypeID, operation,
		).First(&permission).Error; err != nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Permission denied",
				"details": fiber.Map{
					"role":        apiRole.Name,
					"operation":   operation,
					"contentType": identifier,
				},
			})
		}

		c.Locals("contentTypeID", contentTypeID)
		return c.Next()
	}
}
