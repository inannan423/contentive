package models

import "github.com/google/uuid"

type PermissionType string

const (
	// ContentType Permissions
	CreateContentType PermissionType = "create_content_type"
	ReadContentType   PermissionType = "read_content_type"
	UpdateContentType PermissionType = "update_content_type"
	DeleteContentType PermissionType = "delete_content_type"

	// Content Permissions
	CreateContent PermissionType = "create_content"
	ReadContent   PermissionType = "read_content"
	UpdateContent PermissionType = "update_content"
	DeleteContent PermissionType = "delete_content"

	// User Permissions
	ManageUsers   PermissionType = "manage_users"
	ManageRoles   PermissionType = "manage_roles"
	ViewAuditLogs PermissionType = "view_audit_logs"
)

type Permission struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	Name        string         `gorm:"unique;not null"`
	Type        PermissionType `gorm:"type:varchar(50);not null"`
	Description string
}
