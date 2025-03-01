package routes

import (
	"contentive/internal/handlers"
	"contentive/internal/middlewares"
	"contentive/internal/models"

	"github.com/gofiber/fiber/v2"
)

func RegisterAPIRoleRoutes(app *fiber.App) {
	// Routes for managing API roles (requires admin permission)
	apiRoles := app.Group("/admin/api-roles")
	apiRoles.Use(middlewares.AuthMiddleware())
	apiRoles.Use(middlewares.RequirePermission(models.ManageRoles))

	// API role management
	apiRoles.Get("/", handlers.GetAPIRoles)
	apiRoles.Get("/:id", handlers.GetAPIRole)
	apiRoles.Post("/", handlers.CreateAPIRole)
	apiRoles.Put("/:id", handlers.UpdateAPIRole)
	apiRoles.Delete("/:id", handlers.DeleteAPIRole)
	apiRoles.Post("/:id/regenerate-key", handlers.RegenerateAPIKey)

	// API permission management
	apiRoles.Get("/:id/permissions", handlers.GetAPIRolePermissions)
	apiRoles.Put("/:id/permissions", handlers.BatchUpdateAPIRolePermissions)
	apiRoles.Put("/:id/permission", handlers.UpdateAPIRolePermission)
}
