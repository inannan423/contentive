package database

import (
	"contentive/internal/models"
	"fmt"
)

type dbSchemaValidator struct{}

func (v *dbSchemaValidator) ValidateTargetSchema(slug string) error {
	var schema models.Schema
	if err := DB.Where("slug = ?", slug).First(&schema).Error; err != nil {
		return fmt.Errorf("target schema '%s' does not exist", slug)
	}
	return nil
}

func InitSchemaValidator() {
	models.SetSchemaValidator(&dbSchemaValidator{})
}
