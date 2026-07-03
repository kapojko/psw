package cli

import (
	"fmt"
	"os/exec"
	"strings"
)

// CopyToClipboard copies text to the Windows clipboard
func CopyToClipboard(text string) error {
	// Use PowerShell's Set-Clipboard cmdlet
	cmd := exec.Command("powershell", "-NoProfile", "-Command", "Set-Clipboard -Value $input")
	cmd.Stdin = strings.NewReader(text)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("clipboard error: %s: %w", string(output), err)
	}

	return nil
}
