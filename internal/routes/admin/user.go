package adminroutes

import (
	"contentive/internal/handler"
	"contentive/internal/middleware"
	"contentive/internal/models"

	"github.com/gofiber/fiber/v2"
)

func RegisterAdminUserRoutes(app *fiber.App) {
	admin := app.Group("/admin")

	auth := admin.Group("/auth")
	// Login
	auth.Post("/login", handler.AdminUserLogin)

	admin.Use(middleware.AuthenticateAdminUserJWT())

	users := admin.Group("/users")
	// Get all users
	users.Get("/", middleware.RequireRole(
		models.AdminUserRoleViewer,
	), handler.GetAllAdminUsers)

	user := admin.Group("/user")
	// Get user by id
	user.Get("/:id", middleware.RequireRole(
		models.AdminUserRoleViewer,
	), handler.GetAdminUserById)

	// Create user
	user.Post("/", middleware.RequireRole(
		models.AdminUserRoleSuperAdmin,
	), handler.CreateAdminUser)

	// Update user
	user.Put("/:id", middleware.RequireRole(
		models.AdminUserRoleSuperAdmin,
	), handler.UpdateAdminUser)

	// Delete user
	user.Delete("/:id", middleware.RequireRole(
		models.AdminUserRoleSuperAdmin,
	), handler.DeleteAdminUser)
}
