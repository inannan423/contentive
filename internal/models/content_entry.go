package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type ContentEntry struct {
	ID            uuid.UUID      `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Slug          string         `gorm:"unique;not null"`
	ContentTypeID uuid.UUID      `gorm:"type:uuid;not null"`
	Data          datatypes.JSON `gorm:"type:jsonb"`
	CreatedAt     time.Time      `gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime"`
}
