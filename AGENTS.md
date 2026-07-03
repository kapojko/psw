# psw - PowerShell What?

CLI utility that sends PowerShell-related questions to an LLM and returns commands.

## Dependencies

- github.com/spf13/cobra - CLI framework
- github.com/sashabaranov/go-openai - OpenAI-compatible API client

## Build & Run

```bash
go build ./cmd/psw
./psw --help
```

## Important Rule

When modifying code structure, always update `docs/CODE.md` to reflect changes.
