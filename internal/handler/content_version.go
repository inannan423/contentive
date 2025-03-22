package handler

import (
	"contentive/internal/database"
	"contentive/internal/logger"
	"contentive/internal/models"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// ListContentVersions returns a list of content versions for a given content entry
func ListContentVersions(c *fiber.Ctx) error {
	contentID := c.Params("content_id")

	var contentEntry models.ContentEntry
	if err := database.DB.Where("id = ?", contentID).First(&contentEntry).Error; err != nil {
		logger.Error("Content not found: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Content not found",
		})
	}

	var versions []models.ContentVersion
	if err := database.DB.Where("content_entry_id = ?", contentID).
		Order("version DESC").
		Find(&versions).Error; err != nil {
		logger.Error("Failed to fetch content versions: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch content versions",
		})
	}

	return c.JSON(versions)
}

// GetContentVersion returns a specific version of a content entry
func GetContentVersion(c *fiber.Ctx) error {
	contentID := c.Params("content_id")
	versionStr := c.Params("version")

	version, err := strconv.Atoi(versionStr)
	if err != nil {
		logger.Error("Invalid version number: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid version number",
		})
	}

	var contentVersion models.ContentVersion
	if err := database.DB.Where("content_entry_id = ? AND version = ?", contentID, version).
		First(&contentVersion).Error; err != nil {
		logger.Error("Content version not found: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Content version not found",
		})
	}

	return c.JSON(contentVersion)
}

func RestoreContentVersion(c *fiber.Ctx) error {
	contentID := c.Params("content_id")
	versionStr := c.Params("version")

	version, err := strconv.Atoi(versionStr)
	if err != nil {
		logger.Error("Invalid version number: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid version number",
		})
	}

	// Use transaction to ensure atomicity
	tx := database.DB.Begin()
	if tx.Error != nil {
		logger.Error("Failed to start transaction: %v", tx.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	// Get content entry
	var contentEntry models.ContentEntry
	if err := tx.Where("id = ?", contentID).First(&contentEntry).Error; err != nil {
		tx.Rollback()
		logger.Error("Content not found: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Content not found",
		})
	}

	// Get version to restore
	var versionToRestore models.ContentVersion
	if err := tx.Where("content_entry_id = ? AND version = ?", contentID, version).
		First(&versionToRestore).Error; err != nil {
		tx.Rollback()
		logger.Error("Content version not found: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Content version not found",
		})
	}

	// Update content entry data
	contentEntry.Data = versionToRestore.Data

	// Get current highest version number
	var maxVersion struct {
		MaxVersion int
	}
	if err := tx.Model(&models.ContentVersion{}).
		Select("MAX(version) as max_version").
		Where("content_entry_id = ?", contentID).
		Scan(&maxVersion).Error; err != nil {
		tx.Rollback()
		logger.Error("Failed to get highest version number: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	newVersionNumber := maxVersion.MaxVersion + 1

	if err := tx.Save(&contentEntry).Error; err != nil {
		tx.Rollback()
		logger.Error("Failed to update content: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update content",
		})
	}

	// Create new version (based on restored version)
	var userID uuid.UUID
	var userType models.ContentEntryUserByType

	if adminUser, ok := c.Locals("user").(models.AdminUser); ok {
		userID = adminUser.ID
		userType = models.ContentEntryUserByTypeAdmin
		logger.AdminAction(
			adminUser.ID,
			adminUser.Name,
			"RESTORE_CONTENT_VERSION",
			fmt.Sprintf("Content %s restored to version %d", contentID, version),
		)
	} else if apiUser, ok := c.Locals("user").(models.APIUser); ok {
		userID = apiUser.ID
		userType = models.ContentEntryUserByTypeAPI
		logger.APIAction(
			apiUser.ID,
			apiUser.Name,
			"RESTORE_CONTENT_VERSION",
			fmt.Sprintf("Content %s restored to version %d", contentID, version),
		)
	} else {
		tx.Rollback()
		logger.Error("User not found in context")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	// Update content entry updater information
	contentEntry.UpdatedBy = &userID
	contentEntry.UpdatedByType = userType

	if err := tx.Save(&contentEntry).Error; err != nil {
		tx.Rollback()
		logger.Error("Failed to update content metadata: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update content",
		})
	}

	newVersion := models.ContentVersion{
		ContentEntryID: contentEntry.ID,
		Version:        newVersionNumber,
		Data:           versionToRestore.Data,
		CreatedByID:    &userID,
	}

	if err := tx.Create(&newVersion).Error; err != nil {
		tx.Rollback()
		logger.Error("Failed to create content version: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create content version",
		})
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		logger.Error("Failed to commit transaction: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("Content restored to version %d", version),
		"content": contentEntry,
		"version": newVersion,
	})
}

