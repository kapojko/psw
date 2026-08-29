package llm

import (
	"context"
	"fmt"

	openai "github.com/sashabaranov/go-openai"

	"github.com/kapojko/psw/internal/config"
)

// OpenAICompatibleClient implements the Client interface for an OpenAI-compatible
// provider (e.g. a locally running OpenAI-compatible server).
type OpenAICompatibleClient struct {
	config *config.OpenAICompatibleConfig
	client *openai.Client
}

// NewOpenAICompatibleClient creates a new OpenAI-compatible client
func NewOpenAICompatibleClient(cfg *config.OpenAICompatibleConfig) *OpenAICompatibleClient {
	clientConfig := openai.DefaultConfig(cfg.APIKey)
	clientConfig.BaseURL = cfg.GetBaseURL()

	return &OpenAICompatibleClient{
		config: cfg,
		client: openai.NewClientWithConfig(clientConfig),
	}
}

// ChatCompletion implements Client.ChatCompletion
func (c *OpenAICompatibleClient) ChatCompletion(ctx context.Context, model string, messages []Message) (string, error) {
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
		return "", fmt.Errorf("openai-compatible API error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from openai-compatible provider")
	}

	return resp.Choices[0].Message.Content, nil
}

// ListModels implements Client.ListModels
func (c *OpenAICompatibleClient) ListModels(ctx context.Context) ([]Model, error) {
	models, err := c.client.ListModels(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list openai-compatible models: %w", err)
	}

	result := make([]Model, 0, len(models.Models))
	for _, m := range models.Models {
		result = append(result, Model{
			ID:          m.ID,
			DisplayName: m.ID,
			Provider:    config.ProviderOpenAICompatible,
		})
	}

	return result, nil
}
