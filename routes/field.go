package routes

import (
	"contentive/controllers"

	"github.com/gofiber/fiber/v2"
)

func RegisterFieldRoutes(app *fiber.App) {
	app.Post("/api/fields", controllers.CreateField)
	app.Get("/api/fields", controllers.GetFields)
	app.Get("/api/content-types/:contentTypeId/fields", controllers.GetFieldsByContentType)
}
