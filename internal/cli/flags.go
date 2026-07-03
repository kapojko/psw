package cli

// Flags holds all command-line flags
type Flags struct {
	Model    string
	Setup    bool
	Copy     bool
	Question bool
	Exec     bool
}

var flags Flags
