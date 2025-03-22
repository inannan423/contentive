package models

import (
	"time"

	"github.com/google/uuid"
)

type MediaType string

const (
	MediaTypeImage MediaType = "image"
	MediaTypeVideo MediaType = "video"
	MediaTypeAudio MediaType = "audio"
	MediaTypeFile  MediaType = "file"
)

type Media struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Name      string    `json:"name" gorm:"type:varchar(255);not null"`
	Type      MediaType `json:"type" gorm:"type:varchar(255);not null"`
	MimeType  string    `json:"mime_type" gorm:"type:varchar(255);not null"`
	Size      int64     `json:"size" gorm:"type:bigint;not null"`
	Path      string    `json:"path" gorm:"type:varchar(255);not null"`
	URL       string    `json:"url" gorm:"type:varchar(255);not null"`
	Width     *int      `json:"width,omitempty"`
	Height    *int      `json:"height,omitempty"`
	Duration  *int      `json:"duration,omitempty"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	CreatedBy uuid.UUID `json:"created_by" gorm:"type:uuid;not null"`
}
