package convert

import (
	"fmt"
)

// Options holds conversion settings
type Options struct {
	InputFile  string
	OutputFile string
	From       string
	To         string
	DryRun     bool
}

// Run executes a conversion job
func Run(opts Options) error {
	// Validate formats
	if err := ValidateFormats(opts.From, opts.To); err != nil {
		return err
	}

	// Resolve paths and check for overwrite
	in, out, err := ResolvePaths(opts)
	if err != nil {
		return err
	}
	opts.InputFile, opts.OutputFile = in, out

	// Dry-run
	if opts.DryRun {
		fmt.Printf("[Dry-run] Would convert %s (%s) -> %s (%s)\n",
			opts.InputFile, opts.From, opts.OutputFile, opts.To)
		return nil
	}

	// Get handlers using GetFormat() for case-insensitive lookup
	fromHandler, ok := GetFormat(opts.From)
	if !ok {
		return fmt.Errorf("no reader registered for format: %s", opts.From)
	}
	toHandler, ok := GetFormat(opts.To)
	if !ok {
		return fmt.Errorf("no writer registered for format: %s", opts.To)
	}

	// Read input
	data, err := fromHandler.ReaderFn(opts.InputFile)
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}

	// Write output
	if err := toHandler.WriterFn(opts.OutputFile, data); err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}

	fmt.Printf("Successfully converted %s (%s) -> %s (%s)\n",
		opts.InputFile, opts.From, opts.OutputFile, opts.To)

	return nil
}
