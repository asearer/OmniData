package cmd

import (
	"fmt"

	"omnidata/internal/convert"
	"omnidata/internal/inspect"
	"omnidata/internal/output"

	"github.com/spf13/cobra"
)

var (
	diffFile1      string
	diffFile2      string
	diffFormat1    string
	diffFormat2    string
	diffOutputFile string
	diffOutputFmt  string
)

// diffCmd defines the "diff" subcommand for the CLI.
var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Compare schemas between two data files",
	Long: `Compare the schemas of two data files and show differences.
Displays added, removed, and changed columns.`,
	Example: `
  omnidata diff -1 data1.csv -2 data2.csv --format1 csv --format2 csv
  omnidata diff -1 old.json -2 new.json --format1 json --format2 json`,
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := inspect.DiffOptions{
			File1:   diffFile1,
			File2:   diffFile2,
			Format1: diffFormat1,
			Format2: diffFormat2,
		}

		// If output format is specified, use formatter
		if diffOutputFmt != "" {
			return runDiffWithOutput(opts, diffOutputFmt, diffOutputFile)
		}

		if err := inspect.RunDiff(opts); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(diffCmd)

	diffCmd.Flags().StringVarP(&diffFile1, "file1", "1", "", "First file to compare")
	diffCmd.Flags().StringVarP(&diffFile2, "file2", "2", "", "Second file to compare")
	diffCmd.Flags().StringVar(&diffFormat1, "format1", "", "Format of first file (csv/json/xml/xlsx)")
	diffCmd.Flags().StringVar(&diffFormat2, "format2", "", "Format of second file (csv/json/xml/xlsx)")
	diffCmd.Flags().StringVarP(&diffOutputFile, "output", "o", "", "Output file path (optional, '-' for STDOUT)")
	diffCmd.Flags().StringVar(&diffOutputFmt, "output-format", "", "Output format (markdown/html/json)")

	diffCmd.MarkFlagRequired("file1")
	diffCmd.MarkFlagRequired("file2")
	diffCmd.MarkFlagRequired("format1")
	diffCmd.MarkFlagRequired("format2")
}

func runDiffWithOutput(opts inspect.DiffOptions, outputFormat, outputFile string) error {
	// Get format handlers
	handler1, ok := convert.GetFormat(opts.Format1)
	if !ok {
		return fmt.Errorf("unsupported format: %s", opts.Format1)
	}

	handler2, ok := convert.GetFormat(opts.Format2)
	if !ok {
		return fmt.Errorf("unsupported format: %s", opts.Format2)
	}

	// Read data from both files
	data1, err := handler1.ReaderFn(opts.File1)
	if err != nil {
		return fmt.Errorf("failed to read file1: %w", err)
	}

	data2, err := handler2.ReaderFn(opts.File2)
	if err != nil {
		return fmt.Errorf("failed to read file2: %w", err)
	}

	// Infer schemas
	schema1, err := inspect.InferSchema(data1, opts.Format1)
	if err != nil {
		return fmt.Errorf("failed to infer schema for file1: %w", err)
	}

	schema2, err := inspect.InferSchema(data2, opts.Format2)
	if err != nil {
		return fmt.Errorf("failed to infer schema for file2: %w", err)
	}

	// Compare schemas
	diff := inspect.CompareSchemas(schema1, schema2)

	// Get formatter
	formatter, err := output.GetFormatter(outputFormat)
	if err != nil {
		return err
	}

	// Format diff
	content, err := formatter.FormatDiff(diff, schema1, schema2)
	if err != nil {
		return fmt.Errorf("failed to format diff: %w", err)
	}

	// Write output
	return output.WriteOutput(content, outputFile)
}
