package cmd

/*
Root command for OmniData CLI.

Responsibilities:
- Define the root command of the CLI.
- Provide global flags (e.g., version).
- Allow subcommands (like "convert") to attach via init().
- Ensure clean exit with appropriate messaging.
*/

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd is the base command of OmniData CLI.
// All subcommands (convert, list, delete, etc.) attach here.
var rootCmd = &cobra.Command{
	Use:   "omnidata",
	Short: "OmniData - Universal Data Translator",
	Long: `OmniData is a CLI tool for converting, validating, and analyzing structured data.
Supports multiple formats like CSV, JSON, XML, XLSX, and is easily extensible.`,
	// SilenceErrors prevents Cobra from printing default errors, allowing custom formatting
	SilenceErrors: true,
	// SilenceUsage prevents Cobra from printing usage on every error
	SilenceUsage: true,
}

// Execute runs the root command; called from main().
func Execute() {
	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		// Print the error in a user-friendly way
		fprintf, err := fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		if err != nil {
			return
		}
		// Exit with a non-zero status code to indicate failure
		os.Exit(1)
	}
}

func init() {
	// ---------------------------
	// Global Persistent Flags
	// ---------------------------
	rootCmd.PersistentFlags().BoolP("version", "v", false, "Show OmniData version")

	// PersistentPreRun executes before any subcommand
	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		// Handle version flag
		versionFlag, err := cmd.Flags().GetBool("version")
		if err != nil {
			// Unexpected error retrieving the flag
			fprintf, err := fmt.Fprintf(os.Stderr, "Failed to read version flag: %v\n", err)
			if err != nil {
				return
			}
			os.Exit(1)
		}
		if versionFlag {
			fmt.Println("OmniData CLI v1.0.0") // Version output
			os.Exit(0)
		}
	}

	// ---------------------------
	// Placeholder for subcommand attachment
	// Example:
	// rootCmd.AddCommand(convertCmd)
	// rootCmd.AddCommand(listCmd)
	// ---------------------------
}
