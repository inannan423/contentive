package handlers

import (
	"contentive/config"
	"contentive/internal/logger"
	"contentive/internal/models"

	"github.com/gofiber/fiber/v2"
)

func GetRoles(c *fiber.Ctx) error {
	var roles []models.Role
	if err := config.DB.Preload("Permissions").Find(&roles).Error; err != nil {
		logger.Error("Error fetching roles %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve roles",
		})
	}

	logger.Info("Retrieved %d roles", len(roles))

	return c.Status(fiber.StatusOK).JSON(roles)
}

func GetRole(c *fiber.Ctx) error {
	roleID := c.Params("id")
	var role models.Role
	if err := config.DB.Preload("Permissions").First(&role, "id = ?", roleID).Error; err != nil {
		logger.Error("Error fetching role %v", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Role not found",
		})
	}

	logger.Info("Retrieved role %s", role.Name)

	return c.Status(fiber.StatusOK).JSON(role)
}

func GetPermissions(c *fiber.Ctx) error {
	var permissions []models.Permission
	if err := config.DB.Find(&permissions).Error; err != nil {
		logger.Error("Error fetching permissions %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve permissions",
		})
	}

	logger.Info("Retrieved %d permissions", len(permissions))

	return c.Status(fiber.StatusOK).JSON(permissions)
}
