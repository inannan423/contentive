package routes

import (
	"contentive/internal/handlers"
	"contentive/internal/middlewares"

	"github.com/gofiber/fiber/v2"
)

func RegisterContentEntryRoutes(app *fiber.App) {
	entries := app.Group("/api/content-types/:identifier/entries")

	// Add a new content entry to a content type
	entries.Post("/",
		middlewares.ValidateContentEntry(),
		handlers.CreateContentEntry,
	)

	// Get all content entries for a content type
	entries.Get("/", handlers.GetContentEntries)

	// Get a single content entry for a content type
	entries.Get("/:slug", handlers.GetContentEntry)

	// Update a content entry for a content type
	entries.Put("/:slug",
		middlewares.ValidateContentEntry(),
		handlers.UpdateContentEntry,
	)

	// Delete a content entry for a content type
	entries.Delete("/:slug", handlers.DeleteContentEntry)
}
