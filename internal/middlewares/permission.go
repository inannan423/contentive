package middlewares

import (
	"contentive/internal/models"
	"log"

	"github.com/gofiber/fiber/v2"
)

// RequirePermission is a middleware that checks if the user has the required permission
func RequirePermission(permissionType models.PermissionType) fiber.Handler {
	return func(c *fiber.Ctx) error {
		log.Printf("=== Permission Check Start ===")
		log.Printf("Request Path: %s", c.Path())
		log.Printf("Request Method: %s", c.Method())

		user := c.Locals("user").(*models.User)

		log.Printf("User details - Name: %s, ID: %s", user.Username, user.ID)
		log.Printf("Role details - Type: %s, ID: %s", user.Role.Type, user.Role.ID)
		log.Printf("Required permission: %s", permissionType)

		if user.Role.Type == "" {
			log.Printf("Warning: User %s has empty role type", user.Username)
		}

		isSuperAdmin := user.Role.Type == models.SuperAdmin
		log.Printf("Is super admin? %v (Role.Type: %s, SuperAdmin const: %s)",
			isSuperAdmin, user.Role.Type, models.SuperAdmin)

		if isSuperAdmin {
			log.Printf("Granting access to super admin user: %s", user.Username)
			log.Printf("=== Permission Check End (Granted) ===")
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

		if !hasPermission {
			log.Printf("Permission denied: User %s does not have %s permission", user.Username, permissionType)
			log.Printf("=== Permission Check End (Denied) ===")
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Permission denied",
			})
		}

		log.Printf("=== Permission Check End (Granted) ===")
		return c.Next()
	}
}
