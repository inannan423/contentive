package routes

import (
	"contentive/internal/handlers"
	"contentive/internal/middlewares"

	"github.com/gofiber/fiber/v2"
)

func RegisterContentEntryRoutes(app *fiber.App) {
	entries := app.Group("/api/content-types/:contentTypeId/entries")

	entries.Post("/",
		middlewares.ValidateContentEntry(),
		handlers.CreateContentEntry,
	)

	entries.Get("/", handlers.GetContentEntries)

	entries.Get("/:id", handlers.GetContentEntry)

	entries.Put("/:id",
		middlewares.ValidateContentEntry(),
		handlers.UpdateContentEntry,
	)

	entries.Delete("/:id", handlers.DeleteContentEntry)
}
