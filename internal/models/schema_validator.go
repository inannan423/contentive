package models

type SchemaValidator interface {
	ValidateTargetSchema(slug string) error
}

var schemaValidator SchemaValidator

func SetSchemaValidator(validator SchemaValidator) {
	schemaValidator = validator
}
