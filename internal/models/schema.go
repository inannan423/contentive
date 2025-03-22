package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type FieldType string

const (
	FieldTypeText      FieldType = "text"
	FieldTypeTextarea  FieldType = "textarea"
	FieldTypeNumber    FieldType = "number"
	FieldTypeDate      FieldType = "date"     // Date only
	FieldTypeDateTime  FieldType = "datetime" // Date and time
	FieldTypeBoolean   FieldType = "boolean"
	FieldTypeRelation  FieldType = "relation"
	FieldTypeMedia     FieldType = "media"
	FieldTypeMediaList FieldType = "media_list"
	FieldTypeSelect    FieldType = "select"
	FieldTypeRichText  FieldType = "richtext"
	FieldTypeEmail     FieldType = "email"
	FieldTypePassword  FieldType = "password"
)

type SchemaType string

const (
	SchemaTypeSingle SchemaType = "single"
	SchemaTypeList   SchemaType = "list"
)

// FieldDefinition is a struct that represents the definition of a field
type FieldDefinition struct {
	ID       uuid.UUID              `json:"id" gorm:"type:uuid;"`
	Name     string                 `json:"name"`
	Type     FieldType              `json:"type"`
	Required bool                   `json:"required"`
	Options  map[string]interface{} `json:"options,omitempty"` // Extend the field definition with more options, like min, max, relation, etc.
}

// Schema is a struct that represents the schema of a content type
type Schema struct {
	ID        uuid.UUID      `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name      string         `json:"name" gorm:"unique;not null"`
	Type      SchemaType     `json:"type" gorm:"type:varchar(10);not null"`
	Slug      string         `json:"slug" gorm:"unique;not null"`
	Fields    datatypes.JSON `json:"fields" gorm:"type:jsonb;not null"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
}

// Reserved Words List: Schema name cannot be one of the following:
var reservedWords = []string{
	"media", "media_list",
}

