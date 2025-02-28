package models

import (
	"github.com/google/uuid"
)

type RoleType string

const (
	SuperAdmin   RoleType = "super_admin"
	ContentAdmin RoleType = "content_admin"
	Editor       RoleType = "editor"
	Viewer       RoleType = "viewer"
)

type Role struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name        string    `gorm:"unique;not null"`
	Type        RoleType  `gorm:"type:varchar(20);not null"`
	Description string
	Permissions []Permission `gorm:"many2many:role_permissions;"`
}
