package prompt

import (
	"strings"

	"github.com/kapojko/psw/internal/llm"
)

const defaultSystemPrompt = `You are a PowerShell assistant. The user is working in a PowerShell command-line environment on Windows.

CRITICAL: You MUST respond in EXACTLY this format:

<command>
[the PowerShell command here]
</command>
<explanation>
[brief 1-2 sentence explanation if command is complex, otherwise leave empty]
</explanation>

Rules:
- Provide correct PowerShell syntax for commands
- Use PowerShell-native cmdlets and syntax (not cmd.exe or bash)
- For file paths, assume Windows-style paths
- Keep explanations brief (1-2 sentences max)
- If command is simple, explanation can be empty but tags must still be present
- Never wrap commands in code blocks - use the XML tags above
- Never add any text before <command> or after </explanation>`

const questionSystemPrompt = `You are a helpful assistant answering questions in a terminal environment.

Rules:
- Keep answers brief and concise (2-4 sentences max)
- Terminal is not suitable for reading long texts
- Be direct and to the point
- Use plain text, avoid markdown formatting
- If showing code, keep it minimal`

// GetSystemPrompt returns the system prompt for PowerShell commands
func GetSystemPrompt() string {
	return defaultSystemPrompt
}

// GetQuestionPrompt returns the system prompt for general questions
func GetQuestionPrompt() string {
	return questionSystemPrompt
}

// BuildMessages creates the message array for the LLM request (PowerShell mode)
func BuildMessages(userPrompt string) []llm.Message {
	return []llm.Message{
		{
			Role:    "system",
			Content: GetSystemPrompt(),
		},
		{
			Role:    "user",
			Content: userPrompt,
		},
	}
}

// BuildQuestionMessages creates the message array for general questions
func BuildQuestionMessages(userPrompt string) []llm.Message {
	return []llm.Message{
		{
			Role:    "system",
			Content: GetQuestionPrompt(),
		},
		{
			Role:    "user",
			Content: userPrompt,
		},
	}
}

// ParseResponse extracts command and explanation from structured LLM response
func ParseResponse(response string) (command, explanation string) {
	// Extract command from <command>...</command> tags
	command = extractTag(response, "command")
	explanation = extractTag(response, "explanation")

	// Fallback: if no tags found, use entire response as command
	if command == "" {
		command = strings.TrimSpace(response)
	}

	return command, strings.TrimSpace(explanation)
}

func extractTag(text, tagName string) string {
	startTag := "<" + tagName + ">"
	endTag := "</" + tagName + ">"

	startIdx := strings.Index(text, startTag)
	if startIdx == -1 {
		return ""
	}

	startIdx += len(startTag)
	endIdx := strings.Index(text[startIdx:], endTag)
	if endIdx == -1 {
		return ""
	}

	return strings.TrimSpace(text[startIdx : startIdx+endIdx])
}
