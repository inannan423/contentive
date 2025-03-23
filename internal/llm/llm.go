package llm

import "context"

type LLMRequest struct {
	Prompt      string                 `json:"prompt"`
	MaxTokens   int                    `json:"max_tokens,omitempty"`
	Temperature float64                `json:"temperature,omitempty"`
	TopP        float64                `json:"top_p,omitempty"`
	Model       string                 `json:"model,omitempty"`
	Stream      bool                   `json:"stream,omitempty"`
	ExtraParams map[string]interface{} `json:"extra_params,omitempty"`
}

type LLMResponse struct {
	Text         string `json:"text"`
	FinishReason string `json:"finish_reason"` // the reson why the generation finished
}

type LLMStreamResponse struct {
	Text         string `json:"text"`
	FinishReason string `json:"finish_reason,omitempty"`
	Done         bool   `json:"done"`
}

type LLMProvider interface {
	// Get full response
	Chat(ctx context.Context, req LLMRequest) (*LLMResponse, error)
	// Get stream response, the channel will be closed when the generation is finished
	ChatStream(ctx context.Context, req LLMRequest) (<-chan LLMStreamResponse, error)
}

var provider LLMProvider

// SetProvider sets the LLM provider
func SetProvider(p LLMProvider) {
	provider = p
}

// Get LLM provider
func GetProvider() LLMProvider {
	return provider
}
