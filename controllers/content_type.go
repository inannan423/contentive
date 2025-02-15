package controllers

import (
	"contentive/config"
	"contentive/models"
	"log"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func CreateContentType(c *fiber.Ctx) error {
	var contentType models.ContentType
	if err := c.BodyParser(&contentType); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if contentType.Name == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Name is required"})
	} else if !contentType.Type.IsValid() {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Type value"})
	}

	log.Println(contentType)

	if err := config.DB.Create(&contentType).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create content type"})
	}

	return c.Status(http.StatusOK).JSON(contentType)
}

func GetContentTypes(c *fiber.Ctx) error {
	var contentTypes []models.ContentType
	if err := config.DB.Find(&contentTypes).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch content types"})
	}

	return c.Status(http.StatusOK).JSON(contentTypes)
}
