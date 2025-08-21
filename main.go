package main

/*
OmniData - Universal Data Translator

Entry point for the OmniData CLI tool.

Responsibilities:
- Execute the root command of the CLI.
- Handle errors gracefully and exit with a non-zero status code if needed.
*/

import (
	"fmt"
	"os"

	"omnidata/cmd" // Import the CLI commands using the module path
)

func main() {
	// Execute the root command (defined in cmd/root.go)
	if err := cmd.Execute(); err != nil {
		// Print errors to STDERR
		fmt.Fprintln(os.Stderr, "Error:", err)
		// Exit with non-zero status code
		os.Exit(1)
	}
}
