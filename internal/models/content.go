package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type ContentEntryUserByType string

const (
	ContentEntryUserByTypeAdmin ContentEntryUserByType = "admin"
	ContentEntryUserByTypeAPI   ContentEntryUserByType = "api"
)

// ContentEntry represents a single entry of a content type
type ContentEntry struct {
	ID             uuid.UUID              `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Slug           string                 `json:"slug" gorm:"unique;not null"`
	ContentTypeID  uuid.UUID              `json:"content_type_id" gorm:"type:uuid;not null"`
	Data           datatypes.JSON         `json:"data" gorm:"type:jsonb"`
	IsPublished    bool                   `json:"is_published" gorm:"default:false"`
	PublishedAt    *time.Time             `json:"published_at"`
	CreatedAt      time.Time              `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time              `json:"updated_at" gorm:"autoUpdateTime"`
	PublishedBy    *uuid.UUID             `json:"published_by" gorm:"type:uuid"`
	CreatedByType  ContentEntryUserByType `json:"created_by_type" gorm:"not null"`
	UpdatedBy      *uuid.UUID             `json:"updated_by" gorm:"type:uuid"`
	UpdatedByType  ContentEntryUserByType `json:"updated_by_type" gorm:"not null"`
	CurrentVersion int                    `json:"current_version" gorm:"default:1"`
	Versions       []ContentVersion       `json:"versions,omitempty" gorm:"foreignKey:ContentEntryID"`
	Status         string                 `json:"status" gorm:"type:varchar(20);default:'draft'"`
}

// ContentVersion represents a version of a content entry
type ContentVersion struct {
	ID             uuid.UUID      `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	ContentEntryID uuid.UUID      `json:"content_entry_id" gorm:"type:uuid;not null"`
	Version        int            `json:"version" gorm:"not null"`
	Data           datatypes.JSON `json:"data" gorm:"type:jsonb;not null"`
	CreatedByID    *uuid.UUID     `json:"created_by_id" gorm:"type:uuid"`
	CreatedAt      time.Time      `json:"created_at" gorm:"autoCreateTime"`
	Comment        string         `json:"comment" gorm:"type:text"`
	Status         string         `json:"status" gorm:"type:varchar(20);default:'draft';not null"`
}
