package routes

import (
	"contentive/internal/handlers"
	"contentive/internal/middlewares"
	"contentive/internal/models"

	"github.com/gofiber/fiber/v2"
)

func RegisterContentTypeRoutes(app *fiber.App) {
	contentTypes := app.Group("/api/content-types")

	// Protect all routes in this group
	contentTypes.Use(middlewares.AuthMiddleware())

	// POST /api/content-types - Create a new content type
	contentTypes.Post("/",
		middlewares.RequirePermission(models.CreateContentType),
		middlewares.ValidateContentType(),
		handlers.CreateContentType,
	)

	// GET /api/content-types - Get all content types
	contentTypes.Get("/",
		middlewares.RequirePermission(models.ReadContentType),
		handlers.GetAllContentTypes,
	)

	// GET /api/content-types/:contentTypeId - Get a content type by ID
	contentTypes.Get("/:identifier",
		middlewares.RequirePermission(models.ReadContentType),
		handlers.GetContentType,
	)

	contentTypes.Put("/:identifier",
		middlewares.RequirePermission(models.UpdateContentType),
		middlewares.ValidateContentType(),
		handlers.UpdateContentType,
	)

	contentTypes.Post("/:identifier/fields",
		middlewares.RequirePermission(models.UpdateContentType),
		middlewares.ValidateField(),
		handlers.AddField,
	)

	// PUT /api/content-types/:contentTypeId/fields/:id - Update a field in a content type
	contentTypes.Put("/:identifier/fields/:id",
		middlewares.RequirePermission(models.UpdateContentType),
		middlewares.ValidateField(),
		handlers.UpdateField,
	)

	// DELETE /api/content-types/:contentTypeId/fields/:id - Delete a field from a content type
	contentTypes.Delete("/:identifier/fields/:id",
		middlewares.RequirePermission(models.DeleteContentType),
		handlers.DeleteField,
	)
}
