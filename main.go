package main

/*
OmniData - Universal Data Translator

Entry point for the OmniData CLI tool.

Responsibilities:
- Execute the root command of the CLI.
- Handle errors gracefully and exit with a non-zero status code if needed.
- Ensure proper initialization of subcommands and format handlers.
*/

import (
	"omnidata/cmd" // Import the CLI commands using the module path
)

func main() {
	// Execute the root command (defined in cmd/root.go)
	// Errors are handled internally by cmd.Execute()
	cmd.Execute()
}
