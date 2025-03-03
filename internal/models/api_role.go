package models

import (
	"time"

	"github.com/google/uuid"
)

type APIRoleType string

const (
	AuthenticatedUser APIRoleType = "authenticated_user"
	PublicUser        APIRoleType = "public_user"
	CustomAPIRole     APIRoleType = "custom"
)

type APIRole struct {
	ID          uuid.UUID   `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Name        string      `gorm:"unique;not null"`
	Type        APIRoleType `gorm:"type:varchar(50);not null"`
	Description string
	APIKey      string          `gorm:"unique"`
	IsSystem    bool            `gorm:"default:false"` // System roles cannot be deleted
	ExpiresAt   *time.Time      `gorm:"default:null"`
	CreatedAt   time.Time       `gorm:"autoCreateTime"`
	UpdatedAt   time.Time       `gorm:"autoUpdateTime"`
	Permissions []APIPermission `gorm:"foreignKey:APIRoleID"`
}
