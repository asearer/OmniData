package convert

import (
	"fmt"
)

/*
Options holds all settings for a conversion job.

Fields:
- InputFile: path to the input file; use "-" for STDIN.
- OutputFile: path to the output file; use "-" for STDOUT.
- From: source format name (csv, json, xml, xlsx).
- To: target format name (csv, json, xml, xlsx).
- DryRun: if true, simulates conversion without writing output.
- Stream: if true, uses streaming mode for large files (memory-efficient).
*/
type Options struct {
	InputFile  string
	OutputFile string
	From       string
	To         string
	DryRun     bool
	Stream     bool
}

/*
Run executes a conversion job based on Options.

Responsibilities:
- Validate formats and file paths.
- Handle dry-run simulations.
- Read from the source format and write to the target format.
- Wrap errors with detailed context.
- Support cross-platform STDIN/STDOUT.
*/
func Run(opts Options) error {
	// ---------------------------
	// Step 1: Validate formats
	// ---------------------------
	if err := ValidateFormats(opts.From, opts.To); err != nil {
		return fmt.Errorf("invalid format selection: %w", err)
	}

	// ---------------------------
	// Step 2: Resolve paths
	// ---------------------------
	inputPath, outputPath, err := ResolvePaths(opts)
	if err != nil {
		return fmt.Errorf("failed to resolve paths: %w", err)
	}
	opts.InputFile = inputPath
	opts.OutputFile = outputPath

	// ---------------------------
	// Step 3: Dry-run mode
	// ---------------------------
	if opts.DryRun {
		// Simulate reading input to detect early errors
		fromHandler, ok := GetFormat(opts.From)
		if !ok {
			return fmt.Errorf("[dry-run] no reader registered for format: %s", opts.From)
		}
		if _, err := fromHandler.ReaderFn(opts.InputFile); err != nil {
			return fmt.Errorf("[dry-run] failed to read input: %w", err)
		}

		fmt.Printf("[Dry-run] Conversion simulation succeeded: %s (%s) -> %s (%s)\n",
			opts.InputFile, opts.From, opts.OutputFile, opts.To)
		return nil
	}

	// ---------------------------
	// Step 4: Get format handlers
	// ---------------------------
	fromHandler, ok := GetFormat(opts.From)
	if !ok {
		return fmt.Errorf("no reader registered for format: %s", opts.From)
	}
	toHandler, ok := GetFormat(opts.To)
	if !ok {
		return fmt.Errorf("no writer registered for format: %s", opts.To)
	}

	// ---------------------------
	// Step 5: Read input data
	// ---------------------------
	data, err := fromHandler.ReaderFn(opts.InputFile)
	if err != nil {
		return fmt.Errorf("failed to read input file '%s': %w", opts.InputFile, err)
	}

	// ---------------------------
	// Step 6: Write output data
	// ---------------------------
	if err := toHandler.WriterFn(opts.OutputFile, data); err != nil {
		return fmt.Errorf("failed to write output file '%s': %w", opts.OutputFile, err)
	}

	// ---------------------------
	// Step 7: Success message
	// ---------------------------
	fmt.Printf("Successfully converted %s (%s) -> %s (%s)\n",
		opts.InputFile, opts.From, opts.OutputFile, opts.To)

	return nil
}
