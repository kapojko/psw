package config

import "fmt"

// ProviderType represents a supported LLM provider
type ProviderType string

const (
	ProviderOpenRouter ProviderType = "openrouter"
	ProviderLMStudio   ProviderType = "lmstudio"
)

// ProviderConfig is the interface all provider configs must implement
type ProviderConfig interface {
	GetType() ProviderType
	IsEnabled() bool
	GetDisplayName() string
}

// OpenRouterConfig holds OpenRouter-specific configuration
type OpenRouterConfig struct {
	APIKey string `json:"api_key,omitempty"`
}

func (c *OpenRouterConfig) GetType() ProviderType {
	return ProviderOpenRouter
}

func (c *OpenRouterConfig) IsEnabled() bool {
	return c.APIKey != ""
}

func (c *OpenRouterConfig) GetDisplayName() string {
	return "OpenRouter"
}

// LMStudioConfig holds LM Studio-specific configuration
type LMStudioConfig struct {
	Enabled bool   `json:"enabled"`
	BaseURL string `json:"base_url"`
}

func (c *LMStudioConfig) GetType() ProviderType {
	return ProviderLMStudio
}

func (c *LMStudioConfig) IsEnabled() bool {
	return c.Enabled
}

func (c *LMStudioConfig) GetDisplayName() string {
	return "LM Studio"
}

// GetBaseURL returns the base URL, defaulting to localhost if not set
func (c *LMStudioConfig) GetBaseURL() string {
	if c.BaseURL == "" {
		return "http://localhost:1234/v1"
	}
	return c.BaseURL
}

// ModelRef represents a reference to a specific model on a specific provider
type ModelRef struct {
	Provider ProviderType `json:"provider"`
	ModelID  string       `json:"model_id"`
}

// String returns a human-readable representation of the model reference
func (m ModelRef) String() string {
	return fmt.Sprintf("%s/%s", m.Provider, m.ModelID)
}

// ParseModelRef parses a "provider/model" string into a ModelRef
func ParseModelRef(s string) (ModelRef, error) {
	for i, c := range s {
		if c == '/' {
			provider := ProviderType(s[:i])
			if provider != ProviderOpenRouter && provider != ProviderLMStudio {
				return ModelRef{}, fmt.Errorf("unknown provider: %s", provider)
			}
			return ModelRef{
				Provider: provider,
				ModelID:  s[i+1:],
			}, nil
		}
	}
	return ModelRef{}, fmt.Errorf("invalid model format (expected provider/model): %s", s)
}

// ProvidersConfig holds configuration for all providers
type ProvidersConfig struct {
	OpenRouter *OpenRouterConfig `json:"openrouter,omitempty"`
	LMStudio   *LMStudioConfig   `json:"lmstudio,omitempty"`
}

// GetProviders returns all provider configs
func (p *ProvidersConfig) GetProviders() []ProviderConfig {
	var providers []ProviderConfig
	if p.OpenRouter != nil {
		providers = append(providers, p.OpenRouter)
	}
	if p.LMStudio != nil {
		providers = append(providers, p.LMStudio)
	}
	return providers
}

// GetProvider returns a specific provider config by type
func (p *ProvidersConfig) GetProvider(providerType ProviderType) ProviderConfig {
	switch providerType {
	case ProviderOpenRouter:
		return p.OpenRouter
	case ProviderLMStudio:
		return p.LMStudio
	default:
		return nil
	}
}
