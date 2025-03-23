package openai

import (
	"bufio"
	"bytes"
	"contentive/internal/llm"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type OpenAIProvider struct {
	APIKey  string
	BaseURL string
	Model   string
	Client  *http.Client
}

type OpenAIChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OpenAIChatRequest represents a request to the OpenAI API
type OpenAIChatRequest struct {
	Model       string              `json:"model"`
	Messages    []OpenAIChatMessage `json:"messages"`
	MaxTokens   int                 `json:"max_tokens,omitempty"`
	Temperature float64             `json:"temperature,omitempty"`
	TopP        float64             `json:"top_p,omitempty"`
	Stream      bool                `json:"stream,omitempty"`
}

// OpenAIChatResponse represents a response from the OpenAI API
type OpenAIChatResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index        int               `json:"index"`
		Message      OpenAIChatMessage `json:"message"`
		FinishReason string            `json:"finish_reason"`
	} `json:"choices"`
}

// OpenAIChatStreamResponse represents a response from the OpenAI API
type OpenAIChatStreamResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index int `json:"index"`
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
}

// NewOpenAIProvider creates a new OpenAIProvider
func NewOpenAIProvider(apiKey, baseURL, model string) *OpenAIProvider {
	if baseURL == "" {
		baseURL = "https://api.openai.com"
	}

	return &OpenAIProvider{
		APIKey:  apiKey,
		BaseURL: baseURL,
		Model:   model,
		Client:  &http.Client{},
	}
}

// Chat sends a chat request to the OpenAI API
func (p *OpenAIProvider) Chat(ctx context.Context, req llm.LLMRequest) (*llm.LLMResponse, error) {
	model := p.Model
	if req.Model != "" {
		model = req.Model
	}

	openAIReq := OpenAIChatRequest{
		Model: model,
		Messages: []OpenAIChatMessage{
			{
				Role:    "user",
				Content: req.Prompt,
			},
		},
		MaxTokens:   req.MaxTokens,
		Temperature: req.Temperature,
		TopP:        req.TopP,
	}

	reqBody, err := json.Marshal(openAIReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	httpReq, err := http.NewRequestWithContext(
		ctx,
		"POST",
		p.BaseURL,
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.APIKey))

	resp, err := p.Client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var openAIResp OpenAIChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&openAIResp); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	if len(openAIResp.Choices) == 0 {
		return nil, fmt.Errorf("no choices in response")
	}

	return &llm.LLMResponse{
		Text:         openAIResp.Choices[0].Message.Content,
		FinishReason: openAIResp.Choices[0].FinishReason,
	}, nil
}

// ChatStream sends a chat request to the OpenAI API
func (p *OpenAIProvider) ChatStream(ctx context.Context, req llm.LLMRequest) (<-chan llm.LLMStreamResponse, error) {
	model := p.Model
	if req.Model != "" {
		model = req.Model
	}

	openAIReq := OpenAIChatRequest{
		Model: model,
		Messages: []OpenAIChatMessage{
			{
				Role:    "user",
				Content: req.Prompt,
			},
		},
		MaxTokens:   req.MaxTokens,
		Temperature: req.Temperature,
		TopP:        req.TopP,
		Stream:      true,
	}

	reqBody, err := json.Marshal(openAIReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	httpReq, err := http.NewRequestWithContext(
		ctx,
		"POST",
		p.BaseURL,
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.APIKey))

	resp, err := p.Client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	responseChan := make(chan llm.LLMStreamResponse)

	go func() {
		defer resp.Body.Close()
		defer close(responseChan)

		reader := bufio.NewReader(resp.Body)

		for {
			select {
			// When the context is done, return
			case <-ctx.Done():
				return
			// Otherwise, read the stream
			default:
				line, err := reader.ReadString('\n')
				if err != nil {
					if err != io.EOF {
						responseChan <- llm.LLMStreamResponse{
							Text: fmt.Sprintf("Error reading stream: %v", err),
							Done: true,
						}
					}
					return
				}

				line = strings.TrimSpace(line)
				if line == "" {
					continue
				}

				if !strings.HasPrefix(line, "data: ") {
					continue
				}

				data := strings.TrimPrefix(line, "data: ")
				if data == "[DONE]" {
					responseChan <- llm.LLMStreamResponse{Done: true}
					return
				}

				var streamResp OpenAIChatStreamResponse
				if err := json.Unmarshal([]byte(data), &streamResp); err != nil {
					continue
				}

				if len(streamResp.Choices) == 0 {
					continue
				}

				choice := streamResp.Choices[0]
				responseChan <- llm.LLMStreamResponse{
					Text:         choice.Delta.Content,
					FinishReason: choice.FinishReason,
					Done:         choice.FinishReason != "",
				}

				if choice.FinishReason != "" {
					return
				}
			}
		}
	}()

	return responseChan, nil
}
