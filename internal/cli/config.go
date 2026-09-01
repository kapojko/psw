package cli

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/kapojko/psw/internal/config"
	"github.com/kapojko/psw/internal/llm"
)

var reader = bufio.NewReader(os.Stdin)

// RunConfig starts the interactive configuration wizard
func RunConfig() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	fmt.Println("=== psw Setup ===")
	fmt.Println()

	// Step 0: Proxy config
	if err := configureProxy(cfg); err != nil {
		return err
	}

	// Step 1: OpenRouter config
	if err := configureOpenRouter(cfg); err != nil {
		return err
	}

	// Step 2: OpenAI Compatible config
	if err := configureOpenAICompatible(cfg); err != nil {
		return err
	}

	// Step 3: StepFun config
	if err := configureStepFun(cfg); err != nil {
		return err
	}

	// Step 4: Default model selection
	if err := configureDefaultModel(cfg); err != nil {
		return err
	}

	// Save config
	if err := cfg.Save(); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Println()
	fmt.Println("Configuration saved successfully!")
	return nil
}

func configureProxy(cfg *config.Config) error {
	fmt.Println("0. Proxy Setup (for OpenRouter API)")

	if cfg.Proxy != nil && cfg.Proxy.Enabled {
		fmt.Printf("   Current proxy: %s\n", cfg.Proxy.URL)
	} else {
		fmt.Println("   Current proxy: Not configured")
	}

	fmt.Println("   Options:")
	fmt.Println("   [1] Configure HTTP proxy")
	fmt.Println("   [2] Don't use proxy")

	choice, err := readChoice(1, 2)
	if err != nil {
		return err
	}

	switch choice {
	case 1:
		fmt.Print("   Enter proxy URL (e.g. http://proxy:8080): ")
		url, err := readLine()
		if err != nil {
			return err
		}
		url = strings.TrimSpace(url)
		if url == "" {
			// Keep existing if empty
			if cfg.Proxy == nil || !cfg.Proxy.Enabled {
				return fmt.Errorf("proxy URL cannot be empty")
			}
		} else {
			cfg.Proxy = &config.ProxyConfig{
				Enabled: true,
				URL:     url,
			}
		}
	case 2:
		cfg.Proxy = nil
	}

	fmt.Println()
	return nil
}

func configureOpenRouter(cfg *config.Config) error {
	fmt.Println("1. OpenRouter Setup")

	if cfg.Providers.OpenRouter.APIKey != "" {
		masked := maskAPIKey(cfg.Providers.OpenRouter.APIKey)
		fmt.Printf("   Current API Key: %s\n", masked)
	} else {
		fmt.Println("   Current API Key: Not configured")
	}

	fmt.Println("   Options:")
	fmt.Println("   [1] Enter/Update API Key")
	fmt.Println("   [2] Remove API Key (disable provider)")
	fmt.Println("   [3] Keep current (press ENTER)")

	choice, err := readChoiceOrKeep(1, 3, 3)
	if err != nil {
		return err
	}

	switch choice {
	case 1:
		fmt.Print("   Enter API Key: ")
		key, err := readLine()
		if err != nil {
			return err
		}
		key = strings.TrimSpace(key)
		if key == "" {
			// Keep existing if empty
			if cfg.Providers.OpenRouter.APIKey == "" {
				return fmt.Errorf("API key cannot be empty")
			}
		} else {
			cfg.Providers.OpenRouter.APIKey = key
		}
	case 2:
		cfg.Providers.OpenRouter.APIKey = ""
	case 3:
		// Keep current
	}

	fmt.Println()
	return nil
}

func configureOpenAICompatible(cfg *config.Config) error {
	fmt.Println("2. OpenAI Compatible Setup")

	if cfg.Providers.OpenAICompatible.BaseURL != "" {
		fmt.Printf("   Current endpoint: %s\n", cfg.Providers.OpenAICompatible.GetBaseURL())
	} else {
		fmt.Println("   Current endpoint: Not configured")
	}

	if cfg.Providers.OpenAICompatible.APIKey != "" {
		fmt.Printf("   Current API Key: %s\n", maskAPIKey(cfg.Providers.OpenAICompatible.APIKey))
	} else {
		fmt.Println("   Current API Key: Not configured")
	}

	fmt.Println("   Options:")
	fmt.Println("   [1] Configure endpoint")
	fmt.Println("   [2] Set API Key (optional)")
	fmt.Println("   [3] Remove API Key")
	fmt.Println("   [4] Keep current (press ENTER)")

	choice, err := readChoiceOrKeep(1, 4, 4)
	if err != nil {
		return err
	}

	switch choice {
	case 1:
		fmt.Print("   Enter endpoint (e.g. http://127.0.0.1:1234/v1): ")
		endpoint, err := readLine()
		if err != nil {
			return err
		}
		endpoint = strings.TrimSpace(endpoint)
		if endpoint == "" {
			return fmt.Errorf("endpoint cannot be empty")
		}
		cfg.Providers.OpenAICompatible.BaseURL = endpoint
	case 2:
		fmt.Print("   Enter API Key: ")
		key, err := readLine()
		if err != nil {
			return err
		}
		key = strings.TrimSpace(key)
		if key == "" {
			return fmt.Errorf("API key cannot be empty")
		}
		cfg.Providers.OpenAICompatible.APIKey = key
	case 3:
		cfg.Providers.OpenAICompatible.APIKey = ""
	case 4:
		// Keep current
	}

	fmt.Println()
	return nil
}

