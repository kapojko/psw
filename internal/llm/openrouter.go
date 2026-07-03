package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	openai "github.com/sashabaranov/go-openai"

	"github.com/kapojko/psw/internal/config"
)

const openRouterBaseURL = "https://openrouter.ai/api/v1"

// OpenRouterClient implements the Client interface for OpenRouter
type OpenRouterClient struct {
	config     *config.OpenRouterConfig
	client     *openai.Client
	httpClient *http.Client
}

// NewOpenRouterClient creates a new OpenRouter client with optional proxy
func NewOpenRouterClient(cfg *config.OpenRouterConfig, proxyCfg *config.ProxyConfig) *OpenRouterClient {
	clientConfig := openai.DefaultConfig(cfg.APIKey)
	clientConfig.BaseURL = openRouterBaseURL

	// Create HTTP client with optional proxy
	httpClient := &http.Client{Timeout: 60 * time.Second}
	if proxyCfg != nil && proxyCfg.Enabled && proxyCfg.URL != "" {
		proxyURL, err := url.Parse(proxyCfg.URL)
		if err == nil {
			httpClient.Transport = &http.Transport{
				Proxy: http.ProxyURL(proxyURL),
			}
		}
	}

	return &OpenRouterClient{
		config:     cfg,
		client:     openai.NewClientWithConfig(clientConfig),
		httpClient: httpClient,
	}
}

// ChatCompletion implements Client.ChatCompletion
func (c *OpenRouterClient) ChatCompletion(ctx context.Context, model string, messages []Message) (string, error) {
	msgs := make([]openai.ChatCompletionMessage, len(messages))
	for i, m := range messages {
		msgs[i] = openai.ChatCompletionMessage{
			Role:    m.Role,
			Content: m.Content,
		}
	}

	resp, err := c.client.CreateChatCompletion(ctx, openai.ChatCompletionRequest{
		Model:    model,
		Messages: msgs,
	})
	if err != nil {
		return "", fmt.Errorf("openrouter API error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from openrouter")
	}

	return resp.Choices[0].Message.Content, nil
}

type openRouterModelsResponse struct {
	Data []openRouterModel `json:"data"`
}

type openRouterModel struct {
	ID string `json:"id"`
}

// ListModels implements Client.ListModels
func (c *OpenRouterClient) ListModels(ctx context.Context) ([]Model, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://openrouter.ai/api/v1/models", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.config.APIKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to list openrouter models: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("openrouter API returned status %d: %s", resp.StatusCode, string(body))
	}

	var modelsResp openRouterModelsResponse
	if err := json.Unmarshal(body, &modelsResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	result := make([]Model, 0, len(modelsResp.Data))
	for _, m := range modelsResp.Data {
		result = append(result, Model{
			ID:          m.ID,
			DisplayName: m.ID,
			Provider:    config.ProviderOpenRouter,
		})
	}

	return result, nil
}
