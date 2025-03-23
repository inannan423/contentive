package main

import (
	"contentive/internal/bootstrap"
	"contentive/internal/config"
	"contentive/internal/database"
	llm "contentive/internal/llm"
	"contentive/internal/llm/openai"
	adminroutes "contentive/internal/routes/admin"
	apiroutes "contentive/internal/routes/api"
	"contentive/internal/storage"
	"contentive/internal/storage/aliyun"
	"contentive/internal/storage/local"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	config.InitConfig()
	database.InitDB()
	database.InitSchemaValidator()
	bootstrap.InitSuperUser()

	// init LLM	Provider
	initLLMProvider()

	// init storage
	initStorageProvider()

	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET,POST,HEAD,PUT,DELETE,PATCH",
	}))

	adminroutes.RegisterAdminUserRoutes(app)
	adminroutes.RegisterAPIUserRoutes(app)
	adminroutes.RegisterAdminSchemaRoutes(app)
	adminroutes.RegisterAdminContentRoutes(app)
	adminroutes.RegisterAdminMediaRoutes(app)

	apiroutes.RegisterAPIContentRoutes(app)
	apiroutes.RegisterAPIMediaRoutes(app)

	app.Listen(":8080")
}

// initStorageProvider initializes the appropriate storage provider based on configuration
func initStorageProvider() {
	switch config.AppConfig.MEDIA_STORAGE_TYPE {
	case "local":
		storage.SetStorageProvider(local.NewLocalStorage(
			config.AppConfig.MEDIA_STORAGE_PATH,
			config.AppConfig.MEDIA_STORAGE_URL,
		))
		log.Println("Local storage provider initialized")
	case "aliyun":
		ossProvider, err := aliyun.NewAliyunOSSStorage(
			config.AppConfig.OSS_REGION_ID,
			config.AppConfig.OSS_ACCESS_KEY_ID,
			config.AppConfig.OSS_ACCESS_KEY_SECRET,
			config.AppConfig.OSS_BUCKET_NAME,
			config.AppConfig.MEDIA_STORAGE_URL,
		)
		if err != nil {
			log.Fatalf("Failed to initialize Aliyun OSS storage provider: %v", err)
		}
		storage.SetStorageProvider(ossProvider)
		log.Println("Aliyun OSS storage provider initialized")
	default:
		log.Fatalf("Unsupported storage type: %s", config.AppConfig.MEDIA_STORAGE_TYPE)
	}
}

// initLLMProvider initializes the appropriate LLM provider based on configuration
func initLLMProvider() {
	switch config.AppConfig.LLM_PROVIDER {
	case "openai":
		llm.SetProvider(openai.NewOpenAIProvider(
			config.AppConfig.LLM_API_KEY,
			config.AppConfig.LLM_BASE_URL,
			config.AppConfig.LLM_MODEL,
		))
		log.Println("OpenAI LLM provider initialized")
	case "qwen":
		// qwen support openai client
		llm.SetProvider(openai.NewOpenAIProvider(
			config.AppConfig.LLM_API_KEY,
			config.AppConfig.LLM_BASE_URL,
			config.AppConfig.LLM_MODEL,
		))
		log.Println("Qwen LLM provider initialized")
	default:
		log.Fatalf("Unsupported LLM provider: %s", config.AppConfig.LLM_PROVIDER)
	}
}
