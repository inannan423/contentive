package handler

import (
	"contentive/internal/database"
	"contentive/internal/logger"
	"contentive/internal/models"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

func isValidContentSlug(slug string) bool {
	return slug == strings.ToLower(slug) &&
		!strings.Contains(slug, " ") &&
		!strings.Contains(slug, "_")
}

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// validateContentData validates the content data against the schema fields
func validateContentData(data map[string]interface{}, fields []models.FieldDefinition) error {
	for _, field := range fields {
		value, exists := data[field.Name]
		if !exists {
			if field.Required {
				return fmt.Errorf("required field %s is missing", field.Name)
			}
			continue
		}

		switch field.Type {
		case models.FieldTypeText, models.FieldTypeTextarea, models.FieldTypeRichText:
			strVal, ok := value.(string)
			if !ok {
				return fmt.Errorf("field '%s' must be a string", field.Name)
			}
			// Check if the value is within the length range
			if maxLen, exists := field.Options["maxLength"]; exists {
				if maxLength, ok := maxLen.(float64); ok {
					if float64(len(strVal)) > maxLength {
						return fmt.Errorf("field '%s' exceeds maximum length of %v", field.Name, maxLength)
					}
				}
			}
			if minLen, exists := field.Options["minLength"]; exists {
				if minLength, ok := minLen.(float64); ok {
					if float64(len(strVal)) < minLength {
						return fmt.Errorf("field '%s' is shorter than minimum length of %v", field.Name, minLength)
					}
				}
			}

		case models.FieldTypeNumber:
			numVal, ok := value.(float64)
			if !ok {
				return fmt.Errorf("field '%s' must be a number", field.Name)
			}
			// Check if the value is within the range
			if min, exists := field.Options["min"]; exists {
				if minVal, ok := min.(float64); ok && numVal < minVal {
					return fmt.Errorf("field '%s' is less than minimum value of %v", field.Name, minVal)
				}
			}
			if max, exists := field.Options["max"]; exists {
				if maxVal, ok := max.(float64); ok && numVal > maxVal {
					return fmt.Errorf("field '%s' exceeds maximum value of %v", field.Name, maxVal)
				}
			}

		case models.FieldTypeBoolean:
			if _, ok := value.(bool); !ok {
				return fmt.Errorf("field '%s' must be a boolean", field.Name)
			}

		case models.FieldTypeDate, models.FieldTypeDateTime:
			strVal, ok := value.(string)
			if !ok {
				return fmt.Errorf("field '%s' must be a valid date string", field.Name)
			}
			var layout string
			if field.Type == models.FieldTypeDate {
				layout = "2006-01-02"
			} else {
				layout = time.RFC3339 // ISO 8601 format
			}

			if _, err := time.Parse(layout, strVal); err != nil {
				if field.Type == models.FieldTypeDate {
					return fmt.Errorf("field '%s' must be a valid date in YYYY-MM-DD format", field.Name)
				}
				return fmt.Errorf("field '%s' must be a valid datetime in ISO 8601 format", field.Name)
			}

		case models.FieldTypeEmail:
			strVal, ok := value.(string)
			if !ok {
				return fmt.Errorf("field '%s' must be a string", field.Name)
			}
			// Check if the value is a valid email address
			if !emailRegex.MatchString(strVal) {
				return fmt.Errorf("field '%s' is not a valid email address", field.Name)
			}

		case models.FieldTypeSelect:
			strVal, ok := value.(string)
			if !ok {
				return fmt.Errorf("field '%s' must be a string", field.Name)
			}
			// Check if the value is in the options list
			if options, exists := field.Options["options"]; exists {
				if optionsList, ok := options.([]interface{}); ok {
					valid := false
					for _, opt := range optionsList {
						if opt == strVal {
							valid = true
							break
						}
					}
					if !valid {
						return fmt.Errorf("field '%s' contains invalid option", field.Name)
					}
				}
			}

		case models.FieldTypeRelation:
			strVal, ok := value.(string)
			if !ok {
				return fmt.Errorf("field '%s' must be a string (slug of the related content)", field.Name)
			}

			targetSchema, ok := field.Options["targetSchema"]
			if !ok {
				return fmt.Errorf("field '%s' missing targetSchema option", field.Name)
			}
			targetSchemaStr, ok := targetSchema.(string)
			if !ok {
				return fmt.Errorf("field '%s' has invalid targetSchema option", field.Name)
			}

			// Check if the target schema exists
			var targetSchemaModel models.Schema
			if err := database.DB.Where("slug = ?", targetSchemaStr).First(&targetSchemaModel).Error; err != nil {
				return fmt.Errorf("field '%s' references non-existent schema '%s'", field.Name, targetSchemaStr)
			}

			// Check if the target schema content exists
			var contentEntry models.ContentEntry
			if err := database.DB.Where("content_type_id = ? AND slug = ?", targetSchemaModel.ID, strVal).First(&contentEntry).Error; err != nil {
				return fmt.Errorf("field '%s' references non-existent content '%s' in schema '%s'", field.Name, strVal, targetSchemaStr)
			}

		case models.FieldTypeMedia:
			if strVal, ok := value.(string); ok {
				var media models.Media
				if err := database.DB.Where("id = ?", strVal).First(&media).Error; err != nil {
					return fmt.Errorf("field '%s' references non-existent media '%s'", field.Name, strVal)
				}
				if mediaType, exists := field.Options["mediaType"]; exists {
					if allowedType, ok := mediaType.(string); ok && string(media.Type) != allowedType {
						return fmt.Errorf("field '%s' requires media of type '%s', but got '%s'", field.Name, allowedType, media.Type)
					}
				}
			} else if arrayVal, ok := value.([]interface{}); ok {
				for _, item := range arrayVal {
					if strVal, ok := item.(string); ok {
						var media models.Media
						if err := database.DB.Where("id = ?", strVal).First(&media).Error; err != nil {
							return fmt.Errorf("field '%s' references non-existent media '%s'", field.Name, strVal)
						}
						if mediaType, exists := field.Options["mediaType"]; exists {
							if allowedType, ok := mediaType.(string); ok && string(media.Type) != allowedType {
								return fmt.Errorf("field '%s' requires media of type '%s', but got '%s'", field.Name, allowedType, media.Type)
							}
						}
					} else {
						return fmt.Errorf("field '%s' must be an array of media IDs", field.Name)
					}
				}
			} else {
				return fmt.Errorf("field '%s' must be a media ID or an array of media IDs", field.Name)
			}

		case models.FieldTypeMediaList:
			arrayVal, ok := value.([]interface{})
			if !ok {
				return fmt.Errorf("field '%s' must be an array of media IDs", field.Name)
			}

			for _, item := range arrayVal {
				strVal, ok := item.(string)
				if !ok {
					return fmt.Errorf("field '%s' must contain only media IDs", field.Name)
				}

				var media models.Media
				if err := database.DB.Where("id = ?", strVal).First(&media).Error; err != nil {
					return fmt.Errorf("field '%s' references non-existent media '%s'", field.Name, strVal)
				}

				if mediaType, exists := field.Options["mediaType"]; exists {
					if allowedType, ok := mediaType.(string); ok && string(media.Type) != allowedType {
						return fmt.Errorf("field '%s' requires media of type '%s', but got '%s'", field.Name, allowedType, media.Type)
					}
				}
			}

		default:
			return fmt.Errorf("unsupported field type '%s'", field.Type)
		}
	}
	return nil
}

// CreateContent creates a new content entry for a given schema
func CreateContent(c *fiber.Ctx) error {
	// Get schema ID from locals
	var schemaID interface{}
	if id := c.Locals("schema_id"); id != nil {
		schemaID = id
	} else {
		// If schema_id is not in locals, get it from params
		schemaID = c.Params("schema_id")
	}

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

	// Check if slug is empty
	if input.Slug == "" {
		logger.Error("Slug is empty")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Slug cannot be empty",
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

	if err := validateContentData(input.Data, fileds); err != nil {
		logger.Error("Content data validation failed: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
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

// ContentQuery represents the query parameters for the GetContent handler
type ContentQuery struct {
	Page     int    `query:"page"`
	PageSize int    `query:"page_size"`
	OrderBy  string `query:"order_by"` // field name such as "created_at" "updated_at" "published_at
	Order    string `query:"order"`    // asc or desc
	Search   string `query:"search"`   // search query
	Status   string `query:"status"`   // is_published or not
}

// GetContent gets all content entries for a given schema
func GetContent(c *fiber.Ctx) error {
	// Get schema ID from locals
	var schemaID interface{}
	id := c.Locals("schema_id")
	fmt.Println(id)
	if id := c.Locals("schema_id"); id != nil {
		schemaID = id
	} else {
		// If schema_id is not in locals, get it from params
		schemaID = c.Params("schema_id")
	}

	query := new(ContentQuery)
	if err := c.QueryParser(query); err != nil {
		logger.Error("Error parsing query parameters: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid query parameters",
		})
	}

	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 10
	} else if query.PageSize > 100 {
		query.PageSize = 100
	}

	// Check if orderBy is valid
	allowedOrderBy := map[string]bool{"created_at": true, "updated_at": true, "slug": true}
	if query.OrderBy == "" {
		query.OrderBy = "created_at"
	} else if !allowedOrderBy[query.OrderBy] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid orderBy parameter",
		})
	}

	if query.Order != "asc" && query.Order != "desc" {
		query.Order = "desc"
	}

	// Check if schema exists
	var schema models.Schema
	if err := database.DB.Where("id =?", schemaID).First(&schema).Error; err != nil {
		logger.Error("Schema not found: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Schema not found",
		})
	}

	db := database.DB.Model(&models.ContentEntry{}).Where("content_type_id = ?", schemaID)

	// Apply filters
	if query.Status != "" {
		switch query.Status {
		case "published":
			db = db.Where("is_published = ?", true)
		case "draft":
			db = db.Where("is_published = ?", false)
		default:
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid status parameter",
			})
		}
	}

	if query.Search != "" {
		searchPattern := "%" + query.Search + "%"
		db = db.Where("slug LIKE ? OR data::text LIKE ?", searchPattern, searchPattern)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		logger.Error("Error counting content entries: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to count content entries",
		})
	}

	offset := (query.Page - 1) * query.PageSize
	var content []models.ContentEntry
	if err := db.Order(fmt.Sprintf("%s %s", query.OrderBy, query.Order)).
		Offset(offset).
		Limit(query.PageSize).
		Find(&content).Error; err != nil {
		logger.Error("Failed to get content: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get content",
		})
	}

	totalPages := (total + int64(query.PageSize) - 1) / int64(query.PageSize)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": content,
		"pagination": fiber.Map{
			"current_page": query.Page,
			"page_size":    query.PageSize,
			"total_pages":  totalPages,
			"total":        total,
		},
		"query": fiber.Map{
			"order_by": query.OrderBy,
			"order":    query.Order,
			"search":   query.Search,
			"status":   query.Status,
		},
	})
}

