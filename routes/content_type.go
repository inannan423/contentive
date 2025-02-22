package routes

import (
	"contentive/controllers"

	"github.com/gofiber/fiber/v2"
)

func RegisterContentTypeRoutes(app *fiber.App) {
	app.Post("/api/content-types", controllers.CreateContentType)
	app.Get("/api/content-types", controllers.GetContentTypes)
	app.Get("/api/content-types/:slug", controllers.GetContentTypeBySlug)
}
