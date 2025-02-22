package models

type FieldType string

const (
	Text     FieldType = "text"
	Number   FieldType = "number"
	Boolean  FieldType = "boolean"
	DateTime FieldType = "datetime"
	Media    FieldType = "media"
)

func (t FieldType) IsValid() bool {
	switch t {
	case Text, Number, Boolean, DateTime, Media:
		return true
	}
	return false
}

type Field struct {
	ID            uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name          string    `json:"name" gorm:"not null"`
	Type          FieldType `json:"type" gorm:"type:varchar(20)"`
	Required      bool      `json:"required" gorm:"default:false"`
	ContentTypeID uint      `json:"content_type_id"`
	ContentType   ContentType
}