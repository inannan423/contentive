package handler

import (
	"contentive/internal/database"
	"contentive/internal/logger"
	"contentive/internal/models"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
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

	for i := range input.Fields {
		if input.Fields[i].ID == uuid.Nil {
			input.Fields[i].ID = uuid.New()
		}
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
	// Get existing schema
	if err := database.DB.Where("id = ?", id).First(&schema).Error; err != nil {
		logger.Error("Schema not found: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Schema not found",
		})
	}

	// Define input struct with pointer fields to support partial updates
	var input struct {
		Name   *string                   `json:"name"`
		Type   *models.SchemaType        `json:"type"`
		Slug   *string                   `json:"slug"`
		Fields *[]models.FieldDefinition `json:"fields"`
	}

	// Parse request body
	if err := c.BodyParser(&input); err != nil {
		logger.Error("Error parsing request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Check uniqueness of new values
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
		if !isValidSlug(*input.Slug) {
			logger.Error("Invalid slug format")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid slug format, must be lowercase with no spaces or underscores",
			})
		}
		newSlug = *input.Slug
	}

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

	// Check if schema with same name or slug exists
	var existingSchema models.Schema
	if err := database.DB.Where("id != ? AND (slug = ? OR name = ?)", id, newSlug, newName).First(&existingSchema).Error; err == nil {
		logger.Error("Schema with slug or name already exists")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Schema with slug or name already exists",
		})
	}

	// Update schema fields
	if input.Name != nil {
		schema.Name = newName
	}
	if input.Type != nil {
		schema.Type = newType
	}
	if input.Slug != nil {
		schema.Slug = newSlug
	}

	// Handle field updates if new field definitions are provided
	if input.Fields != nil {
		// Get existing fields
		var existingFields []models.FieldDefinition
		if err := json.Unmarshal(schema.Fields, &existingFields); err != nil {
			logger.Error("Error unmarshalling existing fields: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal server error",
			})
		}

		// Process field UUIDs
		for i := range *input.Fields {
			field := &(*input.Fields)[i]
			// Generate new ID if field doesn't have one
			if field.ID == uuid.Nil {
				field.ID = uuid.New()
			}
		}

		// Validate new fields
		fieldsJSON, err := json.Marshal(*input.Fields)
		if err != nil {
			logger.Error("Error marshalling fields to JSON: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal server error",
			})
		}
		schema.Fields = datatypes.JSON(fieldsJSON)
		if err := schema.ValidateFields(); err != nil {
			logger.Error("Invalid fields: %v", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid fields: " + err.Error(),
			})
		}

		// Use transaction for field updates
		tx := database.DB.Begin()
		if tx.Error != nil {
			logger.Error("Error starting transaction: %v", tx.Error)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal server error",
			})
		}

		// Handle field changes
		if err := handleFieldChanges(tx, schema.ID, existingFields, *input.Fields); err != nil {
			tx.Rollback()
			logger.Error("Error handling field changes: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal server error",
			})
		}

		// Save updated schema
		if err := tx.Save(&schema).Error; err != nil {
			tx.Rollback()
			logger.Error("Failed to update schema: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update schema",
			})
		}

		// Commit transaction
		if err := tx.Commit().Error; err != nil {
			logger.Error("Error committing transaction: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal server error",
			})
		}
	} else {
		// Save schema directly if no field updates
		if err := database.DB.Save(&schema).Error; err != nil {
			logger.Error("Failed to update schema: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to update schema",
			})
		}
	}

	// Log admin action
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
	// if err := database.DB.Delete(&schema).Error; err != nil {
	// 	logger.Error("Failed to delete schema: %v", err)
	// 	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
	// 		"error": "Failed to delete schema",
	// 	})
	// }

	// Use a transaction to ensure atomicity
	tx := database.DB.Begin()
	if tx.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to start transaction",
		})
	}

	// Delete the content entries associated with the schema.
	if err := tx.Where("content_type_id = ?", schema.ID).Delete(&models.ContentEntry{}).Error; err != nil {
		tx.Rollback()
		logger.Error("Failed to delete content entries: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete content entries",
		})
	}

	// Delete the schema itself.
	if err := tx.Delete(&schema).Error; err != nil {
		tx.Rollback()
		logger.Error("Failed to delete schema: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete schema",
		})
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		logger.Error("Failed to commit transaction: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to commit changes",
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

// handleFieldChanges handles changes in field definitions.
func handleFieldChanges(tx *gorm.DB, schemaID uuid.UUID, oldFields, newFields []models.FieldDefinition) error {
	oldFieldMap := make(map[string]models.FieldDefinition)
	oldFieldIDMap := make(map[uuid.UUID]models.FieldDefinition)
	for _, field := range oldFields {
		oldFieldMap[field.Name] = field
		oldFieldIDMap[field.ID] = field
	}
	newFieldMap := make(map[string]models.FieldDefinition)
	for _, field := range newFields {
		if oldField, exists := oldFieldIDMap[field.ID]; exists {
			field.ID = oldField.ID
		} else if field.ID == uuid.Nil {
			field.ID = uuid.New()
		}
		newFieldMap[field.Name] = field
	}

	batchSize := 100
	var offset int

	for {
		var contents []models.ContentEntry
		if err := tx.Where("content_type_id = ?", schemaID).
			Offset(offset).
			Limit(batchSize).
			Find(&contents).Error; err != nil {
			return fmt.Errorf("failed to fetch content entries: %v", err)
		}

		if len(contents) == 0 {
			break
		}

		logger.Info("Processing %d content entries", len(contents))

		for i := range contents {
			var contentData map[string]interface{}
			if err := json.Unmarshal(contents[i].Data, &contentData); err != nil {
				return fmt.Errorf("failed to unmarshal content data: %v", err)
			}

			modified := false

			for oldFieldName := range oldFieldMap {
				if _, exists := newFieldMap[oldFieldName]; !exists {
					delete(contentData, oldFieldName)
					modified = true
				}
			}

			for newFieldName, newField := range newFieldMap {
				oldField, exists := oldFieldMap[newFieldName]
				if !exists {
					if newField.Required {
						if defaultValue, ok := newField.Options["default"]; ok {
							contentData[newFieldName] = defaultValue
							modified = true
						} else {
							return fmt.Errorf("new required field '%s' has no default value", newFieldName)
						}
					} else {
						if defaultValue, ok := newField.Options["default"]; ok {
							contentData[newFieldName] = defaultValue
						} else {
							contentData[newFieldName] = nil
						}
						modified = true
					}
				} else if oldField.ID != newField.ID {
					if value, ok := contentData[oldField.Name]; ok {
						contentData[newField.Name] = value
						delete(contentData, oldField.Name)
						modified = true
					}
				}
			}

			if modified {
				updatedData, err := json.Marshal(contentData)
				if err != nil {
					return fmt.Errorf("failed to marshal updated content data: %v", err)
				}
				contents[i].Data = datatypes.JSON(updatedData)
			}
		}

		if err := tx.Save(&contents).Error; err != nil {
			return fmt.Errorf("failed to update content entries: %v", err)
		}

		offset += batchSize
	}

	return nil
}
