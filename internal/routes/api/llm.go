package apiroutes

import (
	"contentive/internal/handler"
	"contentive/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterAPILLMRoutes(app *fiber.App) {
	api := app.Group("/api")
	llmRoutes := api.Group("/llm")
	// All API routes require API token authentication
	llmRoutes.Use(middleware.AuthenticateAPIUserToken())
	// Chat
	llmRoutes.Post("/chat", middleware.RequireAPIScope("llm:chat"), handler.LLMChat)

	// Stream Chat
	llmRoutes.Post("/chat/stream", middleware.RequireAPIScope("llm:chat"), handler.LLMChatStream)

	// RAG
	llmRoutes.Post("/knowledge", middleware.RequireAPIScope("llm:rag"), handler.LLMKnowledgeQuery)

	// Stream RAG
	llmRoutes.Post("/knowledge/stream", middleware.RequireAPIScope("llm:rag"), handler.LLMKnowledgeQueryStream)
}
