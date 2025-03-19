package main

import (
	"contentive/internal/bootstrap"
	"contentive/internal/config"
	"contentive/internal/database"
	adminroutes "contentive/internal/routes/admin"

	"github.com/gofiber/fiber/v2"
)

func main() {
	config.InitConfig()
	database.InitDB()
	database.InitSchemaValidator()
	bootstrap.InitSuperUser()

	app := fiber.New()

	adminroutes.RegisterAdminUserRoutes(app)
	adminroutes.RegisterAPIUserRoutes(app)
	adminroutes.RegisterAdminSchemaRoutes(app)
	adminroutes.RegisterAdminContentRoutes(app)

	app.Listen(":8080")
}
