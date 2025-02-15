package main

import (
	"contentive/config"
	"contentive/routes"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	config.InitConfig()

	config.InitDB()

	app := fiber.New()

	routes.RegisterContentTypeRoutes(app)

	log.Fatal(app.Listen(":8080"))
}
