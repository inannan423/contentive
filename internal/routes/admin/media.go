package adminroutes

import (
	"contentive/internal/handler"
	"contentive/internal/middleware"
	"contentive/internal/models"

	"github.com/gofiber/fiber/v2"
)

func RegisterAdminMediaRoutes(app *fiber.App) {
	media := app.Group("/admin/media")
	media.Use(middleware.AuthenticateAdminUserJWT())
	media.Use(middleware.RequireRole(models.AdminUserRoleEditor))

	media.Post("/", handler.UploadMedia)
	media.Get("/:id", handler.GetMedia)
	media.Delete("/:id", handler.DeleteMedia)
	media.Get("/", handler.ListMedia)
}
