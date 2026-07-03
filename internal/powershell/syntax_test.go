package powershell

import (
	"testing"
)

func TestCheckSyntax_ValidCommand(t *testing.T) {
	valid, errs := CheckSyntax("Get-ChildItem")
	if !valid {
		t.Errorf("expected valid syntax, got errors: %v", errs)
	}
	if len(errs) != 0 {
		t.Errorf("expected no errors, got: %v", errs)
	}
}

func TestCheckSyntax_ValidCommandWithArgs(t *testing.T) {
	valid, errs := CheckSyntax("Get-ChildItem -Path C:\\Users -Recurse")
	if !valid {
		t.Errorf("expected valid syntax, got errors: %v", errs)
	}
}

func TestCheckSyntax_MultipleCommands(t *testing.T) {
	valid, errs := CheckSyntax("Get-ChildItem; Write-Host 'hello'")
	if !valid {
		t.Errorf("expected valid syntax, got errors: %v", errs)
	}
}

func TestCheckSyntax_Pipeline(t *testing.T) {
	valid, errs := CheckSyntax("Get-Process | Where-Object { $_.CPU -gt 100 }")
	if !valid {
		t.Errorf("expected valid syntax, got errors: %v", errs)
	}
}

func TestCheckSyntax_InvalidSyntax_MissingValue(t *testing.T) {
	valid, errs := CheckSyntax("{")
	if valid {
		t.Error("expected invalid syntax for unclosed brace")
	}
	if len(errs) == 0 {
		t.Error("expected parse errors, got none")
	}
}

func TestCheckSyntax_InvalidSyntax_UnclosedString(t *testing.T) {
	valid, errs := CheckSyntax("Write-Host 'unclosed string")
	if valid {
		t.Error("expected invalid syntax for unclosed string")
	}
	if len(errs) == 0 {
		t.Error("expected parse errors, got none")
	}
}

func TestCheckSyntax_InvalidSyntax_BadCmdlet(t *testing.T) {
	valid, _ := CheckSyntax(")invalid(")
	if valid {
		t.Error("expected invalid syntax for malformed input")
	}
}

func TestCheckSyntax_EmptyCommand(t *testing.T) {
	valid, _ := CheckSyntax("")
	if !valid {
		t.Error("empty command should be valid (no-op)")
	}
}

func TestPsStringLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"hello", "'hello'"},
		{"it's", "'it''s'"},
		{"a'b'c", "'a''b''c'"},
		{"", "''"},
	}

	for _, tt := range tests {
		result := psStringLiteral(tt.input)
		if result != tt.expected {
			t.Errorf("psStringLiteral(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}
