package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type Content struct {
	ID            uint          `json:"id" gorm:"primaryKey;autoIncrement"`
	ContentTypeID uint          `json:"content_type_id"`
	ContentType   ContentType   `json:"content_type"`
	Data          JSON          `json:"data" gorm:"type:jsonb"`
	IsCollection  bool          `json:"is_collection" gorm:"default:false"`
	Items         []ContentItem `json:"items,omitempty" gorm:"foreignKey:CollectionID"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}

type ContentItem struct {
	ID           uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	CollectionID uint      `json:"collection_id"`
	Collection   Content   `json:"-" gorm:"foreignKey:CollectionID"`
	Data         JSON      `json:"data" gorm:"type:jsonb"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type JSON map[string]interface{}

func (j *JSON) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	*j = make(map[string]interface{})

	if len(bytes) == 0 {
		return nil
	}

	if err := json.Unmarshal(bytes, &j); err != nil {
		return err
	}

	return nil
}

func (j JSON) Value() (driver.Value, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return json.Marshal(j)
}
