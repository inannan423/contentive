package models

import (
	"github.com/google/uuid"
)

type OperationType string

const (
	CreateOperation OperationType = "create"
	ReadOperation   OperationType = "read"
	UpdateOperation OperationType = "update"
	DeleteOperation OperationType = "delete"
)

type APIPermission struct {
	ID            uuid.UUID     `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	APIRoleID     uuid.UUID     `gorm:"type:uuid;not null"`
	ContentTypeID uuid.UUID     `gorm:"type:uuid;not null"`
	ContentType   ContentType   `gorm:"foreignKey:ContentTypeID"`
	Operation     OperationType `gorm:"type:varchar(20);not null"`
	Enabled       bool          `gorm:"default:false"`
}

type APIPermissionKey struct {
	APIRoleID     uuid.UUID
	ContentTypeID uuid.UUID
	Operation     OperationType
}
