package formats

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"omnidata/internal/convert"
)

// init registers the CSV format handler in the global Registry
func init() {
	convert.RegisterFormat("csv", convert.FormatHandler{
		Name:     "csv",
		ReaderFn: readCSV,
		WriterFn: writeCSV,
	})
}

// readCSV reads CSV data from the given path.
// If path is empty, reads from os.Stdin (for "-").
// Returns [][]string containing all rows.
func readCSV(path string) (interface{}, error) {
	var f *os.File
	var err error

	if path == "" {
		// STDIN support
		f = os.Stdin
	} else {
		f, err = os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("failed to open CSV file '%s': %w", path, err)
		}
		defer f.Close()
	}

	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil && err != io.EOF {
		return nil, fmt.Errorf("failed to read CSV data from '%s': %w", path, err)
	}

	return records, nil
}

// writeCSV writes [][]string data to the given path.
// If path is empty, writes to os.Stdout (for "-").
func writeCSV(path string, data interface{}) error {
	records, ok := data.([][]string)
	if !ok {
		return fmt.Errorf("invalid data for CSV writer, expected [][]string")
	}

	var f *os.File
	var err error

	if path == "" {
		// STDOUT support
		f = os.Stdout
	} else {
		f, err = os.Create(path)
		if err != nil {
			return fmt.Errorf("failed to create CSV file '%s': %w", path, err)
		}
		defer f.Close()
	}

	w := csv.NewWriter(f)
	if err := w.WriteAll(records); err != nil {
		return fmt.Errorf("failed to write CSV data to '%s': %w", path, err)
	}
	w.Flush()

	if err := w.Error(); err != nil {
		return fmt.Errorf("CSV writer flush error for '%s': %w", path, err)
	}

	return nil
}
