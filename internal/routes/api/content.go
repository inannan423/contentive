package apiroutes

import (
	"contentive/internal/handler"
	"contentive/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterAPIContentRoutes(app *fiber.App) {
	api := app.Group("/api")
	content := api.Group("/content")

	// All API routes require API token authentication
	content.Use(middleware.AuthenticateAPIUserToken())

	// Routes for content operations
	// Create content - requires {schema}:create scope
	content.Post("/schema/:schema_slug",
		middleware.GetSchemaFromSlug(),
		middleware.RequireSchemaScope("create"),
		handler.CreateContent,
	)

	// Get content - requires {schema}:read scope
	content.Get("/schema/:schema_slug",
		middleware.GetSchemaFromSlug(),
		middleware.RequireSchemaScope("read"),
		handler.GetContent,
	)

	content.Get("/schema/:schema_slug/:content_slug",
		middleware.GetSchemaFromSlug(),
		middleware.GetContentFromSlug(),
		middleware.RequireSchemaScope("read"),
		handler.GetContentById,
	)

	// Update content - requires {schema}:update scope
	content.Put("/schema/:schema_slug/:content_slug",
		middleware.GetSchemaFromSlug(),
		middleware.GetContentFromSlug(),
		middleware.RequireSchemaScope("update"),
		handler.UpdateContent,
	)

	// Delete content - requires {schema}:delete scope
	content.Delete("/schema/:schema_slug/:content_slug",
		middleware.GetSchemaFromSlug(),
		middleware.GetContentFromSlug(),
		middleware.RequireSchemaScope("delete"),
		handler.DeleteContent,
	)

	// Publish content - requires {schema}:publish scope
	content.Post("/schema/:schema_slug/:content_slug/publish",
		middleware.GetSchemaFromSlug(),
		middleware.GetContentFromSlug(),
		middleware.RequireSchemaScope("publish"),
		handler.PublishContent,
	)
}
