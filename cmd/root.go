package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "omnidata",
	Short: "OmniData - Universal Data Translator",
	Long: `OmniData is a CLI tool for converting, validating, and analyzing structured data.
Supports multiple formats like CSV, JSON, XML, XLSX, and is easily extensible.`,
	SilenceErrors: true,
	SilenceUsage:  true,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		// Just print the error — no need to assign to fprintf
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolP("version", "v", false, "Show OmniData version")

	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		versionFlag, err := cmd.Flags().GetBool("version")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to read version flag: %v\n", err)
			os.Exit(1)
		}
		if versionFlag {
			fmt.Println("OmniData CLI v1.0.0")
			os.Exit(0)
		}
	}

	// Attach subcommands here
	// rootCmd.AddCommand(convertCmd)
	// rootCmd.AddCommand(listCmd)
}