func configureStepFun(cfg *config.Config) error {
	fmt.Println("3. StepFun Setup")

	if cfg.Providers.StepFun.APIKey != "" {
		masked := maskAPIKey(cfg.Providers.StepFun.APIKey)
		fmt.Printf("   Current API Key: %s\n", masked)
	} else {
		fmt.Println("   Current API Key: Not configured")
	}

	fmt.Println("   Options:")
	fmt.Println("   [1] Enter/Update API Key")
	fmt.Println("   [2] Remove API Key (disable provider)")
	fmt.Println("   [3] Keep current (press ENTER)")

	choice, err := readChoiceOrKeep(1, 3, 3)
	if err != nil {
		return err
	}

	switch choice {
	case 1:
		fmt.Print("   Enter API Key: ")
		key, err := readLine()
		if err != nil {
			return err
		}
		key = strings.TrimSpace(key)
		if key == "" {
			// Keep existing if empty
			if cfg.Providers.StepFun.APIKey == "" {
				return fmt.Errorf("API key cannot be empty")
			}
		} else {
			cfg.Providers.StepFun.APIKey = key
		}
	case 2:
		cfg.Providers.StepFun.APIKey = ""
	case 3:
		// Keep current
	}

	fmt.Println()
	return nil
}

func configureDefaultModel(cfg *config.Config) error {
	fmt.Println("4. Default Model Selection")
	fmt.Println("   Fetching available models...")

	ctx := context.Background()
	availableModels, err := fetchAvailableModels(ctx, cfg)
	if err != nil {
		fmt.Printf("   Warning: Failed to fetch models: %v\n", err)
		fmt.Println("   Skipping model selection.")
		return nil
	}

	if len(availableModels) == 0 {
		fmt.Println("   No models available. Please configure a provider first.")
		return nil
	}

	// Sort models by provider then ID
	sort.Slice(availableModels, func(i, j int) bool {
		if availableModels[i].Provider != availableModels[j].Provider {
			return availableModels[i].Provider < availableModels[j].Provider
		}
		return availableModels[i].ID < availableModels[j].ID
	})

	fmt.Println("   Available models:")
	for i, m := range availableModels {
		fmt.Printf("   [%d] %s: %s\n", i+1, m.Provider, m.ID)
	}
	fmt.Println("   Press ENTER to keep current selection")

	choice, err := readChoiceOrKeep(1, len(availableModels), 0)
	if err != nil {
		return err
	}

	if choice == 0 {
		// Keep current
		if cfg.DefaultModel != nil {
			fmt.Printf("   Keeping: %s\n", cfg.DefaultModel.String())
		}
	} else {
		selected := availableModels[choice-1]
		cfg.DefaultModel = &config.ModelRef{
			Provider: selected.Provider,
			ModelID:  selected.ID,
		}
		fmt.Printf("   Selected: %s\n", selected.String())
	}

	fmt.Println()
	return nil
}

func fetchAvailableModels(ctx context.Context, cfg *config.Config) ([]llm.Model, error) {
	var allModels []llm.Model

	if cfg.Providers.OpenRouter.IsEnabled() {
		client := llm.NewOpenRouterClient(cfg.Providers.OpenRouter, cfg.Proxy, false)
		models, err := client.ListModels(ctx)
		if err != nil {
			return nil, fmt.Errorf("openrouter: %w", err)
		}
		allModels = append(allModels, models...)
	}

	if cfg.Providers.OpenAICompatible.IsEnabled() {
		client := llm.NewOpenAICompatibleClient(cfg.Providers.OpenAICompatible)
		models, err := client.ListModels(ctx)
		if err != nil {
			// OpenAI Compatible provider might be offline, don't fail
			fmt.Printf("   Warning: OpenAI Compatible provider not available: %v\n", err)
		} else {
			allModels = append(allModels, models...)
		}
	}

	if cfg.Providers.StepFun.IsEnabled() {
		client := llm.NewStepFunClient(cfg.Providers.StepFun, cfg.Proxy, false)
		models, err := client.ListModels(ctx)
		if err != nil {
			// StepFun provider might be offline, don't fail
			fmt.Printf("   Warning: StepFun provider not available: %v\n", err)
		} else {
			allModels = append(allModels, models...)
		}
	}

	return allModels, nil
}

func maskAPIKey(key string) string {
	if len(key) <= 8 {
		return "****"
	}
	return "****" + key[len(key)-4:]
}

// readChoiceOrKeep reads a choice with default value on empty input
func readChoiceOrKeep(min, max, defaultVal int) (int, error) {
	fmt.Print("   Choice: ")
	line, err := reader.ReadString('\n')
	if err != nil {
		return 0, fmt.Errorf("failed to read input: %w", err)
	}

	line = strings.TrimSpace(line)
	// Remove BOM if present
	line = strings.TrimPrefix(line, "\uFEFF")

	// Empty input returns default
	if line == "" {
		return defaultVal, nil
	}

	choice, err := strconv.Atoi(line)
	if err != nil {
		return 0, fmt.Errorf("invalid choice: %s", line)
	}

	if choice < min || choice > max {
		return 0, fmt.Errorf("choice out of range: %d (must be %d-%d)", choice, min, max)
	}

	return choice, nil
}

func readChoice(min, max int) (int, error) {
	return readChoiceOrKeep(min, max, 0)
}

func readLine() (string, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read input: %w", err)
	}
	line = strings.TrimSpace(line)
	line = strings.TrimPrefix(line, "\uFEFF")
	return line, nil
}
