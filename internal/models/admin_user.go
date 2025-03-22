package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

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
	ID          uuid.UUID       `json:"id" gorm:"primaryKey;autoIncrement"`
	Name        string          `json:"name" gorm:"unique"`
	Email       string          `json:"email" gorm:"unique"`
	Password    string          `json:"-" gorm:"column:password"` // json:"-" means it will not be returned in the response
	Role        AdminUserRole   `json:"role"`
	Permissions pq.StringArray  `json:"permissions" gorm:"type:text[]"`
	Status      AdminUserStatus `json:"status"`
	CreatedAt   time.Time       `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time       `json:"updated_at" gorm:"autoUpdateTime"`
	LastLoginAt time.Time       `json:"last_login_at"`
}

// SetPermissions sets the permissions for the user according to the role
func (u *AdminUser) SetPermissions() {
	switch u.Role {
	case AdminUserRoleSuperAdmin:
		u.Permissions = []string{string(AdminUserPermissionAll)}
	case AdminUserRoleAdmin:
		u.Permissions = []string{string(AdminUserPermissionSchemeManagement), string(AdminUserPermissionContentRead), string(AdminUserPermissionContentEdit)}
	case AdminUserRoleEditor:
		u.Permissions = []string{string(AdminUserPermissionContentRead), string(AdminUserPermissionContentEdit)}
	case AdminUserRoleViewer:
		u.Permissions = []string{string(AdminUserPermissionContentRead)}
	}
}

// SetPassword sets the password for the user using bcrypt
func (u *AdminUser) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword checks if the provided password matches the user's password
func (u *AdminUser) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// IsSuperAdmin checks if the user is a super admin
func (u *AdminUser) IsSuperAdmin() bool {
	return u.Role == AdminUserRoleSuperAdmin
}

// BeforeCreate is a GORM hook that is called before creating a new user
func (u *AdminUser) BeforeCreate(tx *gorm.DB) error {
	u.ID = uuid.New()
	if u.Password != "" {
		if err := u.SetPassword(u.Password); err != nil {
			return err
		}
	}
	// if the role is not set, set it to viewer
	if u.Role == "" {
		u.Role = AdminUserRoleViewer
	}
	u.SetPermissions()
	return nil
}

// BeforeUpdate is a GORM hook that is called before updating a user
func (u *AdminUser) BeforeUpdate(tx *gorm.DB) error {
	if tx.Statement.Changed("Password") {
		if err := u.SetPassword(u.Password); err != nil {
			return err
		}
	}
	// if the role is changed, set the permissions
	if tx.Statement.Changed("Role") {
		u.SetPermissions()
	}
	return nil
}

// GetRoleLevel returns the level of the role
func (r AdminUserRole) GetRoleLevel() int {
	switch r {
	case AdminUserRoleSuperAdmin:
		return 4
	case AdminUserRoleAdmin:
		return 3
	case AdminUserRoleEditor:
		return 2
	case AdminUserRoleViewer:
		return 1
	default:
		return 0
	}
}
