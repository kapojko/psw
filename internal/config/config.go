package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config represents the application configuration
type Config struct {
	Proxy        *ProxyConfig    `json:"proxy,omitempty"`
	Providers    ProvidersConfig `json:"providers"`
	DefaultModel *ModelRef       `json:"default_model,omitempty"`
}

// ProxyConfig holds proxy configuration
type ProxyConfig struct {
	Enabled bool   `json:"enabled"`
	URL     string `json:"url"` // e.g. http://proxy:8080
}

// Load reads the configuration from disk
func Load() (*Config, error) {
	path, err := GetConfigPath()
	if err != nil {
		return nil, fmt.Errorf("failed to get config path: %w", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{
				Providers: ProvidersConfig{
					OpenRouter: &OpenRouterConfig{},
					LMStudio:   &LMStudioConfig{BaseURL: "http://localhost:1234/v1"},
				},
			}, nil
		}
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Ensure provider configs exist
	if cfg.Providers.OpenRouter == nil {
		cfg.Providers.OpenRouter = &OpenRouterConfig{}
	}
	if cfg.Providers.LMStudio == nil {
		cfg.Providers.LMStudio = &LMStudioConfig{BaseURL: "http://localhost:1234/v1"}
	}

	return &cfg, nil
}

// Save writes the configuration to disk
func (c *Config) Save() error {
	if err := EnsureConfigDir(); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	path, err := GetConfigPath()
	if err != nil {
		return fmt.Errorf("failed to get config path: %w", err)
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, FilePerm); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}
