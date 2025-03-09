package middlewares

import (
	"contentive/config"
	"contentive/internal/logger"
	"contentive/internal/models"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func ValidateField() fiber.Handler {
	return func(c *fiber.Ctx) error {
		identifier := c.Params("identifier")
		var contentType models.ContentType
		if err := config.DB.First(&contentType, "slug = ?", identifier).Error; err != nil {
			logger.Error("Error fetching content type: %v", err)
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Content type not found",
			})
		}

		c.Locals("contentTypeID", contentType.ID)

		if fieldID := c.Params("id"); fieldID != "" {
			uid, err := uuid.Parse(fieldID)
			if err != nil {
				logger.Error("Error parsing field ID: %v", err)
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Invalid field ID",
				})
			}
			c.Locals("fieldID", uid)
		}

		if c.Method() != "DELETE" {
			var field models.Field
			if err := c.BodyParser(&field); err != nil {
				logger.Error("Error parsing request body: %v", err)
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Invalid request body",
				})
			}

			if field.Label == "" {
				logger.Error("Label cannot be empty")
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Label cannot be empty",
				})
			}

			var existingField models.Field
			query := config.DB.Where("content_type_id = ? AND label = ?", contentType.ID, field.Label)

			if fieldID, ok := c.Locals("fieldID").(uuid.UUID); ok {
				query = query.Where("id != ?", fieldID)
			}

			if err := query.First(&existingField).Error; err == nil {
				logger.Error("Field with label '%s' already exists in this content type", field.Label)
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": fmt.Sprintf("Field with label '%s' already exists in this content type", field.Label),
				})
			}

			if field.Type == "" {
				logger.Error("Type cannot be empty")
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Type cannot be empty",
				})
			}

			// If the field is a relation, validate the relation type
			if field.Type == models.Relation {
				if *field.RelationType == "" {
					logger.Error("Relation type cannot be empty")
					return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
						"error": "Relation type cannot be empty",
					})
				}

				validRelationTypes := []models.RelationTypeEnum{
					models.OneToOne,
					models.OneToMany,
					models.ManyToOne,
					models.ManyToMany,
				}

				isValidRelationType := false
				for _, rt := range validRelationTypes {
					if *field.RelationType == rt {
						isValidRelationType = true
						break
					}
				}

				if !isValidRelationType {
					logger.Error("Invalid relation type")
					return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
						"error": "Invalid relation type",
					})
				}

				// Check if target content type exists
				var targetType models.ContentType
				if err := config.DB.First(&targetType, "id = ?", field.TargetTypeID).Error; err != nil {
					logger.Error("Target content type not found")
					return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
						"error": "Target content type not found",
					})
				}

			}

			if c.Method() == "POST" && !field.Required {
				field.Required = false
			}

			validTypes := map[models.FieldTypeEnum]bool{
				models.Text:     true,
				models.RichText: true,
				models.Number:   true,
				models.Date:     true,
				models.Boolean:  true,
				models.Enum:     true,
				models.Relation: true,
			}

			if !validTypes[field.Type] {
				logger.Error("Invalid field type")
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Invalid field type",
				})
			}

			c.Locals("field", field)
		}

		logger.Info("Validated field")

		return c.Next()
	}
}