func CompareContentVersions(c *fiber.Ctx) error {
	contentID := c.Params("content_id")
	v1Str := c.Query("v1")
	v2Str := c.Query("v2")

	if v1Str == "" || v2Str == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Both version parameters (v1 and v2) are required",
		})
	}

	v1, err := strconv.Atoi(v1Str)
	if err != nil {
		logger.Error("Invalid version number v1: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid version number for v1",
		})
	}

	v2, err := strconv.Atoi(v2Str)
	if err != nil {
		logger.Error("Invalid version number v2: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid version number for v2",
		})
	}

	var version1 models.ContentVersion
	if err := database.DB.Where("content_entry_id = ? AND version = ?", contentID, v1).
		First(&version1).Error; err != nil {
		logger.Error("Version 1 not found: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Version 1 not found",
		})
	}

	var version2 models.ContentVersion
	if err := database.DB.Where("content_entry_id = ? AND version = ?", contentID, v2).
		First(&version2).Error; err != nil {
		logger.Error("Version 2 not found: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Version 2 not found",
		})
	}

	// Parse data from both versions
	var data1 map[string]interface{}
	var data2 map[string]interface{}

	if err := json.Unmarshal(version1.Data, &data1); err != nil {
		logger.Error("Failed to unmarshal version 1 data: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	if err := json.Unmarshal(version2.Data, &data2); err != nil {
		logger.Error("Failed to unmarshal version 2 data: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	// Calculate differences
	differences := calculateDifferences(data1, data2)

	return c.JSON(fiber.Map{
		"v1":          v1,
		"v2":          v2,
		"differences": differences,
		"v1_data":     data1,
		"v2_data":     data2,
	})
}

func calculateDifferences(data1, data2 map[string]interface{}) map[string]interface{} {
	differences := make(map[string]interface{})

	// Check fields that exist in data1 but not in data2 or have different values
	for key, value1 := range data1 {
		if value2, exists := data2[key]; !exists {
			differences[key] = map[string]interface{}{
				"action":    "removed",
				"old_value": value1,
			}
		} else if !compareValues(value1, value2) {
			differences[key] = map[string]interface{}{
				"action":    "changed",
				"old_value": value1,
				"new_value": value2,
			}
		}
	}

	// Check fields that exist in data2 but not in data1
	for key, value2 := range data2 {
		if _, exists := data1[key]; !exists {
			differences[key] = map[string]interface{}{
				"action":    "added",
				"new_value": value2,
			}
		}
	}

	return differences
}

// Compare if two values are equal
func compareValues(v1, v2 interface{}) bool {
	// Simple comparison, can be extended to more complex deep comparison if needed
	return fmt.Sprintf("%v", v1) == fmt.Sprintf("%v", v2)
}

// CreateContentVersion creates a new version of a content entry manually
func CreateContentVersion(c *fiber.Ctx) error {
	contentID := c.Params("content_id")

	var input struct {
		Comment string                 `json:"comment"`
		Data    map[string]interface{} `json:"data"`
		Status  string                 `json:"status"` // draft or published
	}

	if err := c.BodyParser(&input); err != nil {
		logger.Error("Error parsing request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate status
	if input.Status != "" && input.Status != "draft" && input.Status != "published" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid status value, must be either 'draft' or 'published'",
		})
	}

	// Use transaction to ensure atomicity
	tx := database.DB.Begin()
	if tx.Error != nil {
		logger.Error("Failed to start transaction: %v", tx.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	// Get content entry
	var contentEntry models.ContentEntry
	if err := tx.Where("id = ?", contentID).First(&contentEntry).Error; err != nil {
		tx.Rollback()
		logger.Error("Content not found: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Content not found",
		})
	}

	// Get current highest version number
	var maxVersion struct {
		MaxVersion int
	}
	if err := tx.Model(&models.ContentVersion{}).
		Select("MAX(version) as max_version").
		Where("content_entry_id = ?", contentID).
		Scan(&maxVersion).Error; err != nil {
		tx.Rollback()
		logger.Error("Failed to get max version: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	newVersionNumber := maxVersion.MaxVersion + 1
	if newVersionNumber == 0 {
		newVersionNumber = 1 // If no versions exist, start from 1
	}

	// Prepare data
	var dataJSON datatypes.JSON
	if len(input.Data) > 0 {
		// if input.Data is not empty, use it
		jsonBytes, err := json.Marshal(input.Data)
		if err != nil {
			tx.Rollback()
			logger.Error("Failed to marshal data: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to process data",
			})
		}
		dataJSON = datatypes.JSON(jsonBytes)
	} else {
		// if input.Data is empty, use the current data
		dataJSON = contentEntry.Data
	}

	// Create new version
	var userID uuid.UUID

	if adminUser, ok := c.Locals("user").(models.AdminUser); ok {
		userID = adminUser.ID
		logger.AdminAction(
			adminUser.ID,
			adminUser.Name,
			"CREATE_CONTENT_VERSION",
			fmt.Sprintf("Created version %d for content %s", newVersionNumber, contentID),
		)
	} else if apiUser, ok := c.Locals("user").(models.APIUser); ok {
		userID = apiUser.ID
		logger.APIAction(
			apiUser.ID,
			apiUser.Name,
			"CREATE_CONTENT_VERSION",
			fmt.Sprintf("Created version %d for content %s", newVersionNumber, contentID),
		)
	} else {
		tx.Rollback()
		logger.Error("User not found in context")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	newVersion := models.ContentVersion{
		ContentEntryID: contentEntry.ID,
		Version:        newVersionNumber,
		Data:           dataJSON,
		CreatedByID:    &userID,
		Comment:        input.Comment,
		Status:         input.Status,
	}

	if err := tx.Create(&newVersion).Error; err != nil {
		tx.Rollback()
		logger.Error("Failed to create content version: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create content version",
		})
	}

	// Update content entry updater information
	contentEntry.CurrentVersion = newVersionNumber
	if err := tx.Save(&contentEntry).Error; err != nil {
		tx.Rollback()
		logger.Error("Failed to update content entry: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update content entry",
		})
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		logger.Error("Failed to commit transaction: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("Created version %d for content", newVersionNumber),
		"version": newVersion,
	})
}

// DeleteContentVersion deletes a specific version of a content entry
func DeleteContentVersion(c *fiber.Ctx) error {
	contentID := c.Params("content_id")
	versionStr := c.Params("version")

	version, err := strconv.Atoi(versionStr)
	if err != nil {
		logger.Error("Invalid version number: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid version number",
		})
	}

	tx := database.DB.Begin()
	if tx.Error != nil {
		logger.Error("Failed to start transaction: %v", tx.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	var contentEntry models.ContentEntry
	if err := tx.Where("id = ?", contentID).First(&contentEntry).Error; err != nil {
		tx.Rollback()
		logger.Error("Content not found: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Content not found",
		})
	}

	var versionToDelete models.ContentVersion
	if err := tx.Where("content_entry_id = ? AND version = ?", contentID, version).
		First(&versionToDelete).Error; err != nil {
		tx.Rollback()
		logger.Error("Content version not found: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Content version not found",
		})
	}

	var maxVersion struct {
		MaxVersion int
	}
	if err := tx.Model(&models.ContentVersion{}).
		Select("MAX(version) as max_version").
		Where("content_entry_id = ?", contentID).
		Scan(&maxVersion).Error; err != nil {
		tx.Rollback()
		logger.Error("Failed to get max version: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	if version == maxVersion.MaxVersion {
		tx.Rollback()
		logger.Error("Cannot delete the latest version")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot delete the latest version",
		})
	}

	var versionCount int64
	if err := tx.Model(&models.ContentVersion{}).
		Where("content_entry_id = ?", contentID).
		Count(&versionCount).Error; err != nil {
		tx.Rollback()
		logger.Error("Failed to count versions: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	if versionCount <= 1 {
		tx.Rollback()
		logger.Error("Cannot delete the only version")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot delete the only version",
		})
	}

	if err := tx.Delete(&versionToDelete).Error; err != nil {
		tx.Rollback()
		logger.Error("Failed to delete content version: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete content version",
		})
	}

	currentUser := c.Locals("user").(models.AdminUser)
	logger.AdminAction(
		currentUser.ID,
		currentUser.Name,
		"DELETE_CONTENT_VERSION",
		fmt.Sprintf("Deleted version %d for content %s", version, contentID),
	)

	if err := tx.Commit().Error; err != nil {
		logger.Error("Failed to commit transaction: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// PublishContentVersion publishes a specific version of a content entry
func PublishContentVersion(c *fiber.Ctx) error {
	contentID := c.Params("content_id")
	versionStr := c.Params("version")

	version, err := strconv.Atoi(versionStr)
	if err != nil {
		logger.Error("Invalid version number: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid version number",
		})
	}

	tx := database.DB.Begin()
	if tx.Error != nil {
		logger.Error("Failed to start transaction: %v", tx.Error)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	var contentEntry models.ContentEntry
	if err := tx.Where("id = ?", contentID).First(&contentEntry).Error; err != nil {
		tx.Rollback()
		logger.Error("Content not found: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Content not found",
		})
	}

	var versionToPublish models.ContentVersion
	if err := tx.Where("content_entry_id = ? AND version = ?", contentID, version).
		First(&versionToPublish).Error; err != nil {
		tx.Rollback()
		logger.Error("Content version not found: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Content version not found",
		})
	}

	contentEntry.Data = versionToPublish.Data
	contentEntry.IsPublished = true
	now := time.Now()
	contentEntry.PublishedAt = &now

	var userID uuid.UUID

	if adminUser, ok := c.Locals("user").(models.AdminUser); ok {
		userID = adminUser.ID
		contentEntry.PublishedBy = &userID
		logger.AdminAction(
			adminUser.ID,
			adminUser.Name,
			"PUBLISH_CONTENT_VERSION",
			fmt.Sprintf("Published version %d for content %s", version, contentID),
		)
	} else if apiUser, ok := c.Locals("user").(models.APIUser); ok {
		userID = apiUser.ID
		contentEntry.PublishedBy = &userID
		logger.APIAction(
			apiUser.ID,
			apiUser.Name,
			"PUBLISH_CONTENT_VERSION",
			fmt.Sprintf("Published version %d for content %s", version, contentID),
		)
	} else {
		tx.Rollback()
		logger.Error("User not found in context")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	if err := tx.Save(&contentEntry).Error; err != nil {
		tx.Rollback()
		logger.Error("Failed to update content: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update content",
		})
	}

	if err := tx.Commit().Error; err != nil {
		logger.Error("Failed to commit transaction: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	return c.JSON(fiber.Map{
		"message": fmt.Sprintf("Published version %d for content", version),
		"content": contentEntry,
	})
}

