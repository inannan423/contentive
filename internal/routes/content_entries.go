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
	// TODO: Other CRUD operations
}
