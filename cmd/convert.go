package cmd

import (
	"strings"

	"omnidata/internal/convert"

	"github.com/spf13/cobra"
)

var (
	inputFile  string
	outputFile string
	fromFormat string
	toFormat   string
	dryRun     bool
)

var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Convert between data formats",
	Long: `Convert structured data between formats.
Currently supports CSV <-> JSON, extensible to other formats.`,
	Example: `
  omnidata convert -i data.csv -o data.json --from csv --to json
  cat data.csv | omnidata convert -i - -o - --from csv --to json
  omnidata convert --list-formats`,
	RunE: func(cmd *cobra.Command, args []string) error {
		opts := convert.Options{
			InputFile:  inputFile,
			OutputFile: outputFile,
			From:       strings.ToLower(fromFormat),
			To:         strings.ToLower(toFormat),
			DryRun:     dryRun,
		}
		return convert.Run(opts)
	},
}

func init() {
	rootCmd.AddCommand(convertCmd)

	convertCmd.Flags().StringVarP(&inputFile, "input", "i", "", "Input file path ('-' for STDIN)")
	convertCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file path ('-' for STDOUT)")
	convertCmd.Flags().StringVar(&fromFormat, "from", "", "Source format (csv/json)")
	convertCmd.Flags().StringVar(&toFormat, "to", "", "Target format (csv/json)")
	convertCmd.Flags().BoolVarP(&dryRun, "dry-run", "d", false, "Preview conversion without writing output")

	convertCmd.MarkFlagRequired("input")
	convertCmd.MarkFlagRequired("output")
}
