package cmd

/*
Root command for OmniData CLI.

Responsibilities:
- Define the root command of the CLI.
- Provide global flags (e.g., version).
- Allow subcommands (like "convert") to attach via init().
*/

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd is the base command of OmniData
var rootCmd = &cobra.Command{
	Use:   "OmniData",
	Short: "OmniData - Universal Data Translator",
	Long: `OmniData is a CLI tool for converting, validating, and analyzing structured data.
It supports multiple formats like CSV, JSON, and can be extended to XML, Excel, etc.`,
}

// Execute runs the root command; called from main()
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Global persistent version flag
	rootCmd.PersistentFlags().BoolP("version", "v", false, "Show OmniData version")

	// PersistentPreRun executes before any subcommand
	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		versionFlag, _ := cmd.Flags().GetBool("version")
		if versionFlag {
			fmt.Println("OmniData CLI v1.0.0") // Version info
			os.Exit(0)
		}
	}
}
