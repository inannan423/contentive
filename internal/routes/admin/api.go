package adminroutes

import (
	"contentive/internal/handler"
	"contentive/internal/middleware"
	"contentive/internal/models"

	"github.com/gofiber/fiber/v2"
)

func RegisterAPIUserRoutes(app *fiber.App) {
	admin := app.Group("/admin")
	api := admin.Group("/api")
	api.Use(middleware.AuthenticateAdminUserJWT())
	api.Use(middleware.RequireRole(
		models.AdminUserRoleSuperAdmin,
	))

	// Create a new API user
	api.Post("/", handler.CreateAPIUser)

	// Get all API users
	api.Get("/", handler.GetAPIUsers)

	// Get an API user by ID
	api.Get("/:id", handler.GetAPIUserByID)

	// Update an API user
	api.Put("/:id", handler.UpdateAPIUser)

	// Delete an API user
	api.Delete("/:id", handler.DeleteAPIUser)

	// Regenerate token
	api.Post("/regenerate-token/:id", handler.RegenerateAPIUserToken)
}
