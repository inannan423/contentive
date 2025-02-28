package routes

import (
	"contentive/internal/handlers"
	"contentive/internal/middlewares"
	"contentive/internal/models"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoleRoutes(app *fiber.App) {
	roles := app.Group("/admin/roles")
	roles.Use(middlewares.AuthMiddleware())

	roles.Get("/",
		middlewares.RequirePermission(models.ManageRoles),
		handlers.GetRoles,
	)

	roles.Get("/:id",
		middlewares.RequirePermission(models.ManageRoles),
		handlers.GetRole,
	)

	permissions := app.Group("/admin/permissions")
	permissions.Use(middlewares.AuthMiddleware())
	permissions.Get("/",
		middlewares.RequirePermission(models.ManageRoles),
		handlers.GetPermissions,
	)
}
