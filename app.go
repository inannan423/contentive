package main

import (
	"contentive/internal/bootstrap"
	"contentive/internal/config"
	"contentive/internal/database"
	adminroutes "contentive/internal/routes/admin"
	"contentive/internal/storage"
	"contentive/internal/storage/local"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	config.InitConfig()
	database.InitDB()
	database.InitSchemaValidator()
	bootstrap.InitSuperUser()

	// init storage
	if config.AppConfig.MEDIA_STORAGE_TYPE == "local" {
		storage.SetStorageProvider(local.NewLocalStorage(
			config.AppConfig.MEDIA_STORAGE_PATH,
			config.AppConfig.MEDIA_STORAGE_PATH,
		))
	} else {
		log.Fatal("Unsupported storage type")
	}

	app := fiber.New()

	adminroutes.RegisterAdminUserRoutes(app)
	adminroutes.RegisterAPIUserRoutes(app)
	adminroutes.RegisterAdminSchemaRoutes(app)
	adminroutes.RegisterAdminContentRoutes(app)
	adminroutes.RegisterAdminMediaRoutes(app)

	app.Listen(":8080")
}
