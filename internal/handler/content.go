package handler

import (
	"contentive/internal/database"
	"contentive/internal/logger"
	"contentive/internal/models"
	"encoding/json"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// TODO: validateContentData

// CreateContent creates a new content entry for a given schema
func CreateContent(c *fiber.Ctx) error {
	schemaID := c.Params("schema_id")

	// Check if schema exists
	var schema models.Schema
	if err := database.DB.Where("id = ?", schemaID).First(&schema).Error; err != nil {
		logger.Error("Schema not found: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Schema not found",
		})
	}

	// If Schema Type is single, check if there is already a content entry
	if schema.Type == models.SchemaTypeSingle {
		var existingContent models.ContentEntry
		if err := database.DB.Where("content_type_id =?", schemaID).First(&existingContent).Error; err == nil {
			logger.Error("Single schema already has a content entry")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "This is a single schema, you can't create more than one content entry",
			})
		}
	}

	var input struct {
		Slug string                 `json:"slug"`
		Data map[string]interface{} `json:"data"`
	}

	if err := c.BodyParser(&input); err != nil {
		logger.Error("Error parsing request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Check if slug is valid
	if !isValidContentSlug(input.Slug) {
		logger.Error("Invalid slug format")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid slug format, must be lowercase, no spaces or underscores",
		})
	}

	// Check if slug already exists
	var existingContent models.ContentEntry
	if err := database.DB.Where("slug = ? AND content_type_id = ?", input.Slug, schemaID).First(&existingContent).Error; err == nil {
		logger.Error("Content with slug already exists")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Content with slug already exists",
		})
	}

	// Unmarshal the data into the schema's fields
	var fileds []models.FieldDefinition
	if err := json.Unmarshal(schema.Fields, &fileds); err != nil {
		logger.Error("Error unmarshalling schema fields: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	// Check required fields
	for _, field := range fileds {
		if field.Required {
			if _, exists := input.Data[field.Name]; !exists {
				logger.Error("Required field %s is missing", field.Name)
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Required field is missing: " + field.Name,
				})
			}
		}
	}

	// Turn data to json
	dataJson, err := json.Marshal(input.Data)
	if err != nil {
		logger.Error("Error marshalling data: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	currentUser := c.Locals("user")
	// Check if user is admin
	if currentUser == nil {
		logger.Error("User not found")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Check if user is admin
	var userType models.ContentEntryUserByType

	var userID uuid.UUID
	if adminUser, ok := currentUser.(models.AdminUser); ok {
		userType = models.ContentEntryUserByTypeAdmin
		userID = adminUser.ID
	} else if apiUser, ok := currentUser.(models.APIUser); ok {
		userType = models.ContentEntryUserByTypeAPI
		userID = apiUser.ID
	} else {
		logger.Error("Invalid user type")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user type",
		})
	}

	// Create content entry
	content := models.ContentEntry{
		Slug:          input.Slug,
		Data:          datatypes.JSON(dataJson),
		ContentTypeID: schema.ID,
		IsPublished:   false,
		CreatedByType: userType,
		UpdatedByType: userType,
		UpdatedBy:     &userID,
	}

	if err := database.DB.Create(&content).Error; err != nil {
		logger.Error("Failed to create content: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create content",
		})
	}
	// If user is admin, log admin action
	if userType == models.ContentEntryUserByTypeAdmin {
		adminUser := currentUser.(models.AdminUser)
		logger.AdminAction(
			adminUser.ID,
			adminUser.Name,
			"CREATE_CONTENT",
			"Created content for schema: "+schema.Name+" with slug: "+input.Slug,
		)
	} else {
		apiUser := currentUser.(models.APIUser)
		logger.APIAction(
			apiUser.ID,
			apiUser.Name,
			"CREATE_CONTENT",
			"Created content for schema: "+schema.Name+" with slug: "+input.Slug,
		)
	}

	return c.Status(fiber.StatusCreated).JSON(content)
}

func isValidContentSlug(slug string) bool {
	return slug == strings.ToLower(slug) &&
		!strings.Contains(slug, " ") &&
		!strings.Contains(slug, "_")
}
