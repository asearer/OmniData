package cmd

import (
	"fmt"

	"omnidata/internal/convert"
	"omnidata/internal/inspect"
	"omnidata/internal/output"

	"github.com/spf13/cobra"
)

var (
	peekInputFile  string
	peekFormat     string
	peekRows       int
	peekShowStats  bool
	peekOutputFile string
	peekOutputFmt  string
)

// peekCmd defines the "peek" subcommand for the CLI.
var peekCmd = &cobra.Command{
	Use:   "peek",
	Short: "Preview data and show schema information",
	Long: `Preview the first rows of a data file and display schema information.
Shows column names, types, and statistics.`,
	Example: `
  omnidata peek -i data.csv --format csv
  omnidata peek -i data.json --format json --rows 10 --stats
  cat data.csv | omnidata peek -i - --format csv`,
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := inspect.PeekOptions{
			InputFile: peekInputFile,
			Format:    peekFormat,
			Rows:      peekRows,
			ShowStats: peekShowStats,
		}

		// If output format is specified, use formatter
		if peekOutputFmt != "" {
			return runPeekWithOutput(opts, peekOutputFmt, peekOutputFile)
		}

		if err := inspect.RunPeek(opts); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(peekCmd)

	peekCmd.Flags().StringVarP(&peekInputFile, "input", "i", "", "Input file path ('-' for STDIN)")
	peekCmd.Flags().StringVar(&peekFormat, "format", "", "Data format (csv/json/xml/xlsx)")
	peekCmd.Flags().IntVarP(&peekRows, "rows", "n", 5, "Number of preview rows to show")
	peekCmd.Flags().BoolVar(&peekShowStats, "stats", false, "Show detailed column statistics")
	peekCmd.Flags().StringVarP(&peekOutputFile, "output", "o", "", "Output file path (optional, '-' for STDOUT)")
	peekCmd.Flags().StringVar(&peekOutputFmt, "output-format", "", "Output format (markdown/html/json)")

	err := peekCmd.MarkFlagRequired("input")
	if err != nil {
		return
	}
	err = peekCmd.MarkFlagRequired("format")
	if err != nil {
		return
	}
}

func runPeekWithOutput(opts inspect.PeekOptions, outputFormat, outputFile string) error {
	// Get format handler
	handler, ok := convert.GetFormat(opts.Format)
	if !ok {
		return fmt.Errorf("unsupported format: %s", opts.Format)
	}

	// Read data and infer schema
	data, err := handler.ReaderFn(opts.InputFile)
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}

	schema, err := inspect.InferSchema(data, opts.Format)
	if err != nil {
		return fmt.Errorf("failed to infer schema: %w", err)
	}

	// Get formatter
	formatter, err := output.GetFormatter(outputFormat)
	if err != nil {
		return err
	}

	// Format schema
	content, err := formatter.FormatSchema(schema)
	if err != nil {
		return fmt.Errorf("failed to format schema: %w", err)
	}

	// Write output
	return output.WriteOutput(content, outputFile)
}
