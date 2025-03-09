package main

import (
	"contentive/internal/bootstrap"
	"contentive/internal/config"
	adminroutes "contentive/internal/routes/admin"

	"github.com/gofiber/fiber/v2"
)

func main() {
	config.InitConfig()
	config.InitDB()
	bootstrap.InitSuperUser()

	app := fiber.New()

	adminroutes.RegisterAdminUserRoutes(app)

	app.Listen(":8080")
}
