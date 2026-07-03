# psw Code Structure

## Directory Layout

```
psw/
├── cmd/psw/main.go           # Application entry point
├── internal/
│   ├── cli/                  # Command-line interface handling
│   │   ├── flags.go          # Flag definitions
│   │   ├── config.go         # Interactive configuration wizard
│   │   ├── clipboard.go      # Clipboard functionality
│   │   └── root.go           # Root command and execution logic
│   ├── config/               # Configuration management
│   │   ├── config.go         # Config struct, load/save
│   │   ├── providers.go      # Provider types and configuration
│   │   └── paths.go          # Config file path resolution
│   ├── llm/                  # LLM client implementations
│   │   ├── client.go         # Client interface definition
│   │   ├── openrouter.go     # OpenRouter provider
│   │   └── lmstudio.go       # LM Studio provider
│   ├── powershell/           # PowerShell interaction (no script execution)
│   │   ├── syntax.go         # Syntax validation via Parser API
│   │   └── execute.go        # Command execution with output streaming
│   └── prompt/               # Prompt construction
│       └── system.go         # System prompt and message building
├── docs/
│   └── CODE.md               # This file
├── AGENTS.md                 # Project overview for agents
├── README.md                 # User documentation
├── go.mod                    # Go module definition
└── go.sum                    # Dependency checksums
```

## Package Descriptions

### cmd/psw
Entry point. Imports and executes the CLI package.

### internal/cli
Handles all user interaction:
- **flags.go**: Defines CLI flags (--setup, --copy, --model, --exec, --help)
- **config.go**: Interactive setup wizard with 4 steps:
  0. Proxy setup (optional HTTP proxy for OpenRouter)
  1. OpenRouter API key setup
  2. LM Studio enable/disable
  3. Default model selection (fetches available models)
- **clipboard.go**: Copy command to Windows clipboard using PowerShell
- **root.go**: Main command execution - parses args, loads config, creates LLM client, sends request, parses response, checks syntax, optionally executes

### internal/config
Manages application configuration:
- **config.go**: Main Config struct with Load/Save methods. Config stored in `%APPDATA%/psw/config.json`
- **providers.go**: Provider types (OpenRouter, LM Studio), ModelRef type with parsing
- **paths.go**: Cross-platform config directory resolution

### internal/llm
LLM client abstraction:
- **client.go**: Client interface with ChatCompletion and ListModels methods
- **openrouter.go**: OpenRouter implementation using go-openai library with direct HTTP for model listing
- **lmstudio.go**: LM Studio implementation (local OpenAI-compatible API)

### internal/powershell
PowerShell interaction utilities — syntax validation and command execution:
- **syntax.go**: Validates PowerShell command syntax using `[System.Management.Automation.Language.Parser]::ParseInput()` without executing scripts
- **execute.go**: Executes PowerShell commands with real-time output streaming to host console

### internal/prompt
Prompt construction and response parsing:
- **system.go**: Default system prompt with structured XML output format, response parser

## Key Interfaces

### llm.Client
```go
type Client interface {
    ChatCompletion(ctx context.Context, model string, messages []Message) (string, error)
    ListModels(ctx context.Context) ([]Model, error)
}
```

### config.ProviderConfig
```go
type ProviderConfig interface {
    GetType() ProviderType
    IsEnabled() bool
    GetDisplayName() string
}
```

## Response Format

LLM responses use structured XML format:
```xml
<command>
Get-ChildItem -Recurse
</command>
<explanation>
Lists all files recursively in the current directory
</explanation>
```

## Adding a New Provider

1. Create `internal/llm/newprovider.go` implementing the `Client` interface
2. Add provider config struct in `internal/config/providers.go`
3. Update `ProvidersConfig` struct and `GetProvider()` method
4. Update `llm.NewClient()` factory function
5. Update `cli/config.go` to add provider configuration step

## Configuration Format

```json
{
  "proxy": {
    "enabled": true,
    "url": "http://proxy:8080"
  },
  "providers": {
    "openrouter": {
      "api_key": "sk-..."
    },
    "lmstudio": {
      "enabled": true,
      "base_url": "http://localhost:1234/v1"
    }
  },
  "default_model": {
    "provider": "openrouter",
    "model_id": "anthropic/claude-3.5-sonnet"
  }
}
```

## Dependencies

- **cobra**: CLI framework
- **go-openai**: OpenAI-compatible API client
