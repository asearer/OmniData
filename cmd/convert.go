package cmd

/*
Convert command for OmniData CLI.

Responsibilities:
- Convert structured data between supported formats.
- Support input/output files or STDIN/STDOUT.
- Support dry-run mode for previewing conversions.
- Validate formats and file existence before conversion.
- Extensible for future formats like XML, Excel, etc.
*/

import (
	"fmt"
	"os"
	"strings"

	"omnidata/internal/data" // Correct module import

	"github.com/spf13/cobra"
)

var (
	inputFile  string
	outputFile string
	fromFormat string
	toFormat   string
	dryRun     bool
)

// FormatHandler defines metadata for a supported format
type FormatHandler struct {
	Name string
}

// Supported formats map
var supportedFormats = map[string]FormatHandler{
	"csv":  {Name: "csv"},
	"json": {Name: "json"},
	// Future formats:
	// "xml":  {Name: "xml"},
	// "xlsx": {Name: "xlsx"},
}

// validateFormats checks if the source and target formats are supported
func validateFormats(from, to string) error {
	if _, ok := supportedFormats[from]; !ok {
		return fmt.Errorf("unsupported source format: %s", from)
	}
	if _, ok := supportedFormats[to]; !ok {
		return fmt.Errorf("unsupported target format: %s", to)
	}
	if from == to {
		return fmt.Errorf("source and target formats are the same")
	}
	return nil
}

// convertCmd defines the "convert" subcommand
var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Convert between data formats",
	Long: `Convert structured data between formats.
Currently supports CSV <-> JSON, extensible to other formats.`,
	Run: func(cmd *cobra.Command, args []string) {
		fromFormat = strings.ToLower(fromFormat)
		toFormat = strings.ToLower(toFormat)

		// Validate formats
		if err := validateFormats(fromFormat, toFormat); err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}

		// Handle STDIN/STDOUT
		if inputFile == "-" {
			inputFile = "/dev/stdin"
		} else {
			if _, err := os.Stat(inputFile); os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "Input file does not exist: %s\n", inputFile)
				os.Exit(1)
			}
		}

		if outputFile == "-" {
			outputFile = "/dev/stdout"
		} else if !dryRun {
			f, err := os.Create(outputFile)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Cannot create output file: %v\n", err)
				os.Exit(1)
			}
			f.Close()
		}

		// Dry-run mode
		if dryRun {
			fmt.Printf("[Dry-run] Would convert %s (%s) -> %s (%s)\n",
				inputFile, fromFormat, outputFile, toFormat)
			return
		}

		// Perform conversion
		if err := data.Convert(inputFile, outputFile, fromFormat, toFormat); err != nil {
			fmt.Fprintf(os.Stderr, "Conversion failed: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Successfully converted %s (%s) -> %s (%s)\n",
			inputFile, fromFormat, outputFile, toFormat)
	},
}

func init() {
	// Attach convertCmd as a subcommand of rootCmd
	rootCmd.AddCommand(convertCmd)

	// Define flags
	convertCmd.Flags().StringVarP(&inputFile, "input", "i", "", "Input file path ('-' for STDIN)")
	convertCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output file path ('-' for STDOUT)")
	convertCmd.Flags().StringVar(&fromFormat, "from", "csv", "Source format (csv/json)")
	convertCmd.Flags().StringVar(&toFormat, "to", "json", "Target format (csv/json)")
	convertCmd.Flags().BoolVarP(&dryRun, "dry-run", "d", false, "Preview conversion without writing output")

	// Mark required flags
	convertCmd.MarkFlagRequired("input")
	convertCmd.MarkFlagRequired("output")
}
