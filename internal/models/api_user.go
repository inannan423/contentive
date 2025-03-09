package models

import (
	"contentive/internal/logger"
	"gorm.io/gorm"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"

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
	ID          uuid.UUID      `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"unique"`
	Description string         `json:"description"`
	Token       string         `json:"token" gorm:"unique"`
	ExpireAt    *time.Time     `json:"expire_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	Scopes      pq.StringArray `json:"scopes" gorm:"type:text[]"`
	Status      APIUserStatus  `json:"status"`
}

// Claims is a struct that contains the claims for the JWT token
type Claims struct {
	UserID uuid.UUID `json:"user_id"`
	Name   string    `json:"name"`
	Scopes pq.StringArray
	Status string `json:"status"`
	jwt.RegisteredClaims
}

// GenerateAPIUserJWT generates a JWT token for the API user
func (u *APIUser) GenerateAPIUserJWT() (string, error) {
	// Only set the expiresAt if it is not nil
	var expiresAt *jwt.NumericDate
	if u.ExpireAt != nil {
		expiresAt = jwt.NewNumericDate(*u.ExpireAt)
	}

	claims := Claims{
		UserID: u.ID,
		Name:   u.Name,
		Scopes: u.Scopes,
		Status: string(u.Status),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: expiresAt,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	return SignAPIUserToken(claims)
}

// BeforeCreate is a GORM hook that is called before creating a new user
func (u *APIUser) BeforeCreate(*gorm.DB) error {
	u.ID = uuid.New()
	// Use GenerateAPIUserJWT to generate a new token
	token, err := u.GenerateAPIUserJWT()
	if err != nil {
		logger.Error("Failed to generate token: %v", err)
		log.Fatal(err)
	}
	u.Token = token
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	return nil
}

var Secret string

func SetSecret(secret string) {
	Secret = secret
}

func SignAPIUserToken(claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(Secret))
}

func ValidateAPIUserToken(tokenString string, claims jwt.Claims) error {
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(Secret), nil
	})
	return err
}
