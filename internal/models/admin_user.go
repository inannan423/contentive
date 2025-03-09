package models

import "time"

type AdminUserRole string

const (
	AdminUserRoleSuperAdmin AdminUserRole = "super_admin"
	AdminUserRoleAdmin      AdminUserRole = "admin"
	AdminUserRoleEditor     AdminUserRole = "editor"
	AdminUserRoleViewer     AdminUserRole = "viewer"
)

// AdminUserPermission: determines what the user can do
type AdminUserPermission string

const (
	AdminUserPermissionAll              AdminUserPermission = "all"
	AdminUserPermissionSchemeManagement AdminUserPermission = "schema:manage"
	AdminUserPermissionContentRead      AdminUserPermission = "content:read"
	AdminUserPermissionContentEdit      AdminUserPermission = "content:edit"
)

type AdminUserStatus string

const (
	AdminUserStatusActive   AdminUserStatus = "active"
	AdminUserStatusInactive AdminUserStatus = "inactive"
)

type AdminUser struct {
	ID          uint                  `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string                `json:"name" gorm:"unique"`
	Email       string                `json:"email" gorm:"unique"`
	Password    string                `json:"password"`
	Role        AdminUserRole         `json:"role"`
	Permissions []AdminUserPermission `json:"permissions" gorm:"type:varchar(255)[]"`
	Status      AdminUserStatus       `json:"status"`
	CreatedAt   time.Time             `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time             `json:"updated_at" gorm:"autoUpdateTime"`
	LastLoginAt time.Time             `json:"last_login_at"`
}
