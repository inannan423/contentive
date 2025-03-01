package main

import (
	"contentive/config"
	"contentive/internal/bootstrap"
	"contentive/internal/routes"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	config.InitConfig()
	config.InitDB()

	bootstrap.InitRolesAndPermissions()
	bootstrap.InitSuperAdmin()
	bootstrap.InitAPIRoles()
	bootstrap.InitDefaultAPIPermissions()

	app := fiber.New()

	routes.RegisterUserRoutes(app)
	routes.RegisterContentTypeRoutes(app)
	routes.RegisterContentEntryRoutes(app)
	routes.RegisterRoleRoutes(app)
	routes.RegisterAPIRoleRoutes(app)

	log.Fatal(app.Listen(":8080"))
}
