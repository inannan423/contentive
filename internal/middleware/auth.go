package middleware

import (
	"contentive/internal/config"
	"contentive/internal/logger"
	"contentive/internal/models"
	"contentive/internal/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func AuthenticateAdminUserJWT() fiber.Handler {
	return func(c *fiber.Ctx) error {

		authHeader := c.Get("Authorization")
		if authHeader == "" {
			logger.Error("Missing Authorization header")
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing Authorization header",
			})
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		claims, err := utils.ValidateAdminUserToken(tokenString)
		if err != nil {
			logger.Error("Invalid token: %v", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token",
			})
		}

		// Set user ID in context
		var user models.AdminUser
		if err := config.DB.First(&user, claims.UserID).Error; err != nil {
			logger.Error("Failed to fetch user: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to fetch user",
			})
		}

		// Set user in context
		c.Locals("user", user)
		return c.Next()
	}
}

// RequireRole is a middleware that checks if the user has the required role or higher
func RequireRole(requiredRole models.AdminUserRole) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user").(models.AdminUser)

		// Check if user has required role or higher
		if user.Role.GetRoleLevel() >= requiredRole.GetRoleLevel() {
			return c.Next()
		}

		logger.Error("Access denied for user %s with role %s, required role: %s",
			user.Email, user.Role, requiredRole)
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Access denied",
		})
	}
}
