package adminroutes

import (
	"contentive/internal/handler"
	"contentive/internal/middleware"
	"contentive/internal/models"

	"github.com/gofiber/fiber/v2"
)

func RegisterAdminContentRoutes(app *fiber.App) {
	content := app.Group("/admin/content")
	content.Use(middleware.AuthenticateAdminUserJWT())
	content.Use(middleware.RequireRole(models.AdminUserRoleEditor))

	// Create content
	content.Post("/schema/:schema_id", handler.CreateContent)

	// Get content
	content.Get("/schema/:schema_id", handler.GetContent)

	// Get content by id
	content.Get("/schema/:schema_id/:content_id", handler.GetContentById)

	// Update content
	content.Put("/schema/:schema_id/:content_id", handler.UpdateContent)

	// Delete content
	content.Delete("/schema/:schema_id/:content_id", handler.DeleteContent)

	// Publish content
	content.Post("/schema/:schema_id/:content_id/publish", handler.PublishContent)
}
