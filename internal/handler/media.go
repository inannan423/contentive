package handler

import (
	"contentive/internal/database"
	"contentive/internal/logger"
	"contentive/internal/models"
	"contentive/internal/storage"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type MediaQuery struct {
	Page     int    `query:"page"`
	PageSize int    `query:"page_size"`
	Type     string `query:"type"`
	Search   string `query:"search"`
}

func UploadMedia(c *fiber.Ctx) error {
	// Get user from context
	var userID uuid.UUID
	if adminUser, ok := c.Locals("user").(models.AdminUser); ok {
		userID = adminUser.ID
	} else if apiUser, ok := c.Locals("user").(models.APIUser); ok {
		userID = apiUser.ID
	} else {
		logger.Error("Invalid user type in context")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Invalid user type",
		})
	}

	file, err := c.FormFile("file")
	if err != nil {
		logger.Error("Failed to get file: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No file uploaded",
		})
	}

	contentType := file.Header.Get("Content-Type")
	if contentType == "" || contentType == "application/octet-stream" || contentType == "Other" {
		ext := strings.ToLower(filepath.Ext(file.Filename))
		switch ext {
		case ".jpg", ".jpeg":
			contentType = "image/jpeg"
		case ".png":
			contentType = "image/png"
		case ".gif":
			contentType = "image/gif"
		case ".pdf":
			contentType = "application/pdf"
		case ".doc", ".docx":
			contentType = "application/msword"
		case ".xls", ".xlsx":
			contentType = "application/vnd.ms-excel"
		case ".mp4":
			contentType = "video/mp4"
		case ".mp3":
			contentType = "audio/mpeg"
		default:
			contentType = "application/octet-stream"
		}
	}

	mediaType := getMediaType(contentType)
	if mediaType == "" {
		logger.Error("Invalid file type: %s", contentType)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid file type",
		})
	}

	storageProvider := storage.GetStorageProvider()
	url, err := storageProvider.Upload(file, "media")
	if err != nil {
		logger.Error("Failed to upload file: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to upload file",
		})
	}

	media := models.Media{
		Name:      file.Filename,
		Type:      mediaType,
		MimeType:  contentType,
		Size:      file.Size,
		Path:      url,
		URL:       url,
		CreatedBy: userID,
	}

	if err := database.DB.Create(&media).Error; err != nil {
		logger.Error("Failed to create media record: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create media record",
		})
	}

	// Update logging to handle both user types
	if adminUser, ok := c.Locals("user").(models.AdminUser); ok {
		logger.AdminAction(adminUser.ID, adminUser.Name, "UPLOAD_MEDIA", "Uploaded media: "+media.Name)
	} else {
		logger.Info("API user %s uploaded media: %s", userID, media.Name)
	}

	return c.Status(fiber.StatusCreated).JSON(media)
}

func GetMedia(c *fiber.Ctx) error {
	id := c.Params("id")

	var media models.Media
	if err := database.DB.First(&media, "id = ?", id).Error; err != nil {
		logger.Error("Media not found: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Media not found",
		})
	}

	return c.JSON(media)
}

func DeleteMedia(c *fiber.Ctx) error {
	currentUser := c.Locals("user").(models.AdminUser)
	id := c.Params("id")

	var media models.Media
	if err := database.DB.First(&media, "id = ?", id).Error; err != nil {
		logger.Error("Media not found: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Media not found",
		})
	}

	storageProvider := storage.GetStorageProvider()
	if err := storageProvider.Delete(media.Path); err != nil {
		logger.Error("Failed to delete file: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete file",
		})
	}

	if err := database.DB.Delete(&media).Error; err != nil {
		logger.Error("Failed to delete media record: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete media record",
		})
	}

	logger.AdminAction(currentUser.ID, currentUser.Name, "DELETE_MEDIA", "Deleted media: "+media.Name)

	return c.SendStatus(fiber.StatusNoContent)
}

func ListMedia(c *fiber.Ctx) error {
	var query MediaQuery
	if err := c.QueryParser(&query); err != nil {
		logger.Error("Failed to parse query: %v", err)
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

	var media []models.Media
	var total int64
	db := database.DB.Model(&models.Media{})

	if query.Search != "" {
		db = db.Where("name ILIKE ?", "%"+query.Search+"%")
	}
	if query.Type != "" {
		db = db.Where("type = ?", query.Type)
	}

	if err := db.Count(&total).Error; err != nil {
		logger.Error("Failed to count media: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to count media",
		})
	}

	offset := (query.Page - 1) * query.PageSize
	if err := db.Offset(offset).Limit(query.PageSize).Find(&media).Error; err != nil {
		logger.Error("Failed to fetch media: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch media",
		})
	}

	return c.JSON(fiber.Map{
		"data": media,
		"meta": fiber.Map{
			"total":     total,
			"page":      query.Page,
			"page_size": query.PageSize,
		},
	})
}

// getMediaType returns the media type for the given MIME type
func getMediaType(mimeType string) models.MediaType {
	fmt.Printf("getMediaType: %s\n", mimeType)
	switch {
	case strings.HasPrefix(mimeType, "image/"):
		return models.MediaTypeImage
	case strings.HasPrefix(mimeType, "video/"):
		return models.MediaTypeVideo
	case strings.HasPrefix(mimeType, "audio/"):
		return models.MediaTypeAudio
	case strings.HasPrefix(mimeType, "application/") || strings.HasPrefix(mimeType, "text/"):
		return models.MediaTypeFile
	default:
		return ""
	}
}
