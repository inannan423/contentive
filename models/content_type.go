package models

type ContentTypeEnum string

const (
	Single     ContentTypeEnum = "single"
	Collection ContentTypeEnum = "collection"
)

// IsValid checks if the ContentTypeEnum value is valid.
func (t ContentTypeEnum) IsValid() bool {
	switch t {
	case Single, Collection:
		return true
	}
	return false
}

type ContentType struct {
	ID     uint            `json:"id" gorm:"primaryKey;autoIncrement"`
	Slug   string          `json:"slug" gorm:"unique"`
	Type   ContentTypeEnum `json:"type" gorm:"type:varchar(20)"`
	Name   string          `json:"name" gorm:"unique"`
	Fields []Field         `json:"fields" gorm:"foreignKey:ContentTypeID"`
}
