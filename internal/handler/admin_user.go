package handler

import (
	"contentive/internal/config"
	"contentive/internal/logger"
	"contentive/internal/models"

	"github.com/gofiber/fiber/v2"
)

// GetAllAdminUsers returns all admin users
func GetAllAdminUsers(c *fiber.Ctx) error {
	var users []models.AdminUser

	if err := config.DB.Find(&users).Error; err != nil {
		logger.Error("Failed to fetch users: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch users",
		})
	}

	currentUser := c.Locals("user").(models.AdminUser)
	logger.AdminAction(
		currentUser.ID,
		currentUser.Name,
		"GET_ALL_USERS",
		"Retrieved all users list",
	)

	return c.Status(fiber.StatusOK).JSON(users)
}
