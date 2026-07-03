package powershell

import (
	"fmt"
	"os"
	"os/exec"
)

// Execute runs a PowerShell command, streaming output and errors
// directly to the host console (os.Stdout / os.Stderr).
func Execute(command string) error {
	cmd := exec.Command("powershell", "-NoProfile", "-Command", command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("command execution failed: %w", err)
	}

	return nil
}
