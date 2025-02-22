package controllers

import (
	"contentive/config"
	"contentive/models"
	"net/http"
	"net/url"
	"time"

	"github.com/gofiber/fiber/v2"
)

func validateSingleContent(data models.JSON, fields []models.Field) error {
	for _, field := range fields {
		value, exists := data[field.Name]
		if field.Required && !exists {
			return fiber.NewError(fiber.StatusBadRequest, "Required field missing: "+field.Name)
		}

		if exists {
			switch field.Type {
			case models.Number:
				switch value.(type) {
				case float64, int, int64, float32:
				default:
					return fiber.NewError(fiber.StatusBadRequest, field.Name+" must be a number")
				}
			case models.Boolean:
				if _, ok := value.(bool); !ok {
					return fiber.NewError(fiber.StatusBadRequest, field.Name+" must be a boolean")
				}
			case models.DateTime:
				if str, ok := value.(string); ok {
					if _, err := time.Parse(time.RFC3339, str); err != nil {
						return fiber.NewError(fiber.StatusBadRequest, field.Name+" must be a valid ISO 8601 datetime")
					}
				} else {
					return fiber.NewError(fiber.StatusBadRequest, field.Name+" must be a datetime string")
				}
			case models.Text:
				if _, ok := value.(string); !ok {
					return fiber.NewError(fiber.StatusBadRequest, field.Name+" must be a string")
				}
			case models.Media:
				if str, ok := value.(string); ok {
					if _, err := url.Parse(str); err != nil {
						return fiber.NewError(fiber.StatusBadRequest, field.Name+" must be a valid URL")
					}
				} else {
					return fiber.NewError(fiber.StatusBadRequest, field.Name+" must be a URL string")
				}
			}
		}
	}
	return nil
}

func CreateContent(c *fiber.Ctx) error {
	var content models.Content
	if err := c.BodyParser(&content); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	var contentType models.ContentType
	if err := config.DB.Preload("Fields").First(&contentType, content.ContentTypeID).Error; err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Content type not found"})
	}

	// Check if the slug already exists
	if err := config.DB.Where("slug =?", contentType.Slug).First(&models.Content{}).Error; err == nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Slug already exists"})
	}

	content.IsCollection = contentType.Type == models.Collection

	if content.IsCollection {
		if items, ok := content.Data["items"].([]interface{}); ok {
			delete(content.Data, "items")

			if err := config.DB.Create(&content).Error; err != nil {
				return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create collection"})
			}

			// 用于追踪已使用的 slugs
			usedSlugs := make(map[string]bool)

			for _, item := range items {
				if itemData, ok := item.(map[string]interface{}); ok {
					// 检查 item 是否包含 slug
					slug, hasSlug := itemData["slug"].(string)
					if !hasSlug || slug == "" {
						return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Each item must have a valid slug"})
					}

					// 检查 slug 是否在当前请求中重复
					if usedSlugs[slug] {
						return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Duplicate slug found: " + slug})
					}
					usedSlugs[slug] = true

					// 检查数据库中是否已存在相同的 slug
					var existingItem models.ContentItem
					if err := config.DB.Where("collection_id = ? AND slug = ?", content.ID, slug).First(&existingItem).Error; err == nil {
						return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Item slug already exists: " + slug})
					}

					contentItem := models.ContentItem{
						CollectionID: content.ID,
						Slug:        slug,
						Data:        models.JSON(itemData),
					}

					if err := validateSingleContent(contentItem.Data, contentType.Fields); err != nil {
						return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
					}

					if err := config.DB.Create(&contentItem).Error; err != nil {
						return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create content item"})
					}
				}
			}
		}
	} else {
		if err := validateSingleContent(content.Data, contentType.Fields); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		if err := config.DB.Create(&content).Error; err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create content"})
		}
	}

	if err := config.DB.Preload("ContentType.Fields").Preload("Items").First(&content, content.ID).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch created content"})
	}

	return c.Status(http.StatusOK).JSON(content)
}

func UpdateContent(c *fiber.Ctx) error {
	contentID := c.Params("id")
	var content models.Content

	if err := config.DB.Preload("ContentType.Fields").Preload("Items").First(&content, contentID).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Content not found"})
	}

	var updateData models.Content
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if content.IsCollection {
		if items, ok := updateData.Data["items"].([]interface{}); ok {
			if err := config.DB.Where("collection_id = ?", content.ID).Delete(&models.ContentItem{}).Error; err != nil {
				return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete existing items"})
			}

			delete(updateData.Data, "items")
			content.Data = updateData.Data

			// 用于追踪已使用的 slugs
			usedSlugs := make(map[string]bool)

			for _, item := range items {
				if itemData, ok := item.(map[string]interface{}); ok {
					// 检查 item 是否包含 slug
					slug, hasSlug := itemData["slug"].(string)
					if !hasSlug || slug == "" {
						return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Each item must have a valid slug"})
					}

					// 检查 slug 是否在当前请求中重复
					if usedSlugs[slug] {
						return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Duplicate slug found: " + slug})
					}
					usedSlugs[slug] = true

					contentItem := models.ContentItem{
						CollectionID: content.ID,
						Slug:        slug,
						Data:        models.JSON(itemData),
					}

					if err := validateSingleContent(contentItem.Data, content.ContentType.Fields); err != nil {
						return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
					}
					if err := config.DB.Create(&contentItem).Error; err != nil {
						return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create content item"})
					}
				}
			}
		}
	} else {
		if err := validateSingleContent(updateData.Data, content.ContentType.Fields); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		content.Data = updateData.Data
	}

	if err := config.DB.Save(&content).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update content"})
	}

	if err := config.DB.Preload("ContentType.Fields").Preload("Items").First(&content, content.ID).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch updated content"})
	}

	return c.Status(http.StatusOK).JSON(content)
}

func GetContents(c *fiber.Ctx) error {
	var contents []models.Content
	if err := config.DB.Preload("ContentType.Fields").Preload("Items").Find(&contents).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch contents"})
	}

	return c.Status(http.StatusOK).JSON(contents)
}

func GetContentsByType(c *fiber.Ctx) error {
	contentTypeID := c.Params("contentTypeId")
	var contents []models.Content

	if err := config.DB.Preload("ContentType.Fields").Preload("Items").Where("content_type_id = ?", contentTypeID).Find(&contents).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch contents"})
	}

	return c.Status(http.StatusOK).JSON(contents)
}

func DeleteContent(c *fiber.Ctx) error {
	contentID := c.Params("id")

	if err := config.DB.Where("collection_id = ?", contentID).Delete(&models.ContentItem{}).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete content items"})
	}

	if err := config.DB.Delete(&models.Content{}, contentID).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete content"})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"message": "Content deleted successfully"})
}

func GetContentItemBySlug(c *fiber.Ctx) error {
	collectionID := c.Params("collectionId")
	slug := c.Params("slug")
	var contentItem models.ContentItem
	if err := config.DB.Where("collection_id =? AND slug =?", collectionID, slug).First(&contentItem).Error; err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Content item not found"})
	}
	return c.Status(http.StatusOK).JSON(contentItem)
}
