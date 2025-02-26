package models

import (
	"time"

	"github.com/google/uuid"
)

type ActionType string

const (
	Create ActionType = "create"
	Update ActionType = "update"
	Delete ActionType = "delete"
	Login  ActionType = "login"
)

type AuditLog struct {
	ID        uuid.UUID  `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	UserID    uuid.UUID  `gorm:"type:uuid;not null"`
	Action    ActionType `gorm:"type:varchar(20);not null"`
	Resource  string     `gorm:"not null"`
	Details   string
	CreatedAt time.Time `gorm:"autoCreateTime"`
}
