package models

import (
	"time"

	"github.com/google/uuid"
)

type FieldTypeEnum string

const (
	Text     FieldTypeEnum = "text"
	RichText FieldTypeEnum = "rich_text"
	Number   FieldTypeEnum = "number"
	Date     FieldTypeEnum = "date"
	Boolean  FieldTypeEnum = "boolean"
	Enum     FieldTypeEnum = "enum"
	Relation FieldTypeEnum = "relation"
)

type Field struct {
	ID            uuid.UUID     `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	ContentTypeID uuid.UUID     `gorm:"type:uuid;not null"`
	Label         string        `gorm:"not null"`
	Type          FieldTypeEnum `gorm:"not null"`
	CreatedAt     time.Time     `gorm:"autoCreateTime"`
	UpdatedAt     time.Time     `gorm:"autoUpdateTime"`
}
