package handlers

import (
	"contentive/config"
	"contentive/internal/models"
	"contentive/internal/utils"

	"github.com/gofiber/fiber/v2"
)

func GetAPIRoles(c *fiber.Ctx) error {
	var roles []models.APIRole
	if err := config.DB.Find(&roles).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch roles",
		})
	}

	return c.Status(fiber.StatusOK).JSON(roles)
}

func GetAPIRole(c *fiber.Ctx) error {
	id := c.Params("id")
	var role models.APIRole
	if err := config.DB.Preload("Permissions.ContentType").First(&role, "id = ?", id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Role not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(role)
}

func CreateAPIRole(c *fiber.Ctx) error {
	var input struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	apiKey, err := utils.GenerateAPIKey()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate API key",
		})
	}

	role := models.APIRole{
		Name:        input.Name,
		Type:        models.CustomAPIRole,
		Description: input.Description,
		APIKey:      apiKey,
		IsSystem:    false,
	}

	if err := config.DB.Create(&role).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create role",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(role)
}

func UpdateAPIRole(c *fiber.Ctx) error {
	id := c.Params("id")
	var role models.APIRole
	if err := config.DB.First(&role, "id =?", id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Role not found",
		})
	}

	if role.IsSystem {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Cannot update system role",
		})
	}

	var input struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	if input.Name != "" {
		role.Name = input.Name
	}
	if input.Description != "" {
		role.Description = input.Description
	}

	if err := config.DB.Save(&role).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update role",
		})
	}

	return c.Status(fiber.StatusOK).JSON(role)
}

func DeleteAPIRole(c *fiber.Ctx) error {
	id := c.Params("id")
	var role models.APIRole
	if err := config.DB.First(&role, "id =?", id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Role not found",
		})
	}

	if role.IsSystem {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Cannot delete system role",
		})
	}

	if err := config.DB.Delete(&role).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete role",
		})
	}

	return c.Status(fiber.StatusNoContent).Send(nil)
}

func RegenerateAPIKey(c *fiber.Ctx) error {
	id := c.Params("id")
	var role models.APIRole
	if err := config.DB.First(&role, "id =?", id).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Role not found",
		})
	}

	if role.Type == models.PublicUser {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot regenerate API key for public user role",
		})
	}

	apiKey, err := utils.GenerateAPIKey()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate API key",
		})
	}

	role.APIKey = apiKey
	if err := config.DB.Save(&role).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update role",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"id":      role.ID,
		"name":    role.Name,
		"api_key": apiKey,
	})
}
