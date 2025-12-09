package convert

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"strings"
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
- Support transparent Gzip compression/decompression.
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
		// For SQL, we don't open the file
		if opts.From != "sql" {
			// Try opening check (simplified for dry-run)
			if opts.InputFile != "-" {
				f, err := os.Open(opts.InputFile)
				if err != nil {
					return fmt.Errorf("[dry-run] failed to read input: %w", err)
				}
				f.Close()
			}
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
	// Step 5: Prepare Input Reader
	// ---------------------------
	var reader io.ReadCloser
	if opts.From != "sql" {
		if opts.InputFile == "-" {
			reader = os.Stdin
		} else {
			f, err := os.Open(opts.InputFile)
			if err != nil {
				return fmt.Errorf("failed to open input file: %w", err)
			}
			reader = f
		}

		// Handle Gzip input
		if strings.HasSuffix(opts.InputFile, ".gz") {
			gzReader, err := gzip.NewReader(reader)
			if err != nil {
				reader.Close()
				return fmt.Errorf("failed to create gzip reader: %w", err)
			}
			// Wrap to close both
			reader = &readCloserWrapper{Reader: gzReader, Closer: reader}
		}
	} else {
		// For SQL, we don't provide a reader, the handler manages the connection string
		reader = nil
	}

	if reader != nil && opts.InputFile != "-" {
		defer reader.Close()
	}

	// ---------------------------
	// Step 6: Read input data
	// ---------------------------
	data, err := fromHandler.ReaderFn(reader, opts.InputFile)
	if err != nil {
		return fmt.Errorf("failed to read input '%s': %w", opts.InputFile, err)
	}

	// ---------------------------
	// Step 7: Prepare Output Writer
	// ---------------------------
	var writer io.WriteCloser
	if opts.To != "sql" {
		if opts.OutputFile == "-" {
			writer = os.Stdout
		} else {
			f, err := os.Create(opts.OutputFile)
			if err != nil {
				return fmt.Errorf("failed to create output file: %w", err)
			}
			writer = f
		}

		// Handle Gzip output
		if strings.HasSuffix(opts.OutputFile, ".gz") {
			gzWriter := gzip.NewWriter(writer)
			// Wrap to close both (gzip flush + file close)
			writer = &writeCloserWrapper{Writer: gzWriter, Closer: writer}
		}
	} else {
		// For SQL, writer is nil
		writer = nil
	}

	if writer != nil && opts.OutputFile != "-" {
		defer writer.Close()
	}

	// ---------------------------
	// Step 8: Write output data
	// ---------------------------
	if err := toHandler.WriterFn(writer, opts.OutputFile, data); err != nil {
		return fmt.Errorf("failed to write output '%s': %w", opts.OutputFile, err)
	}

	// ---------------------------
	// Step 9: Success message
	// ---------------------------
	fmt.Printf("Successfully converted %s (%s) -> %s (%s)\n",
		opts.InputFile, opts.From, opts.OutputFile, opts.To)

	return nil
}

// readCloserWrapper wraps a Reader (e.g. gzip) and an underlying Closer (e.g. file)
type readCloserWrapper struct {
	io.Reader
	Closer io.Closer
}

func (w *readCloserWrapper) Close() error {
	// First convert Reader to Closer if possible (e.g. gzip reader should be closed?)
	// gzip.Reader.Close is actually important if it has checksums, but it implements ReadCloser since go1.10?
	// Actually gzip.Reader.Close just closes the underlying reader? No, it validates implementation.
	if rc, ok := w.Reader.(io.Closer); ok {
		if err := rc.Close(); err != nil {
			w.Closer.Close()
			return err
		}
	}
	return w.Closer.Close()
}

// writeCloserWrapper wraps a Writer (e.g. gzip) and an underlying Closer (e.g. file)
type writeCloserWrapper struct {
	io.Writer
	Closer io.Closer
}

func (w *writeCloserWrapper) Close() error {
	// Must close gzip writer to flush footer
	if wc, ok := w.Writer.(io.Closer); ok {
		if err := wc.Close(); err != nil {
			w.Closer.Close()
			return err
		}
	}
	return w.Closer.Close()
}
