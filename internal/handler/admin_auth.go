package handler

import (
	"contentive/internal/database"
	"contentive/internal/logger"
	"contentive/internal/models"
	"contentive/internal/utils"
	"time"

	"github.com/gofiber/fiber/v2"
)

func AdminUserLogin(c *fiber.Ctx) error {
	// Parse input
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&input); err != nil {
		logger.Error("Failed to parse input: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	var user models.AdminUser
	if err := database.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		logger.Error("Failed to fetch user: %v", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	if !user.CheckPassword(input.Password) {
		logger.Error("Invalid password for user: %s", user.Email)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	token, err := utils.GenerateAdminUserToken(&user)
	if err != nil {
		logger.Error("Failed to generate token: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate token"})
	}

	user.LastLoginAt = time.Now()
	if err := database.DB.Save(&user).Error; err != nil {
		logger.Error("Failed to update user: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update user"})
	}

	logger.AdminAction(user.ID, user.Email, "LOGIN", "Successful login")
	return c.JSON(fiber.Map{"token": token, "user": fiber.Map{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
		"role":  string(user.Role),
	}})
}
