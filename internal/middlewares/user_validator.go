package middlewares

import (
	"contentive/config"
	"contentive/internal/models"
	"log"

	"github.com/gofiber/fiber/v2"
)

func ValidateSuperAdminOperation() fiber.Handler {
	return func(c *fiber.Ctx) error {
		targetUserID := c.Params("id")
		var targetUser models.User
		if err := config.DB.Preload("Role").First(&targetUser, "id = ?", targetUserID).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
		}

		currentUser := c.Locals("user").(*models.User)
		isSuperAdmin := targetUser.Role.Type == models.SuperAdmin

		log.Printf("Validating super admin operation - Target user: %s (ID: %s)", targetUser.Username, targetUser.ID)
		log.Printf("Current user: %s (ID: %s)", currentUser.Username, currentUser.ID)
		log.Printf("Is target super admin? %v", isSuperAdmin)

		if isSuperAdmin {
			if targetUser.ID != currentUser.ID {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"error": "Cannot modify other super admin users",
				})
			}

			var input struct {
				RoleID string `json:"role_id"`
			}
			if err := c.BodyParser(&input); err == nil && input.RoleID != "" {
				if input.RoleID != targetUser.RoleID.String() {
					return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
						"error": "Cannot change super admin role",
					})
				}
			}
		}

		c.Locals("targetUser", &targetUser)
		return c.Next()
	}
}
