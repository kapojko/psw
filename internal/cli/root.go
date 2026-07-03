package cli

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/kapojko/psw/internal/config"
	"github.com/kapojko/psw/internal/llm"
	"github.com/kapojko/psw/internal/powershell"
	"github.com/kapojko/psw/internal/prompt"
)

var (
	commandColor        = color.New(color.FgHiMagenta) // orange (bright yellow)
	descriptionColor    = color.New(color.FgWhite)     // dark grey
	copiedColor         = color.New(color.FgHiBlue)    // medium grey
	invalidCommandColor = color.New(color.FgRed)       // red for invalid syntax
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
  psw -q "what is GOPATH?"
  psw -m openrouter/anthropic/claude-3.5-sonnet how to zip a folder
  psw -e "list files"       # Get command and optionally execute it`,
		Args: cobra.ArbitraryArgs,
		RunE: runRoot,
	}

	cmd.Flags().StringVarP(&flags.Model, "model", "m", "", "Override default model (format: provider/model)")
	cmd.Flags().BoolVarP(&flags.Setup, "setup", "s", false, "Run interactive setup wizard")
	cmd.Flags().BoolVarP(&flags.Copy, "copy", "c", false, "Copy command to clipboard")
	cmd.Flags().BoolVarP(&flags.Question, "question", "q", false, "General question mode (not PowerShell-specific)")
	cmd.Flags().BoolVarP(&flags.Exec, "exec", "e", false, "Execute the command after syntax check")

	// Add --help flag for compatibility
	cmd.Flags().BoolP("help", "h", false, "Help for psw")

	return cmd
}

func runRoot(cmd *cobra.Command, args []string) error {
	if flags.Setup {
		return RunConfig()
	}

	// Validate mutually exclusive flags
	if flags.Copy && flags.Exec {
		return fmt.Errorf("cannot use --copy (-c) and --exec (-e) together")
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
		printOutput(cfg.LastRequest.Command, "", true, true, nil)
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

	// Build messages based on mode
	var messages []llm.Message
	if flags.Question {
		messages = prompt.BuildQuestionMessages(userPrompt)
	} else {
		messages = prompt.BuildMessages(userPrompt)
	}

	// Send request
	ctx := context.Background()
	response, err := client.ChatCompletion(ctx, modelRef.ModelID, messages)
	if err != nil {
		return fmt.Errorf("LLM request failed: %w", err)
	}

	// Handle response based on mode
	if flags.Question {
		// Question mode: display response as-is
		fmt.Println(response)
	} else {
		// PowerShell mode: parse and display command/explanation
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

		// Check syntax
		syntaxValid, syntaxErrors := powershell.CheckSyntax(command)

		if flags.Copy {
			// Copy command to clipboard
			if err := CopyToClipboard(command); err != nil {
				return fmt.Errorf("failed to copy to clipboard: %w", err)
			}
			printOutput(command, explanation, true, syntaxValid, syntaxErrors)
		} else if flags.Exec {
			// Execute mode: display command with syntax info
			printOutput(command, explanation, false, syntaxValid, syntaxErrors)

			if !syntaxValid {
				return nil
			}

			// Prompt for execution
			fmt.Print("Execute command? [e/y to execute, any other to cancel]: ")
			reader := bufio.NewReader(os.Stdin)
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(strings.ToLower(input))

			if input == "e" || input == "y" {
				fmt.Println()
				if err := powershell.Execute(command); err != nil {
					return fmt.Errorf("execution failed: %w", err)
				}
			}
		} else {
			// Display full response
			printOutput(command, explanation, false, syntaxValid, syntaxErrors)
		}
	}

	return nil
}

func printOutput(command, explanation string, copied bool, syntaxValid bool, syntaxErrors []string) {
	if !syntaxValid {
		invalidCommandColor.Print("! Syntax error")
		if len(syntaxErrors) > 0 {
			invalidCommandColor.Printf(": %s", strings.TrimSpace(syntaxErrors[0]))
		}
		fmt.Println()
		invalidCommandColor.Println(command)
	} else {
		commandColor.Println(command)
	}
	if copied {
		copiedColor.Println("[Copied to clipboard]")
	}
	if explanation != "" {
		fmt.Println()
		descriptionColor.Println(explanation)
	}
}

// Execute runs the CLI
func Execute() {
	cmd := NewRootCommand()
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
