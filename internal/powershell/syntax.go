package powershell

import (
	"fmt"
	"os/exec"
	"strings"
)

// CheckSyntax validates PowerShell command syntax without executing it.
// Uses [System.Management.Automation.Language.Parser]::ParseInput() to parse
// the command into an AST and return any parse errors.
func CheckSyntax(command string) (bool, []string) {
	psScript := fmt.Sprintf(`
$errors = $null
[System.Management.Automation.Language.Parser]::ParseInput(%s, [ref]$null, [ref]$errors) | Out-Null
if ($errors.Count -gt 0) {
    $errors | ForEach-Object { $_.Message }
}
`, psStringLiteral(command))

	cmd := exec.Command("powershell", "-NoProfile", "-Command", psScript)
	output, err := cmd.CombinedOutput()
	if err != nil {
		// PowerShell process error — treat as invalid syntax
		return false, []string{fmt.Sprintf("syntax check failed: %s", strings.TrimSpace(string(output)))}
	}

	outputStr := strings.TrimSpace(string(output))
	if outputStr == "" {
		return true, nil
	}

	// Parse errors were returned
	return false, strings.Split(outputStr, "\n")
}

// psStringLiteral wraps a string as a PowerShell single-quoted string literal,
// escaping embedded single quotes by doubling them.
func psStringLiteral(s string) string {
	escaped := strings.ReplaceAll(s, "'", "''")
	return "'" + escaped + "'"
}
