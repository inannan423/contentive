package main

import (
	"contentive/config"
	"contentive/internal/bootstrap"
	"contentive/internal/routes"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	config.InitConfig()
	config.InitDB()

	bootstrap.InitRolesAndPermissions()
	bootstrap.InitSuperAdmin()
	bootstrap.InitAPIRoles()
	bootstrap.InitDefaultAPIPermissions()

	app := fiber.New()

	// Add CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowCredentials: true,
	}))

	routes.RegisterUserRoutes(app)
	routes.RegisterContentTypeRoutes(app)
	routes.RegisterContentEntryRoutes(app)
	routes.RegisterRoleRoutes(app)
	routes.RegisterAPIRoleRoutes(app)

	log.Fatal(app.Listen(":8080"))
}
