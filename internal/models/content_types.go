package models

import (
	"time"

	"github.com/google/uuid"
)

type ContentTypeEnum string

const (
	Single     ContentTypeEnum = "single"
	Collection ContentTypeEnum = "collection"
)

type ContentType struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Name      string    `json:"name" gorm:"not null"`
	Type      string    `json:"type" gorm:"not null"`
	CreatedAt time.Time `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updatedAt" gorm:"autoUpdateTime"`
	Fields    []Field   `json:"fields" gorm:"foreignKey:ContentTypeID"`
}
