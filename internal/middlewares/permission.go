package middlewares

import (
	"contentive/internal/models"

	"github.com/gofiber/fiber/v2"
)

// RequirePermission is a middleware that checks if the user has the required permission
func RequirePermission(permissionType models.PermissionType) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Load the user from the context(set by AuthMiddleware)
		user := c.Locals("user").(*models.User)

		// If the user is a super admin, skip the permission check
		if user.Role.Type == models.SuperAdmin {
			return c.Next()
		}

		// Check if the user has the required permission
		hasPermission := false
		for _, permission := range user.Role.Permissions {
			if permission.Type == permissionType {
				hasPermission = true
				break
			}
		}

		// If the user does not have the required permission, return a 403 Forbidden response
		if !hasPermission {
			return c.Status(403).JSON(fiber.Map{
				"error": "Permission denied",
			})
		}

		return c.Next()
	}
}
