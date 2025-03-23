package handler

import (
	"bufio"
	"contentive/internal/llm"
	"contentive/internal/logger"
	"context"
	"encoding/json"
	"time"

	"github.com/gofiber/fiber/v2"
)

// LLMChatRequest is a struct that contains the request for the LLM chat
type LLMChatRequest struct {
	Prompt      string  `json:"prompt"`
	MaxTokens   int     `json:"max_tokens,omitempty"`
	Temperature float64 `json:"temperature,omitempty"`
	TopP        float64 `json:"top_p,omitempty"`
	Model       string  `json:"model,omitempty"`
}

// LLMChat is a handler that handles the LLM chat request
func LLMChat(c *fiber.Ctx) error {
	var req LLMChatRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Error("Failed to parse LLM chat request: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Prompt == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Prompt is required",
		})
	}

	llmReq := llm.LLMRequest{
		Prompt:      req.Prompt,
		MaxTokens:   req.MaxTokens,
		Temperature: req.Temperature,
		TopP:        req.TopP,
		Model:       req.Model,
	}

	// Get the LLM provider
	provider := llm.GetProvider()
	if provider == nil {
		logger.Error("LLM provider not initialized")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "LLM service not available",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := provider.Chat(ctx, llmReq)
	if err != nil {
		logger.Error("LLM chat error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get response from LLM",
		})
	}

	return c.JSON(resp)
}

// LLMChatStream is a handler that handles the LLM chat stream request
func LLMChatStream(c *fiber.Ctx) error {
	var req LLMChatRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Error("Failed to parse LLM chat stream request: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Prompt == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Prompt is required",
		})
	}

	llmReq := llm.LLMRequest{
		Prompt:      req.Prompt,
		MaxTokens:   req.MaxTokens,
		Temperature: req.Temperature,
		TopP:        req.TopP,
		Model:       req.Model,
		Stream:      true,
	}

	provider := llm.GetProvider()
	if provider == nil {
		logger.Error("LLM provider not initialized")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "LLM service not available",
		})
	}

	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")

	done := make(chan bool)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	stream, err := provider.ChatStream(ctx, llmReq)
	if err != nil {
		logger.Error("LLM chat stream error: %v", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get stream from LLM",
		})
	}

	c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		defer close(done)

		for {
			select {
			case <-ctx.Done():
				data, _ := json.Marshal(fiber.Map{"done": true, "error": "Request timeout"})
				w.Write([]byte("data: " + string(data) + "\n\n"))
				w.Flush()
				return
			case resp, ok := <-stream:
				if !ok {
					data, _ := json.Marshal(fiber.Map{"done": true})
					w.Write([]byte("data: " + string(data) + "\n\n"))
					w.Flush()
					return
				}

				data, _ := json.Marshal(resp)
				w.Write([]byte("data: " + string(data) + "\n\n"))
				w.Flush()

				if resp.Done {
					return
				}
			}
		}
	})

	<-done
	return nil
}

// LLMKnowledgeQuery is a handler that handles the LLM knowledge query request
func LLMKnowledgeQuery(c *fiber.Ctx) error {
	// TODO: RAG query endpoint
	return c.JSON(fiber.Map{
		"message": "Knowledge query endpoint - To be implemented",
	})
}

// LLMKnowledgeQueryStream is a handler that handles the LLM knowledge query stream request
func LLMKnowledgeQueryStream(c *fiber.Ctx) error {
	// TODO: RAG query stream endpoint
	return c.JSON(fiber.Map{
		"message": "Knowledge query stream endpoint - To be implemented",
	})
}
