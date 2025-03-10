package handler

import (
	"contentive/internal/database"
	"contentive/internal/logger"
	"contentive/internal/models"
	"encoding/json"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/datatypes"
)

func CreateSchema(c *fiber.Ctx) error {
	var input struct {
		Name   string                   `json:"name"`
		Type   models.SchemaType        `json:"type"`
		Slug   string                   `json:"slug"`
		Fields []models.FieldDefinition `json:"fields"`
	}

	if err := c.BodyParser(&input); err != nil {
		logger.Error("Error parsing request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Check required fields
	if input.Name == "" || input.Type == "" || input.Slug == "" {
		logger.Error("Missing required fields: name, type, slug")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing required fields: name, type, slug",
		})
	}

	// Check types
	if input.Type != models.SchemaTypeList && input.Type != models.SchemaTypeSingle {
		logger.Error("Invalid schema type: %s", input.Type)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid schema type, must be list or single",
		})
	}

	if !isValidSlug(input.Slug) {
		logger.Error("Invalid slug format")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid slug format, must be lowercase, no spaces or underscores",
		})
	}

	// Check if slug or name already exists
	var existingSchema models.Schema
	if err := database.DB.Where("slug = ? OR name = ?", input.Slug, input.Name).First(&existingSchema).Error; err == nil {
		logger.Error("Schema with slug or name already exists")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Schema with slug or name already exists",
		})
	}

	// Turn fields to JSON
	fieldsJSON, err := json.Marshal(input.Fields)
	if err != nil {
		logger.Error("Error marshalling fields to JSON: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	// Create the schema
	schema := models.Schema{
		Name:   input.Name,
		Type:   input.Type,
		Slug:   input.Slug,
		Fields: datatypes.JSON(fieldsJSON),
	}

	if err := schema.ValidateFields(); err != nil {
		logger.Error("Invalid fields: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid fields: " + err.Error(),
		})
	}

	if err := database.DB.Create(&schema).Error; err != nil {
		logger.Error("Failed to create schema: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create schema",
		})
	}

	currentUser := c.Locals("user").(models.AdminUser)
	logger.AdminAction(
		currentUser.ID,
		currentUser.Name,
		"CREATE_SCHEMA",
		"Created schema: "+schema.Name,
	)

	return c.Status(fiber.StatusCreated).JSON(schema)
}

func isValidSlug(slug string) bool {
	return slug == strings.ToLower(slug) &&
		!strings.Contains(slug, " ") &&
		!strings.Contains(slug, "_")
}

// GetSchema retrieves a single schema by its ID or slug.
// The route parameter "id" can be either the schema's ID or its slug.
func GetSchema(c *fiber.Ctx) error {
	idOrSlug := c.Params("id")
	var schema models.Schema
	if err := database.DB.Where("id = ? OR slug = ?", idOrSlug, idOrSlug).First(&schema).Error; err != nil {
		logger.Error("Schema not found: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Schema not found",
		})
	}
	return c.JSON(schema)
}

// ListSchemas retrieves all schemas from the database.
func ListSchemas(c *fiber.Ctx) error {
	var schemas []models.Schema
	if err := database.DB.Find(&schemas).Error; err != nil {
		logger.Error("Error fetching schemas: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve schemas",
		})
	}
	return c.JSON(schemas)
}

// UpdateSchema supports partial update of a schema.
// Required fields (Name, Type, Slug) cannot be updated to an empty value.
// Only the provided fields will be updated.
func UpdateSchema(c *fiber.Ctx) error {
	id := c.Params("id")
	var schema models.Schema
	// Retrieve the existing schema by ID.
	if err := database.DB.Where("id = ?", id).First(&schema).Error; err != nil {
		logger.Error("Schema not found: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Schema not found",
		})
	}

	// Define an input struct with pointer fields for partial update.
	var input struct {
		Name   *string                   `json:"name"`
		Type   *models.SchemaType        `json:"type"`
		Slug   *string                   `json:"slug"`
		Fields *[]models.FieldDefinition `json:"fields"`
	}

	// Parse the request body.
	if err := c.BodyParser(&input); err != nil {
		logger.Error("Error parsing request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Determine new values for uniqueness check.
	newName := schema.Name
	if input.Name != nil {
		if *input.Name == "" {
			logger.Error("Name cannot be empty")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Name cannot be empty",
			})
		}
		newName = *input.Name
	}
	newSlug := schema.Slug
	if input.Slug != nil {
		if *input.Slug == "" {
			logger.Error("Slug cannot be empty")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Slug cannot be empty",
			})
		}
		// Validate slug format.
		if !isValidSlug(*input.Slug) {
			logger.Error("Invalid slug format")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid slug format, must be lowercase with no spaces or underscores",
			})
		}
		newSlug = *input.Slug
	}
	// If Type is provided, check that it's valid.
	newType := schema.Type
	if input.Type != nil {
		if *input.Type != models.SchemaTypeList && *input.Type != models.SchemaTypeSingle {
			logger.Error("Invalid schema type: %s", *input.Type)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid schema type, must be 'list' or 'single'",
			})
		}
		newType = *input.Type
	}

	// Check if a schema with the same slug or name already exists (excluding the current one).
	var existingSchema models.Schema
	if err := database.DB.Where("id != ? AND (slug = ? OR name = ?)", id, newSlug, newName).First(&existingSchema).Error; err == nil {
		logger.Error("Schema with slug or name already exists")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Schema with slug or name already exists",
		})
	}

	// Update the schema fields if provided.
	if input.Name != nil {
		schema.Name = newName
	}
	if input.Type != nil {
		schema.Type = newType
	}
	if input.Slug != nil {
		schema.Slug = newSlug
	}
	if input.Fields != nil {
		// Marshal the provided fields into JSON.
		fieldsJSON, err := json.Marshal(*input.Fields)
		if err != nil {
			logger.Error("Error marshalling fields to JSON: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal server error",
			})
		}
		schema.Fields = datatypes.JSON(fieldsJSON)
		// Validate the updated fields.
		if err := schema.ValidateFields(); err != nil {
			logger.Error("Invalid fields: %v", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid fields: " + err.Error(),
			})
		}
	}

	// Save the updated schema to the database.
	if err := database.DB.Save(&schema).Error; err != nil {
		logger.Error("Failed to update schema: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update schema",
		})
	}

	// Log the admin action.
	currentUser := c.Locals("user").(models.AdminUser)
	logger.AdminAction(
		currentUser.ID,
		currentUser.Name,
		"UPDATE_SCHEMA",
		"Updated schema: "+schema.Name,
	)

	return c.JSON(schema)
}

// DeleteSchema deletes an existing schema identified by the "id" route parameter.
func DeleteSchema(c *fiber.Ctx) error {
	id := c.Params("id")
	var schema models.Schema
	if err := database.DB.Where("id = ?", id).First(&schema).Error; err != nil {
		logger.Error("Schema not found: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Schema not found",
		})
	}

	// Delete the schema from the database.
	if err := database.DB.Delete(&schema).Error; err != nil {
		logger.Error("Failed to delete schema: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete schema",
		})
	}

	// Log the admin action.
	currentUser := c.Locals("user").(models.AdminUser)
	logger.AdminAction(
		currentUser.ID,
		currentUser.Name,
		"DELETE_SCHEMA",
		"Deleted schema: "+schema.Name,
	)

	// Return a No Content status.
	return c.SendStatus(fiber.StatusNoContent)
}
