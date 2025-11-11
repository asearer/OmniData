package cmd

import (
	"strings"

	"omnidata/internal/convert"

	"github.com/spf13/cobra"
)

var (
	// inputFile and outputFile hold the paths for conversion.
	// Use "-" to indicate STDIN/STDOUT, respectively.
	inputFile  string
	outputFile string
	fromFormat string
	toFormat   string
	dryRun     bool
	stream     bool
)

// convertCmd defines the "convert" subcommand for the CLI.
// Responsibilities:
// - Dispatch conversion to the internal engine.
// - Handle dry-run mode.
// - Validate input/output paths and formats.
// - Provide helpful CLI usage and examples.
var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Convert between data formats",
	Long: `Convert structured data between supported formats.
Supports CSV, JSON, XML, XLSX (extensible).`,
	Example: `
  omnidata convert -i data.csv -o data.json --from csv --to json
  cat data.csv | omnidata convert -i - -o - --from csv --to json
  omnidata convert --list-formats`,
	// RunE allows returning errors to Cobra which prints them and exits with code 1
	RunE: func(cmd *cobra.Command, args []string) error {
		// Prepare conversion options
		opts := convert.Options{
			InputFile:  inputFile,
			OutputFile: outputFile,
			From:       strings.ToLower(fromFormat),
			To:         strings.ToLower(toFormat),
			DryRun:     dryRun,
			Stream:     stream,
		}

		// Delegate actual conversion to the internal convert engine
		// This will handle validation, reading, writing, and dry-run simulation
		if err := convert.Run(opts); err != nil {
			// Wrap the error with context before returning to Cobra
			return err
		}

		return nil
	},
}

func init() {
	// Register the convert subcommand under root
	rootCmd.AddCommand(convertCmd)

	// ---------------------------
	// Define CLI flags
	// ---------------------------
	convertCmd.Flags().StringVarP(&inputFile, "input", "i", "", "Input file path ('-' for STDIN)")
	convertCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file path ('-' for STDOUT)")
	convertCmd.Flags().StringVar(&fromFormat, "from", "", "Source format (csv/json/xml/xlsx)")
	convertCmd.Flags().StringVar(&toFormat, "to", "", "Target format (csv/json/xml/xlsx)")
	convertCmd.Flags().BoolVarP(&dryRun, "dry-run", "d", false, "Preview conversion without writing output")
	convertCmd.Flags().BoolVarP(&stream, "stream", "s", false, "Use streaming mode for large files (memory-efficient)")

	// Mark required flags for input/output
	err := convertCmd.MarkFlagRequired("input")
	if err != nil {
		return
	}
	err = convertCmd.MarkFlagRequired("output")
	if err != nil {
		return
	}
}
