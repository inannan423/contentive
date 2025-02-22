package controllers

import (
	"contentive/config"
	"contentive/models"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func GetContentItems(c *fiber.Ctx) error {
	collectionID := c.Params("collectionId")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset := (page - 1) * limit

	var items []models.ContentItem
	var total int64

	if err := config.DB.Model(&models.ContentItem{}).Where("collection_id = ?", collectionID).Count(&total).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to count items"})
	}

	if err := config.DB.Where("collection_id = ?", collectionID).Offset(offset).Limit(limit).Find(&items).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch items"})
	}

	return c.JSON(fiber.Map{
		"items": items,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func CreateContentItem(c *fiber.Ctx) error {
	var item models.ContentItem
	if err := c.BodyParser(&item); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	var collection models.Content
	if err := config.DB.Preload("ContentType.Fields").First(&collection, item.CollectionID).Error; err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Collection not found"})
	}

	if err := validateSingleContent(item.Data, collection.ContentType.Fields); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := config.DB.Create(&item).Error; err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create item"})
	}

	return c.Status(http.StatusOK).JSON(item)
}
