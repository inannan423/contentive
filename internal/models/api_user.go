package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

// APIUserScope is a type for API user scopes.
// ["blog:read", "blog:create", "blog:update", "blog:delete"]
type APIUserScope string

type APIUserStatus string

const (
	APIUserStatusActive   APIUserStatus = "active"
	APIUserStatusInactive APIUserStatus = "inactive"
	APIUserStatusExpired  APIUserStatus = "expired"
)

type APIUser struct {
	ID        uuid.UUID      `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"unique"`
	Token     string         `json:"token" gorm:"unique"`
	ExpireAt  *time.Time     `json:"expire_at"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	Scopes    pq.StringArray `json:"scopes" gorm:"type:text[]"`
	Status    APIUserStatus  `json:"status"`
}