// GetContentById gets a content entry by ID
func GetContentById(c *fiber.Ctx) error {
	schemaID := c.Params("schema_id")
	contentID := c.Params("content_id")
	// Check if schema exists
	var schema models.Schema
	if err := database.DB.Where("id =?", schemaID).First(&schema).Error; err != nil {
		logger.Error("Schema not found: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Schema not found",
		})
	}
	// Check if content exists and belongs to the schema
	var content models.ContentEntry
	if err := database.DB.Where("id =? AND content_type_id =?", contentID, schemaID).First(&content).Error; err != nil {
		logger.Error("Content not found: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Content not found",
		})
	}
	return c.Status(fiber.StatusOK).JSON(content)
}

// GetContentBySlug gets a content entry by slug
func GetContentBySlug(c *fiber.Ctx) error {
	schema := c.Locals("schema").(models.Schema)
	contentSlug := c.Params("content_slug")

	var content models.ContentEntry
	if err := database.DB.Where("content_type_id = ? AND slug = ?", schema.ID, contentSlug).First(&content).Error; err != nil {
		logger.Error("Content not found: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Content not found",
		})
	}

	// Check if user is API user and content is not published
	if _, ok := c.Locals("user").(models.APIUser); ok && !content.IsPublished {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Content not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(content)
}

