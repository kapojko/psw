# psw - PowerShell What?

A CLI utility that helps you with PowerShell commands. Ask a question in natural language and get the correct PowerShell command.

## Installation

```bash
go install github.com/kapojko/psw/cmd/psw@latest
```

Or build from source:

```bash
git clone https://github.com/kapojko/psw.git
cd psw
go build ./cmd/psw
```

## Quick Start

1. Run setup wizard:
   ```powershell
   psw -s
   ```

2. Ask questions:
   ```powershell
   psw "./dir1" get used size
   psw list all files in current directory
   psw how to zip a folder
   ```

## Usage

```
psw [flags] [prompt...]
```

### Flags

- `-s, --setup` - Run interactive setup wizard
- `-c, --copy` - Copy command to clipboard
- `-m, --model string` - Override default model (format: `provider/model`)
- `-h, --help` - Show help

### Examples

```powershell
# Get disk usage
psw "./dir1" get used size

# List files
psw list all files recursively

# Copy command to clipboard
psw -c "list files in current directory"

# Copy last command to clipboard (no prompt needed)
psw -c

# Use specific model
psw -m openrouter/anthropic/claude-3.5-sonnet how to rename multiple files
```

## Configuration

Run `psw -s` to set up:

1. **Proxy** - Optional HTTP proxy for OpenRouter API
2. **OpenRouter** - Enter API key from [openrouter.ai](https://openrouter.ai)
3. **LM Studio** - Enable if running [LM Studio](https://lmstudio.ai) locally
4. **Default Model** - Select from available models

Press ENTER at any prompt to keep the existing value.

Config is saved to `%APPDATA%/psw/config.json`

### Providers

| Provider | Setup |
|----------|-------|
| OpenRouter | API key from openrouter.ai |
| LM Studio | Local server on localhost:1234 |

## How It Works

1. Your question is combined with a system prompt optimized for PowerShell
2. The prompt is sent to the configured LLM provider
3. The response (PowerShell command + brief explanation) is displayed

## License

MIT
