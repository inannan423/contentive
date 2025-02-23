package routes

import (
	"contentive/internal/handlers"
	"contentive/internal/middlewares"

	"github.com/gofiber/fiber/v2"
)

func RegisterContentTypeRoutes(app *fiber.App) {
	contentTypes := app.Group("/api/content-types")

	// POST /api/content-types - Create a new content type
	contentTypes.Post("/",
		middlewares.ValidateContentType(),
		handlers.CreateContentType,
	)

	// GET /api/content-types - Get all content types
	contentTypes.Get("/", handlers.GetAllContentTypes)

	// GET /api/content-types/:contentTypeId - Get a content type by ID
	contentTypes.Get("/:contentTypeId", handlers.GetContentType)

	// PUT /api/content-types/:contentTypeId - Update a content type by ID
	contentTypes.Put("/:contentTypeId",
		middlewares.ValidateContentType(),
		handlers.UpdateContentType,
	)

	// POST /api/content-types/:contentTypeId/fields - Add a new field to a content type
	contentTypes.Post("/:contentTypeId/fields",
		middlewares.ValidateField(),
		handlers.AddField,
	)

	// PUT /api/content-types/:contentTypeId/fields/:id - Update a field in a content type
	contentTypes.Put("/:contentTypeId/fields/:id",
		middlewares.ValidateField(),
		handlers.UpdateField,
	)

	// DELETE /api/content-types/:contentTypeId/fields/:id - Delete a field from a content type
	contentTypes.Delete("/:contentTypeId/fields/:id", handlers.DeleteField)
}