// ValidateFields validates the schema fields with type-specific rules.
func (s *Schema) ValidateFields() error {
	var fields []FieldDefinition
	if err := json.Unmarshal(s.Fields, &fields); err != nil {
		return fmt.Errorf("invalid fields format: %v", err)
	}

	// Check for duplicate field names.
	fieldNames := make(map[string]bool)
	for _, field := range fields {
		if field.Name == "" {
			return errors.New("field name cannot be empty")
		}
		if fieldNames[field.Name] {
			return fmt.Errorf("duplicate field name: %s", field.Name)
		}
		// If the field name is a reserved word, return an error.
		for _, word := range reservedWords {
			if field.Name == word {
				return fmt.Errorf("field name '%s' is a reserved word", field.Name)
			}
		}
		fieldNames[field.Name] = true

		// Validate that the field type is allowed.
		allowedTypes := []FieldType{
			FieldTypeText, FieldTypeTextarea, FieldTypeNumber,
			FieldTypeDate, FieldTypeDateTime, FieldTypeBoolean,
			FieldTypeRelation, FieldTypeMedia, FieldTypeSelect,
			FieldTypeRichText, FieldTypeEmail, FieldTypePassword,
			FieldTypeMediaList,
		}
		validType := false
		for _, t := range allowedTypes {
			if field.Type == t {
				validType = true
				break
			}
		}
		if !validType {
			return fmt.Errorf("invalid field type: %s for field %s", field.Type, field.Name)
		}

		// Type-specific validations.
		switch field.Type {
		// Validate text-based fields.
		case FieldTypeText, FieldTypeTextarea:
			// If maxLength is provided, it must be a positive integer.
			if maxLength, exists := field.Options["maxLength"]; exists {
				v, ok := maxLength.(float64)
				if !ok || v <= 0 || v != float64(int(v)) {
					return fmt.Errorf("%s field %s: 'maxLength' must be a positive integer", field.Type, field.Name)
				}
				// If minLength is provided, it must be a non-negative integer and not greater than maxLength.
				if minLength, exists := field.Options["minLength"]; exists {
					vMin, ok := minLength.(float64)
					if !ok || vMin < 0 || vMin != float64(int(vMin)) {
						return fmt.Errorf("%s field %s: 'minLength' must be a non-negative integer", field.Type, field.Name)
					}
					if vMin > v {
						return fmt.Errorf("%s field %s: 'minLength' cannot be greater than 'maxLength'", field.Type, field.Name)
					}
				}
			}

		// Validate number fields.
		case FieldTypeNumber:
			var minVal, maxVal float64
			var hasMin, hasMax bool
			// Check the 'min' option.
			if min, exists := field.Options["min"]; exists {
				v, ok := min.(float64)
				if !ok {
					return fmt.Errorf("number field %s: 'min' must be a number", field.Name)
				}
				minVal = v
				hasMin = true
			}
			// Check the 'max' option.
			if max, exists := field.Options["max"]; exists {
				v, ok := max.(float64)
				if !ok {
					return fmt.Errorf("number field %s: 'max' must be a number", field.Name)
				}
				maxVal = v
				hasMax = true
			}
			// If both min and max are provided, ensure min <= max.
			if hasMin && hasMax && minVal > maxVal {
				return fmt.Errorf("number field %s: 'min' cannot be greater than 'max'", field.Name)
			}
			// Check the 'precision' option.
			if precision, exists := field.Options["precision"]; exists {
				v, ok := precision.(float64)
				if !ok || v < 0 || v != float64(int(v)) {
					return fmt.Errorf("number field %s: 'precision' must be a non-negative integer", field.Name)
				}
			}

		// Validate date and datetime fields.
		case FieldTypeDate, FieldTypeDateTime:
			// If format is provided, it must be a non-empty string.
			if format, exists := field.Options["format"]; exists {
				v, ok := format.(string)
				if !ok || v == "" {
					return fmt.Errorf("%s field %s: 'format' must be a non-empty string", field.Type, field.Name)
				}
			}

		// Validate boolean fields.
		case FieldTypeBoolean:
			// If a default value is provided, it must be a boolean.
			if def, exists := field.Options["default"]; exists {
				if _, ok := def.(bool); !ok {
					return fmt.Errorf("boolean field %s: 'default' must be a boolean", field.Name)
				}
			}

		// Validate select fields.
		case FieldTypeSelect:
			// 'choices' is required and must be a non-empty array.
			choices, ok := field.Options["choices"]
			if !ok {
				return fmt.Errorf("select field %s must have 'choices' option", field.Name)
			}
			choicesSlice, ok := choices.([]interface{})
			if !ok || len(choicesSlice) == 0 {
				return fmt.Errorf("select field %s must have a non-empty 'choices' array", field.Name)
			}
			// Ensure all choices are strings.
			for i, choice := range choicesSlice {
				if _, ok := choice.(string); !ok {
					return fmt.Errorf("select field %s: choice at index %d is not a string", field.Name, i)
				}
			}
			// If default value is provided, ensure it is one of the choices.
			if def, exists := field.Options["default"]; exists {
				defStr, ok := def.(string)
				if !ok {
					return fmt.Errorf("select field %s: default value must be a string", field.Name)
				}
				found := false
				for _, choice := range choicesSlice {
					if choice.(string) == defStr {
						found = true
						break
					}
				}
				if !found {
					return fmt.Errorf("select field %s: default value '%s' is not in choices", field.Name, defStr)
				}
			}

		// Validate relation fields.
		case FieldTypeRelation:
			// 'targetSchema' is required and must be a non-empty string.
			target, ok := field.Options["targetSchema"]
			if !ok {
				return fmt.Errorf("relation field %s must have 'targetSchema' option", field.Name)
			}
			targetStr, ok := target.(string)
			if !ok || targetStr == "" {
				return fmt.Errorf("relation field %s: 'targetSchema' must be a non-empty string", field.Name)
			}

			// Check if the target schema exists.
			if schemaValidator != nil {
				if err := schemaValidator.ValidateTargetSchema(targetStr); err != nil {
					return fmt.Errorf("relation field %s: %v", field.Name, err)
				}
			}

			// 'relationType' is required and must be a valid value.
			relationType, ok := field.Options["relationType"]
			if !ok {
				return fmt.Errorf("relation field %s must have 'relationType' option", field.Name)
			}
			relationTypeStr, ok := relationType.(string)
			if !ok || relationTypeStr == "" {
				return fmt.Errorf("relation field %s: 'relationType' must be a non-empty string", field.Name)
			}
			allowedRelationTypes := []string{"one-to-one", "one-to-many", "many-to-many", "many-to-one"}
			validRelType := false
			for _, rel := range allowedRelationTypes {
				if relationTypeStr == rel {
					validRelType = true
					break
				}
			}
			if !validRelType {
				return fmt.Errorf("relation field %s: invalid relationType '%s'", field.Name, relationTypeStr)
			}

		// Validate media fields.
		case FieldTypeMedia:
			// If allowedTypes is provided, it must be a non-empty array of strings.
			if allowedTypesVal, exists := field.Options["allowedTypes"]; exists {
				arr, ok := allowedTypesVal.([]interface{})
				if !ok || len(arr) == 0 {
					return fmt.Errorf("media field %s: 'allowedTypes' must be a non-empty array", field.Name)
				}
				for i, v := range arr {
					if _, ok := v.(string); !ok {
						return fmt.Errorf("media field %s: allowedTypes at index %d is not a string", field.Name, i)
					}
				}
			}
			// If maxSize is provided, it must be a positive integer.
			if maxSize, exists := field.Options["maxSize"]; exists {
				v, ok := maxSize.(float64)
				if !ok || v <= 0 || v != float64(int(v)) {
					return fmt.Errorf("media field %s: 'maxSize' must be a positive integer", field.Name)
				}
			}

		case FieldTypeMediaList:
			// If allowedTypes is provided, it must be a non-empty array of strings.
			if allowedTypesVal, exists := field.Options["allowedTypes"]; exists {
				arr, ok := allowedTypesVal.([]interface{})
				if !ok || len(arr) == 0 {
					return fmt.Errorf("media_list field %s: 'allowedTypes' must be a non-empty array", field.Name)
				}
				for i, v := range arr {
					if _, ok := v.(string); !ok {
						return fmt.Errorf("media_list field %s: allowedTypes at index %d is not a string", field.Name, i)
					}
				}
			}
			// If maxSize is provided, it must be a positive integer.
			if maxSize, exists := field.Options["maxSize"]; exists {
				v, ok := maxSize.(float64)
				if !ok || v <= 0 || v != float64(int(v)) {
					return fmt.Errorf("media_list field %s:'maxSize' must be a positive integer", field.Name)
				}
			}

		// Validate rich text fields.
		case FieldTypeRichText:
			// If default value is provided, it must be a string.
			if def, exists := field.Options["default"]; exists {
				if _, ok := def.(string); !ok {
					return fmt.Errorf("richtext field %s: 'default' must be a string", field.Name)
				}
			}

		// Validate email fields.
		case FieldTypeEmail:
			// If default value is provided, it must be a string.
			if def, exists := field.Options["default"]; exists {
				if _, ok := def.(string); !ok {
					return fmt.Errorf("email field %s: 'default' must be a string", field.Name)
				}
				// Email format validation.
				if !strings.Contains(def.(string), "@") {
					return fmt.Errorf("email field %s: 'default' must be a valid email address", field.Name)
				}
			}

		// Validate password fields.
		case FieldTypePassword:
			// Password fields should not have a default value.
			if _, exists := field.Options["default"]; exists {
				return fmt.Errorf("password field %s: default value is not allowed", field.Name)
			}
		}
	}

	return nil
}
