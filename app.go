package main

import (
	"contentive/internal/bootstrap"
	"contentive/internal/config"
)

func main() {
	config.InitConfig()
	config.InitDB()
	bootstrap.InitSuperUser()

	// app := fiber.New()
}
