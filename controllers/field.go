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
	} else if field.Slug == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Slug is required"})
	} else if !field.Type.IsValid() {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid field type"})
	}

	// Check if the field name already exists in the given content type
	var existingField models.Field
	if err := config.DB.Where("content_type_id = ? AND name = ?", field.ContentTypeID, field.Name).First(&existingField).Error; err == nil {
		return c.Status(http.StatusConflict).JSON(fiber.Map{"error": "Field name already exists in this content type"})
	}

	// Check if the slug already exists
	if err := config.DB.Where("slug =?", field.Slug).First(&models.Field{}).Error; err == nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Slug already exists"})
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

	// Preload the fields for the created content type
	if err := config.DB.Preload("ContentType.Fields").First(&field, field.ID).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch created field with content type fields"})
	}

	return c.Status(http.StatusOK).JSON(field)
}

func GetFields(c *fiber.Ctx) error {
	var fields []models.Field
	if err := config.DB.Preload("ContentType").Find(&fields).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch fields"})
	}

	// Preload the fields for each content type
	if err := config.DB.Preload("ContentType.Fields").Find(&fields).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch fields with content type fields"})
	}

	return c.Status(http.StatusOK).JSON(fields)
}

func GetFieldBySlug(c *fiber.Ctx) error {
	fieldSlug := c.Params("slug")
	var field models.Field
	if err := config.DB.Where("slug =?", fieldSlug).First(&field).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Field not found"})
	}
	return c.Status(http.StatusOK).JSON(field)
}

func GetFieldsByContentType(c *fiber.Ctx) error {
	contentTypeID := c.Params("contentTypeId")

	var contentType models.ContentType
	if err := config.DB.First(&contentType, contentTypeID).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch content type"})
	}

	var fields []models.Field
	if err := config.DB.Where("content_type_id = ?", contentTypeID).Find(&fields).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch fields"})
	}

	cleanContentType := models.ContentType{
		ID:   contentType.ID,
		Type: contentType.Type,
		Name: contentType.Name,
	}

	for i := range fields {
		fields[i].ContentType = cleanContentType
	}

	return c.Status(http.StatusOK).JSON(fields)
}
