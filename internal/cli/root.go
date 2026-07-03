package cli

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/kapojko/psw/internal/config"
	"github.com/kapojko/psw/internal/llm"
	"github.com/kapojko/psw/internal/prompt"
)

// NewRootCommand creates the root cobra command
func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "psw [flags] [prompt...]",
		Short: "PowerShell assistant - ask questions about PowerShell commands",
		Long: `psw (powershell what?) is a CLI tool that helps with PowerShell commands.
It sends your question to an LLM and returns the appropriate PowerShell command.

Examples:
  psw "./dir1" get used size
  psw list all files in current directory
  psw -c "list files in current directory"
  psw -c                    # Copy last command to clipboard
  psw -m openrouter/anthropic/claude-3.5-sonnet how to zip a folder`,
		Args: cobra.ArbitraryArgs,
		RunE: runRoot,
	}

	cmd.Flags().StringVarP(&flags.Model, "model", "m", "", "Override default model (format: provider/model)")
	cmd.Flags().BoolVarP(&flags.Setup, "setup", "s", false, "Run interactive setup wizard")
	cmd.Flags().BoolVarP(&flags.Copy, "copy", "c", false, "Copy command to clipboard")

	// Add --help flag for compatibility
	cmd.Flags().BoolP("help", "h", false, "Help for psw")

	return cmd
}

func runRoot(cmd *cobra.Command, args []string) error {
	if flags.Setup {
		return RunConfig()
	}

	// Load config
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Handle -c with empty prompt: copy last command
	userPrompt := strings.Join(args, " ")
	if flags.Copy && userPrompt == "" {
		if cfg.LastRequest == nil || cfg.LastRequest.Command == "" {
			return fmt.Errorf("no previous request found. Run a query first")
		}
		if err := CopyToClipboard(cfg.LastRequest.Command); err != nil {
			return fmt.Errorf("failed to copy to clipboard: %w", err)
		}
		fmt.Println(cfg.LastRequest.Command)
		fmt.Println("[Copied to clipboard]")
		return nil
	}

	if userPrompt == "" {
		return fmt.Errorf("no prompt provided. Usage: psw [flags] [prompt...]")
	}

	// Determine model to use
	var modelRef config.ModelRef
	if flags.Model != "" {
		modelRef, err = config.ParseModelRef(flags.Model)
		if err != nil {
			return fmt.Errorf("invalid model flag: %w", err)
		}
	} else if cfg.DefaultModel != nil {
		modelRef = *cfg.DefaultModel
	} else {
		return fmt.Errorf("no model specified and no default model configured. Run 'psw -s' to set up")
	}

	// Get provider config
	providerCfg := cfg.Providers.GetProvider(modelRef.Provider)
	if providerCfg == nil {
		return fmt.Errorf("provider %s not configured", modelRef.Provider)
	}

	if !providerCfg.IsEnabled() {
		return fmt.Errorf("provider %s is not enabled. Run 'psw -s' to enable it", modelRef.Provider)
	}

	// Create LLM client
	client, err := llm.NewClient(providerCfg, cfg.Proxy)
	if err != nil {
		return fmt.Errorf("failed to create LLM client: %w", err)
	}

	// Build messages
	messages := prompt.BuildMessages(userPrompt)

	// Send request
	ctx := context.Background()
	response, err := client.ChatCompletion(ctx, modelRef.ModelID, messages)
	if err != nil {
		return fmt.Errorf("LLM request failed: %w", err)
	}

	// Parse response to extract command and explanation
	command, explanation := prompt.ParseResponse(response)

	// Save last request/response
	cfg.LastRequest = &config.LastRequest{
		Prompt:  userPrompt,
		Command: command,
	}
	if saveErr := cfg.Save(); saveErr != nil {
		// Non-fatal: just warn
		fmt.Fprintf(os.Stderr, "Warning: failed to save last request: %v\n", saveErr)
	}

	// Display response
	if flags.Copy {
		// Copy command to clipboard
		if err := CopyToClipboard(command); err != nil {
			return fmt.Errorf("failed to copy to clipboard: %w", err)
		}
		fmt.Println(command)
		fmt.Println("[Copied to clipboard]")
	} else {
		// Display full response
		if explanation != "" {
			fmt.Printf("%s\n\n%s\n", command, explanation)
		} else {
			fmt.Println(command)
		}
	}

	return nil
}

// Execute runs the CLI
func Execute() {
	cmd := NewRootCommand()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
