package middleware

import (
	"contentive/internal/database"
	"contentive/internal/logger"
	"contentive/internal/models"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

// GetSchemaFromParams gets schema from params and sets it in locals
func GetSchemaFromParams() fiber.Handler {
	return func(c *fiber.Ctx) error {
		schemaID := c.Params("schema_id")
		var schema models.Schema
		if err := database.DB.First(&schema, "id = ?", schemaID).Error; err != nil {
			logger.Error("Schema not found: %v", err)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Schema not found",
			})
		}

		// Set schema in context for later use
		c.Locals("schema", schema)
		return c.Next()
	}
}

// GetSchemaFromSlug gets schema from slug param and sets it in locals
func GetSchemaFromSlug() fiber.Handler {
	return func(c *fiber.Ctx) error {
		schemaSlug := c.Params("schema_slug")
		var schema models.Schema
		if err := database.DB.First(&schema, "slug = ?", schemaSlug).Error; err != nil {
			logger.Error("Schema not found: %v", err)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Schema not found",
			})
		}

		c.Locals("schema", schema)
		c.Locals("schema_slug", schemaSlug)
		c.Locals("schema_id", schema.ID.String())
		c.Params("schema_id", schema.ID.String())
		fmt.Printf("Schema: %v\n", schema.ID)
		return c.Next()
	}
}

// GetContentFromSlug gets content from slug param and sets it in locals
func GetContentFromSlug() fiber.Handler {
	return func(c *fiber.Ctx) error {
		schemaSlug := c.Params("schema_slug")
		contentSlug := c.Params("content_slug")
		var schema models.Schema
		if err := database.DB.First(&schema, "slug =?", schemaSlug).Error; err != nil {
			logger.Error("Schema not found: %v", err)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Schema not found",
			})
		}
		var content models.ContentEntry
		if err := database.DB.First(&content, "slug =?", contentSlug).Error; err != nil {
			logger.Error("Content not found: %v", err)
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Content not found",
			})
		}
		c.Locals("schema", schema)
		c.Locals("content", content)
		c.Params("schema_id", schema.ID.String())
		c.Params("content_id", content.ID.String())
		return c.Next()
	}
}

// RequireSchemaScope checks if the API user has the required scope for the schema
func RequireSchemaScope(action string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		schema, ok := c.Locals("schema").(models.Schema)
		if !ok {
			logger.Error("Schema not found in context")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Schema not found in context",
			})
		}
		return RequireAPIScope(schema.Slug + ":" + action)(c)
	}
}
