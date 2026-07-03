package cli

// Flags holds all command-line flags
type Flags struct {
	Model  string
	Setup  bool
	Copy   bool
}

var flags Flags
