package controllers

import (
	"contentive/config"
	"contentive/models"
	"net/http"

	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func CreateContentType(c *fiber.Ctx) error {
	var contentType models.ContentType
	if err := c.BodyParser(&contentType); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if contentType.Name == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Name is required"})
	} else if contentType.Slug == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Slug is required"})
	} else if !contentType.Type.IsValid() {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid Type value"})
	}

	// Check if the slug already exists
	if err := config.DB.Where("slug =?", contentType.Slug).First(&models.ContentType{}).Error; err == nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Slug already exists"})
	}

	if err := config.DB.Create(&contentType).Error; err != nil {
		log.Error(err)
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return c.Status(http.StatusConflict).JSON(fiber.Map{"error": "Content type name or id already exists"})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create content type"})
	}

	// Preload the fields for the created content type
	if err := config.DB.Preload("Fields").First(&contentType, contentType.ID).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch content type"})
	}

	return c.Status(http.StatusOK).JSON(contentType)
}

func GetContentTypes(c *fiber.Ctx) error {
	var contentTypes []models.ContentType
	if err := config.DB.Find(&contentTypes).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch content types"})
	}

	// Preload the fields for each content type
	if err := config.DB.Preload("Fields").Find(&contentTypes).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch content types"})
	}

	return c.Status(http.StatusOK).JSON(contentTypes)
}

func GetContentTypeBySlug(c *fiber.Ctx) error {
	slug := c.Params("slug")

	var contentType models.ContentType
	if err := config.DB.Where("slug = ?", slug).First(&contentType).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Content type not found"})
	}
	// Preload the fields for the content type
	if err := config.DB.Preload("Fields").First(&contentType, contentType.ID).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch content type"})
	}
	return c.Status(http.StatusOK).JSON(contentType)
}
