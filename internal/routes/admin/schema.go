package adminroutes

import (
	"contentive/internal/handler"
	"contentive/internal/middleware"
	"contentive/internal/models"

	"github.com/gofiber/fiber/v2"
)

func RegisterAdminSchemaRoutes(app *fiber.App) {
	schema := app.Group("/admin/schema")

	// schema can be accessed only by super admin
	schema.Use(middleware.AuthenticateAdminUserJWT(), middleware.RequireRole(
		models.AdminUserRoleSuperAdmin,
	))

	schema.Post("/", handler.CreateSchema)

	schema.Get("/:id", handler.GetSchema)

	schema.Get("/", handler.ListSchemas)

	schema.Put("/:id", handler.UpdateSchema)

	schema.Delete("/:id", handler.DeleteSchema)
}
