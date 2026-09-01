package llm

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	openai "github.com/sashabaranov/go-openai"

	"github.com/kapojko/psw/internal/config"
)

// stepFunBaseURL is the StepFun StepPlan API endpoint used for chat completions.
// Note: StepFun also exposes a generic OpenAI-compatible endpoint at
// https://api.stepfun.ai/v1, but this client targets the StepPlan API.
const stepFunBaseURL = "https://api.stepfun.ai/step_plan/v1"

// StepFunClient implements the Client interface for StepFun
type StepFunClient struct {
	config     *config.StepFunConfig
	client     *openai.Client
	httpClient *http.Client
	verbose    bool
}

// NewStepFunClient creates a new StepFun client with optional proxy
func NewStepFunClient(cfg *config.StepFunConfig, proxyCfg *config.ProxyConfig, verbose bool) *StepFunClient {
	clientConfig := openai.DefaultConfig(cfg.APIKey)
	clientConfig.BaseURL = stepFunBaseURL

	// Create HTTP client with optional proxy
	httpClient := &http.Client{Timeout: 60 * time.Second}
	if proxyCfg != nil && proxyCfg.Enabled && proxyCfg.URL != "" {
		proxyURL, err := url.Parse(proxyCfg.URL)
		if err == nil {
			httpClient.Transport = &http.Transport{
				Proxy: http.ProxyURL(proxyURL),
			}
			if verbose {
				log.Printf("[DEBUG] Proxy configured: %s", proxyCfg.URL)
			}
		} else if verbose {
			log.Printf("[DEBUG] Failed to parse proxy URL: %v", err)
		}
	} else if verbose {
		log.Printf("[DEBUG] Proxy not enabled or not configured")
	}

	return &StepFunClient{
		config:     cfg,
		client:     openai.NewClientWithConfig(clientConfig),
		httpClient: httpClient,
		verbose:    verbose,
	}
}

// ChatCompletion implements Client.ChatCompletion
func (c *StepFunClient) ChatCompletion(ctx context.Context, model string, messages []Message) (string, error) {
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
		return "", fmt.Errorf("stepfun API error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from stepfun")
	}

	return resp.Choices[0].Message.Content, nil
}

// ListModels implements Client.ListModels
func (c *StepFunClient) ListModels(ctx context.Context) ([]Model, error) {
	models, err := c.client.ListModels(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list stepfun models: %w", err)
	}

	result := make([]Model, 0, len(models.Models))
	for _, m := range models.Models {
		result = append(result, Model{
			ID:          m.ID,
			DisplayName: m.ID,
			Provider:    config.ProviderStepFun,
		})
	}

	return result, nil
}