// UpdateContent updates an existing content entry for a given schema
func UpdateContent(c *fiber.Ctx) error {
	// Get schema ID from locals
	var schemaID interface{}
	if id := c.Locals("schema_id"); id != nil {
		schemaID = id
	} else {
		// If schema_id is not in locals, get it from params
		schemaID = c.Params("schema_id")
	}

	var contentID interface{}
	if id := c.Locals("content_id"); id != nil {
		contentID = id
	} else {
		// If content_id is not in locals, get it from params
		contentID = c.Params("content_id")
	}

	// Check if schema exists
	var schema models.Schema
	if err := database.DB.Where("id = ?", schemaID).First(&schema).Error; err != nil {
		logger.Error("Schema not found: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Schema not found",
		})
	}

	// Check if content exists and belongs to the schema
	var existingContent models.ContentEntry
	if err := database.DB.Where("id = ? AND content_type_id = ?", contentID, schemaID).First(&existingContent).Error; err != nil {
		logger.Error("Content not found: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Content not found",
		})
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

	// Check if slug is provided and valid
	if input.Slug != "" {
		if !isValidContentSlug(input.Slug) {
			logger.Error("Invalid slug format")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid slug format, must be lowercase, no spaces or underscores",
			})
		}

		// Check if new slug already exists (excluding current content)
		if err := database.DB.Where("slug = ? AND content_type_id = ? AND id != ?", input.Slug, schemaID, contentID).
			First(&models.ContentEntry{}).Error; err == nil {
			logger.Error("Content with slug already exists")
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Content with slug already exists",
			})
		}
	}

	// Unmarshal the schema fields
	var fields []models.FieldDefinition
	if err := json.Unmarshal(schema.Fields, &fields); err != nil {
		logger.Error("Error unmarshalling schema fields: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	// If data is provided, validate it
	if input.Data != nil {
		// Merge existing data with new data
		var existingData map[string]interface{}
		if err := json.Unmarshal(existingContent.Data, &existingData); err != nil {
			logger.Error("Error unmarshalling existing content data: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal server error",
			})
		}

		// Update existing data with new data
		for key, value := range input.Data {
			existingData[key] = value
		}

		// Validate all fields after merge
		if err := validateContentData(existingData, fields); err != nil {
			logger.Error("Content data validation failed: %v", err)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// Marshal merged data
		dataJson, err := json.Marshal(existingData)
		if err != nil {
			logger.Error("Error marshalling data: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal server error",
			})
		}
		existingContent.Data = datatypes.JSON(dataJson)
	}

	// Update slug if provided
	if input.Slug != "" {
		existingContent.Slug = input.Slug
	}

	// Get current user
	currentUser := c.Locals("user")
	if currentUser == nil {
		logger.Error("User not found")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Update user information
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

	existingContent.UpdatedByType = userType
	existingContent.UpdatedBy = &userID

	// Save the updated content
	if err := database.DB.Save(&existingContent).Error; err != nil {
		logger.Error("Failed to update content: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update content",
		})
	}

	// Log the action
	if userType == models.ContentEntryUserByTypeAdmin {
		adminUser := currentUser.(models.AdminUser)
		logger.AdminAction(
			adminUser.ID,
			adminUser.Name,
			"UPDATE_CONTENT",
			"Updated content for schema: "+schema.Name+" with slug: "+existingContent.Slug,
		)
	} else {
		apiUser := currentUser.(models.APIUser)
		logger.APIAction(
			apiUser.ID,
			apiUser.Name,
			"UPDATE_CONTENT",
			"Updated content for schema: "+schema.Name+" with slug: "+existingContent.Slug,
		)
	}

	return c.Status(fiber.StatusOK).JSON(existingContent)
}

