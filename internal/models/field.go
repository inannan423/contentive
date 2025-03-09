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

type RelationTypeEnum string

const (
	OneToOne   RelationTypeEnum = "one_to_one"
	OneToMany  RelationTypeEnum = "one_to_many"
	ManyToOne  RelationTypeEnum = "many_to_one"
	ManyToMany RelationTypeEnum = "many_to_many"
)

type Field struct {
	ID            uuid.UUID         `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	ContentTypeID uuid.UUID         `gorm:"type:uuid;not null"`
	Label         string            `gorm:"not null"`
	Type          FieldTypeEnum     `gorm:"not null"`
	Required      bool              `gorm:"not null"`
	CreatedAt     time.Time         `gorm:"autoCreateTime"`
	UpdatedAt     time.Time         `gorm:"autoUpdateTime"`
	RelationType  *RelationTypeEnum `gorm:"type:relation_type_enum"`
	TargetTypeID  *uuid.UUID        `gorm:"type:uuid"`
	TargetType    *ContentType      `gorm:"foreignKey:TargetTypeID"`
}
