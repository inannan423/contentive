package routes

import (
	"contentive/controllers"

	"github.com/gofiber/fiber/v2"
)

func RegisterContentRoutes(app *fiber.App) {
	app.Post("/api/contents", controllers.CreateContent)
	app.Get("/api/contents", controllers.GetContents)
	app.Get("/api/content-types/:contentTypeId/contents", controllers.GetContentsByType)
	app.Put("/api/contents/:id", controllers.UpdateContent)
	app.Delete("/api/contents/:id", controllers.DeleteContent)
	app.Get("/api/contents/:collectionId/items", controllers.GetContentItems)
	app.Post("/api/contents/:collectionId/items", controllers.CreateContentItem)
}
