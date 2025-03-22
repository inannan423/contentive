package middleware

import (
	"contentive/internal/database"
	"contentive/internal/logger"
	"contentive/internal/models"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// AuthenticateAPIUserToken checks if the Authorization header is a valid token
func AuthenticateAPIUserToken() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			logger.Error("Missing Authorization header")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing Authorization header",
			})
		}

		// Check if the Authorization header is a valid token
		if !strings.HasPrefix(authHeader, "Bearer ") {
			logger.Error("Invalid Authorization header format")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid Authorization header format",
			})
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		// Find the api user with the given token
		var apiUser models.APIUser
		if err := database.DB.Where("token = ?", token).First(&apiUser).Error; err != nil {
			logger.Error("Invalid API token: %v", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid API token",
			})
		}

		if apiUser.Status != models.APIUserStatusActive {
			logger.Error("API user is not active")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "API user is not active",
			})
		}

		// Check if the API user token has expired
		if apiUser.ExpireAt != nil && apiUser.ExpireAt.Before(time.Now()) {
			logger.Error("API token has expired")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "API token has expired",
			})
		}

		// Set the API user in the context
		c.Locals("user", apiUser)

		return c.Next()
	}
}

// RequireAPIScope checks if the API user has the required scope
func RequireAPIScope(requiredScope string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the API user from the context
		apiUser, ok := c.Locals("user").(models.APIUser)
		if !ok {
			logger.Error("User is not an API user")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "User is not an API user",
			})
		}

		// Check if the API user has the required scope
		hasScope := false
		for _, scope := range apiUser.Scopes {
			if scope == requiredScope || scope == "*" {
				hasScope = true
				break
			}
		}

		if !hasScope {
			logger.Error("API user does not have the required scope: %s", requiredScope)
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Insufficient permissions",
			})
		}

		return c.Next()
	}
}
