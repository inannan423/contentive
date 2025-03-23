package models

type LLMProvider string

const (
	LLMProviderOpenAI LLMProvider = "openai"
	LLMProviderQwen   LLMProvider = "qwen"
	LLMProviderDoubao LLMProvider = "doubao"
)
