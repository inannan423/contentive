package routes

import (
	"contentive/internal/handlers"
	"contentive/internal/middlewares"
	"contentive/internal/models"

	"github.com/gofiber/fiber/v2"
)

func RegisterContentTypeRoutes(app *fiber.App) {
	contentTypesAdmin := app.Group("/admin/content-types")

	contentTypesAdmin.Use(middlewares.AuthMiddleware())

	contentTypesAdmin.Post("/",
		middlewares.RequirePermission(models.CreateContentType),
		middlewares.ValidateContentType(),
		handlers.CreateContentType,
	)

	contentTypesAdmin.Get("/",
		middlewares.RequirePermission(models.ReadContentType),
		handlers.GetAllContentTypes,
	)

	contentTypesAdmin.Get("/:identifier",
		middlewares.RequirePermission(models.ReadContentType),
		handlers.GetContentType,
	)

	contentTypesAdmin.Put("/:identifier",
		middlewares.RequirePermission(models.UpdateContentType),
		middlewares.ValidateContentType(),
		handlers.UpdateContentType,
	)

	contentTypesAdmin.Post("/:identifier/fields",
		middlewares.RequirePermission(models.UpdateContentType),
		middlewares.ValidateField(),
		handlers.AddField,
	)

	contentTypesAdmin.Put("/:identifier/fields/:id",
		middlewares.RequirePermission(models.UpdateContentType),
		middlewares.ValidateField(),
		handlers.UpdateField,
	)

	contentTypesAdmin.Delete("/:identifier/fields/:id",
		middlewares.RequirePermission(models.DeleteContentType),
		handlers.DeleteField,
	)
}
