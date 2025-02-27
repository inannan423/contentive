package routes

import (
	"contentive/internal/handlers"
	"contentive/internal/middlewares"
	"contentive/internal/models"

	"github.com/gofiber/fiber/v2"
)

func RegisterUserRoutes(app *fiber.App) {
	auth := app.Group("/api/auth")
	auth.Post("/login", handlers.Login)

	users := app.Group("/api/users")
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
		handlers.UpdateUser,
	)
}