// PublishContent publishes or unpublishes a content entry
func PublishContent(c *fiber.Ctx) error {
	// Get schema ID from locals
	var schemaID interface{}
	if id := c.Locals("schema_id"); id != nil {
		schemaID = id
	} else {
		// If schema_id is not in locals, get it from params
		schemaID = c.Params("schema_id")
	}
	var contentID interface{}
	if id := c.Locals("content_id"); id != nil {
		contentID = id
	} else {
		// If content_id is not in locals, get it from params
		contentID = c.Params("content_id")
	}

	// Check if schema exists
	var schema models.Schema
	if err := database.DB.Where("id = ?", schemaID).First(&schema).Error; err != nil {
		logger.Error("Schema not found: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Schema not found",
		})
	}

	// Check if content exists and belongs to the schema
	var content models.ContentEntry
	if err := database.DB.Where("id = ? AND content_type_id = ?", contentID, schemaID).First(&content).Error; err != nil {
		logger.Error("Content not found: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Content not found",
		})
	}

	var input struct {
		IsPublished bool `json:"is_published"`
	}

	if err := c.BodyParser(&input); err != nil {
		logger.Error("Error parsing request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Get current user
	currentUser := c.Locals("user")
	if currentUser == nil {
		logger.Error("User not found")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Get user information
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

	// Update content publish status
	content.IsPublished = input.IsPublished
	content.UpdatedByType = userType
	content.UpdatedBy = &userID

	if input.IsPublished {
		// Set publish information when publishing
		now := time.Now()
		content.PublishedAt = &now
		content.PublishedBy = &userID
	} else {
		// Clear publish information when unpublishing
		content.PublishedAt = nil
		content.PublishedBy = nil
	}

	// Save the updated content
	if err := database.DB.Save(&content).Error; err != nil {
		logger.Error("Failed to update content publish status: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update content publish status",
		})
	}

	// Log the action
	action := "PUBLISH_CONTENT"
	actionDesc := "Published content"
	if !input.IsPublished {
		action = "UNPUBLISH_CONTENT"
		actionDesc = "Unpublished content"
	}

	if userType == models.ContentEntryUserByTypeAdmin {
		adminUser := currentUser.(models.AdminUser)
		logger.AdminAction(
			adminUser.ID,
			adminUser.Name,
			action,
			actionDesc+" for schema: "+schema.Name+" with slug: "+content.Slug,
		)
	} else {
		apiUser := currentUser.(models.APIUser)
		logger.APIAction(
			apiUser.ID,
			apiUser.Name,
			action,
			actionDesc+" for schema: "+schema.Name+" with slug: "+content.Slug,
		)
	}

	return c.Status(fiber.StatusOK).JSON(content)
}

// DeleteContent deletes an existing content entry
func DeleteContent(c *fiber.Ctx) error {

	// Get schema ID from locals
	var schemaID interface{}
	if id := c.Locals("schema_id"); id != nil {
		schemaID = id
	} else {
		// If schema_id is not in locals, get it from params
		schemaID = c.Params("schema_id")
	}
	var contentID interface{}
	if id := c.Locals("content_id"); id != nil {
		contentID = id
	} else {
		// If content_id is not in locals, get it from params
		contentID = c.Params("content_id")
	}

	// Check if schema exists
	var schema models.Schema
	if err := database.DB.Where("id = ?", schemaID).First(&schema).Error; err != nil {
		logger.Error("Schema not found: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Schema not found",
		})
	}

	// Check if content exists and belongs to the schema
	var content models.ContentEntry
	if err := database.DB.Where("id = ? AND content_type_id = ?", contentID, schemaID).First(&content).Error; err != nil {
		logger.Error("Content not found: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Content not found",
		})
	}

	// Get current user
	currentUser := c.Locals("user")
	if currentUser == nil {
		logger.Error("User not found")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Get user information
	var userType models.ContentEntryUserByType
	var userID uuid.UUID
	var userName string

	if adminUser, ok := currentUser.(models.AdminUser); ok {
		userType = models.ContentEntryUserByTypeAdmin
		userID = adminUser.ID
		userName = adminUser.Name
	} else if apiUser, ok := currentUser.(models.APIUser); ok {
		userType = models.ContentEntryUserByTypeAPI
		userID = apiUser.ID
		userName = apiUser.Name
	} else {
		logger.Error("Invalid user type")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user type",
		})
	}

	// Start a transaction
	tx := database.DB.Begin()
	if tx.Error != nil {
		logger.Error("Failed to start transaction: %v", tx.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	// Delete the content
	if err := tx.Delete(&content).Error; err != nil {
		tx.Rollback()
		logger.Error("Failed to delete content: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete content",
		})
	}

	// Log the action before committing the transaction
	if userType == models.ContentEntryUserByTypeAdmin {
		logger.AdminAction(
			userID,
			userName,
			"DELETE_CONTENT",
			"Deleted content for schema: "+schema.Name+" with slug: "+content.Slug,
		)
	} else {
		logger.APIAction(
			userID,
			userName,
			"DELETE_CONTENT",
			"Deleted content for schema: "+schema.Name+" with slug: "+content.Slug,
		)
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		logger.Error("Failed to commit transaction: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete content",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Content deleted successfully",
		"content": content,
	})
}
