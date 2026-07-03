package llm

import (
	"context"
	"fmt"

	openai "github.com/sashabaranov/go-openai"

	"github.com/kapojko/psw/internal/config"
)

// LMStudioClient implements the Client interface for LM Studio
type LMStudioClient struct {
	config *config.LMStudioConfig
	client *openai.Client
}

// NewLMStudioClient creates a new LM Studio client
func NewLMStudioClient(cfg *config.LMStudioConfig) *LMStudioClient {
	clientConfig := openai.DefaultConfig("")
	clientConfig.BaseURL = cfg.GetBaseURL()

	return &LMStudioClient{
		config: cfg,
		client: openai.NewClientWithConfig(clientConfig),
	}
}

// ChatCompletion implements Client.ChatCompletion
func (c *LMStudioClient) ChatCompletion(ctx context.Context, model string, messages []Message) (string, error) {
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
		return "", fmt.Errorf("lmstudio API error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from lmstudio")
	}

	return resp.Choices[0].Message.Content, nil
}

// ListModels implements Client.ListModels
func (c *LMStudioClient) ListModels(ctx context.Context) ([]Model, error) {
	models, err := c.client.ListModels(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list lmstudio models: %w", err)
	}

	result := make([]Model, 0, len(models.Models))
	for _, m := range models.Models {
		result = append(result, Model{
			ID:          m.ID,
			DisplayName: m.ID,
			Provider:    config.ProviderLMStudio,
		})
	}

	return result, nil
}
