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
	ID            uint        `json:"id" gorm:"primaryKey;autoIncrement"`
	Slug          string      `json:"slug" gorm:"not null;uniqueIndex:idx_content_type_field_name"`
	Name          string      `json:"name" gorm:"not null;uniqueIndex:idx_content_type_field_name"`
	Type          FieldType   `json:"type" gorm:"type:varchar(20)"`
	Required      bool        `json:"required" gorm:"default:false"`
	ContentTypeID uint        `json:"content_type_id" gorm:"uniqueIndex:idx_content_type_field_name"`
	ContentType   ContentType `json:"-" gorm:"foreignKey:ContentTypeID"`
}
