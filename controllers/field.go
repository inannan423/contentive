package controllers

import (
	"contentive/config"
	"contentive/models"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func CreateField(c *fiber.Ctx) error {
	var field models.Field
	if err := c.BodyParser(&field); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if field.Name == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Name is required"})
	} else if !field.Type.IsValid() {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid field type"})
	}

	var contentType models.ContentType
	if err := config.DB.First(&contentType, field.ContentTypeID).Error; err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Content type not found"})
	}

	if err := config.DB.Create(&field).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create field"})
	}

	if err := config.DB.Preload("ContentType").First(&field, field.ID).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch created field"})
	}

	return c.Status(http.StatusOK).JSON(field)
}

func GetFields(c *fiber.Ctx) error {
	var fields []models.Field
	if err := config.DB.Preload("ContentType").Find(&fields).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch fields"})
	}

	return c.Status(http.StatusOK).JSON(fields)
}

func GetFieldsByContentType(c *fiber.Ctx) error {
	contentTypeID := c.Params("contentTypeId")
	var fields []models.Field

	if err := config.DB.Preload("ContentType").Where("content_type_id = ?", contentTypeID).Find(&fields).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch fields"})
	}

	return c.Status(http.StatusOK).JSON(fields)
}
