package middlewares

import (
	"contentive/config"
	"contentive/internal/logger"
	"contentive/internal/models"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func APIAuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		logger.Info("API request received - Method: %s, Path: %s", c.Method(), c.Path())

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
				logger.Error("Cannot find public role: %v", err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Cannot find public role",
				})
			}
			c.Locals("apiRole", &publicRole)
			return c.Next()
		}

		var apiRole models.APIRole
		if err := config.DB.Where("api_key = ?", apiKey).First(&apiRole).Error; err != nil {
			logger.Error("Cannot find API role: %v", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid API key",
			})
		}

		// Check if API key has expired
		if apiRole.ExpiresAt != nil && time.Now().After(*apiRole.ExpiresAt) {
			logger.Error("API key has expired: %v", apiRole.ExpiresAt)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":      "API key has expired",
				"expired_at": apiRole.ExpiresAt,
			})
		}

		logger.Info("API key is valid: %v", apiRole.Name)

		c.Locals("apiRole", &apiRole)
		return c.Next()
	}
}

func APIPermissionMiddleware(operation models.OperationType) fiber.Handler {
	return func(c *fiber.Ctx) error {

		logger.Info("API request received - Method: %s, Path: %s", c.Method(), c.Path())

		apiRole, ok := c.Locals("apiRole").(*models.APIRole)
		if !ok {
			logger.Error("Cannot find API role")
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
				logger.Info("Content type identifier is UUID: %v", uuid)
				if err := config.DB.First(&contentType, "id = ?", uuid).Error; err != nil {
					return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
						"error": "Content type not found",
					})
				}
			} else {
				if err := config.DB.First(&contentType, "slug = ?", identifier).Error; err != nil {
					logger.Error("Content type identifier is not UUID: %v", identifier)
					return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
						"error": "Content type not found",
					})
				}
			}
			contentTypeID = contentType.ID
		} else {
			logger.Error("Invalid content type identifier: %v", identifier)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid content type identifier",
			})
		}

		var permission models.APIPermission
		if err := config.DB.Where(
			"api_role_id = ? AND content_type_id = ? AND operation = ? AND enabled = true",
			apiRole.ID, contentTypeID, operation,
		).First(&permission).Error; err != nil {
			logger.Error("Permission denied: %v", err)
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Permission denied",
				"details": fiber.Map{
					"role":        apiRole.Name,
					"operation":   operation,
					"contentType": identifier,
				},
			})
		}

		logger.Info("Permission granted: %v", permission)

		c.Locals("contentTypeID", contentTypeID)
		return c.Next()
	}
}
