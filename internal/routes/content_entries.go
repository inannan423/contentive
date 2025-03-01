package routes

import (
	"contentive/internal/handlers"
	"contentive/internal/middlewares"
	"contentive/internal/models"

	"github.com/gofiber/fiber/v2"
)

func RegisterContentEntryRoutes(app *fiber.App) {
	entries := app.Group("/api/content-types/:identifier/entries")

	entries.Use(middlewares.AuthMiddleware())

	// Add a new content entry to a content type
	entries.Post("/",
		middlewares.RequirePermission(models.CreateContent),
		middlewares.ValidateContentEntry(),
		handlers.CreateContentEntry,
	)

	// Get all content entries for a content type
	entries.Get("/",
		middlewares.RequirePermission(models.ReadContent),
		handlers.GetContentEntries,
	)

	// Get a single content entry for a content type
	entries.Get("/:slug",
		middlewares.RequirePermission(models.ReadContent),
		handlers.GetContentEntry,
	)

	// Update a content entry for a content type
	entries.Put("/:slug",
		middlewares.RequirePermission(models.UpdateContent),
		middlewares.ValidateContentEntry(),
		handlers.UpdateContentEntry,
	)

	// Delete a content entry for a content type
	entries.Delete("/:slug",
		middlewares.RequirePermission(models.DeleteContent),
		handlers.DeleteContentEntry,
	)

	// API routes
	apiEntries := app.Group("/api/content-types/:identifier/entries")
	apiEntries.Use(middlewares.APIAuthMiddleware())

	apiEntries.Post("/",
		middlewares.APIPermissionMiddleware(models.CreateOperation),
		middlewares.ValidateContentEntry(),
		handlers.CreateContentEntry,
	)

	apiEntries.Get("/",
		middlewares.APIPermissionMiddleware(models.ReadOperation),
		handlers.GetContentEntries,
	)

	apiEntries.Get("/:slug",
		middlewares.APIPermissionMiddleware(models.ReadOperation),
		handlers.GetContentEntry,
	)

	apiEntries.Put("/:slug",
		middlewares.APIPermissionMiddleware(models.UpdateOperation),
		middlewares.ValidateContentEntry(),
		handlers.UpdateContentEntry,
	)

	apiEntries.Delete("/:slug",
		middlewares.APIPermissionMiddleware(models.DeleteOperation),
		handlers.DeleteContentEntry,
	)
}
