package llm

import (
	"context"
	"fmt"

	"github.com/kapojko/psw/internal/config"
)

// Message represents a chat message
type Message struct {
	Role    string // "system", "user", "assistant"
	Content string
}

// Model represents an available model
type Model struct {
	ID          string
	DisplayName string
	Provider    config.ProviderType
}

// String returns a human-readable model identifier
func (m Model) String() string {
	return fmt.Sprintf("%s/%s", m.Provider, m.ID)
}

// Client is the interface all LLM providers must implement
type Client interface {
	// ChatCompletion sends messages and returns the completion
	ChatCompletion(ctx context.Context, model string, messages []Message) (string, error)

	// ListModels returns available models from this provider
	ListModels(ctx context.Context) ([]Model, error)
}

// NewClient creates an LLM client for the given provider config
func NewClient(providerConfig config.ProviderConfig, proxyConfig *config.ProxyConfig) (Client, error) {
	switch cfg := providerConfig.(type) {
	case *config.OpenRouterConfig:
		return NewOpenRouterClient(cfg, proxyConfig), nil
	case *config.LMStudioConfig:
		return NewLMStudioClient(cfg), nil
	default:
		return nil, fmt.Errorf("unsupported provider type: %T", providerConfig)
	}
}