// GetContentVersionHistory gets the history of a specific content entry
func GetContentVersionHistory(c *fiber.Ctx) error {
	contentID := c.Params("content_id")

	var contentEntry models.ContentEntry
	if err := database.DB.Where("id = ?", contentID).First(&contentEntry).Error; err != nil {
		logger.Error("Content not found: %v", err)
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Content not found",
		})
	}

	type VersionHistory struct {
		ID          uuid.UUID  `json:"id"`
		Version     int        `json:"version"`
		CreatedAt   time.Time  `json:"created_at"`
		CreatedByID *uuid.UUID `json:"created_by_id"`
		CreatorName string     `json:"creator_name"`
		CreatorType string     `json:"creator_type"`
	}

	var versions []models.ContentVersion
	if err := database.DB.Where("content_entry_id = ?", contentID).
		Order("version DESC").
		Find(&versions).Error; err != nil {
		logger.Error("Failed to fetch content versions: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch content versions",
		})
	}

	var history []VersionHistory
	for _, v := range versions {
		var creatorName string
		var creatorType string

		if v.CreatedByID != nil {
			var adminUser models.AdminUser
			if err := database.DB.Where("id = ?", v.CreatedByID).First(&adminUser).Error; err == nil {
				creatorName = adminUser.Name
				creatorType = "admin"
			} else {
				var apiUser models.APIUser
				if err := database.DB.Where("id = ?", v.CreatedByID).First(&apiUser).Error; err == nil {
					creatorName = apiUser.Name
					creatorType = "api"
				}
			}
		}

		history = append(history, VersionHistory{
			ID:          v.ID,
			Version:     v.Version,
			CreatedAt:   v.CreatedAt,
			CreatedByID: v.CreatedByID,
			CreatorName: creatorName,
			CreatorType: creatorType,
		})
	}

	return c.JSON(fiber.Map{
		"content_id": contentID,
		"history":    history,
	})
}
