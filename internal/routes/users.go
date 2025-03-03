package routes

import (
	"contentive/internal/handlers"
	"contentive/internal/middlewares"
	"contentive/internal/models"

	"github.com/gofiber/fiber/v2"
)

func RegisterUserRoutes(app *fiber.App) {
	auth := app.Group("/admin/auth")
	auth.Post("/login", handlers.Login)

	auth.Get("/validate",
		middlewares.AuthMiddleware(),
		handlers.ValidateToken,
	)

	users := app.Group("/admin/users")
	users.Use(middlewares.AuthMiddleware())

	users.Post("/",
		middlewares.RequirePermission(models.ManageUsers),
		handlers.CreateUser,
	)

	users.Get("/",
		middlewares.RequirePermission(models.ManageUsers),
		handlers.GetUsers,
	)

	users.Put("/:id",
		middlewares.RequirePermission(models.ManageUsers),
		middlewares.ValidateSuperAdminOperation(),
		handlers.UpdateUser,
	)

	users.Delete("/:id",
		middlewares.RequirePermission(models.ManageUsers),
		middlewares.ValidateSuperAdminOperation(),
		handlers.DeleteUser,
	)
}
