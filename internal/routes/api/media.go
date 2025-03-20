package apiroutes

import (
	"contentive/internal/handler"
	"contentive/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterAPIMediaRoutes(app *fiber.App) {
	api := app.Group("/api")
	media := api.Group("/media")

	// All API routes require API token authentication
	media.Use(middleware.AuthenticateAPIUserToken())

	// Routes for media operations
	// Upload media - requires media:create scope
	media.Post("/", middleware.RequireAPIScope("media:create"), handler.UploadMedia)

	// Get media - requires media:read scope
	media.Get("/:id", middleware.RequireAPIScope("media:read"), handler.GetMedia)
	media.Get("/", middleware.RequireAPIScope("media:read"), handler.ListMedia)

	// Delete media - requires media:delete scope
	media.Delete("/:id", middleware.RequireAPIScope("media:delete"), handler.DeleteMedia)
}
