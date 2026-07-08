package cli

// Flags holds all command-line flags
type Flags struct {
	Model    string
	Setup    bool
	Copy     bool
	Question bool
	Exec     bool
	Verbose  bool
}

var flags Flags
