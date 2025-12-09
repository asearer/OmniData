package formats

import (
	"encoding/csv"
	"fmt"
	"io"

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

// readCSV reads CSV data from the given reader.
func readCSV(r io.Reader, resource string) (interface{}, error) {
	if r == nil {
		return nil, fmt.Errorf("readCSV requires a valid reader")
	}

	reader := csv.NewReader(r)

	// Read all records
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV from '%s': %w", resource, err)
	}

	return records, nil
}

// writeCSV writes data as CSV to the given writer.
func writeCSV(w io.Writer, resource string, data interface{}) error {
	if w == nil {
		return fmt.Errorf("writeCSV requires a valid writer")
	}

	records, ok := data.([][]string)
	if !ok {
		return fmt.Errorf("invalid data type for CSV writer, expected [][]string")
	}

	writer := csv.NewWriter(w)
	if err := writer.WriteAll(records); err != nil {
		return fmt.Errorf("failed to write CSV to '%s': %w", resource, err)
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return fmt.Errorf("error flushing CSV writer: %w", err)
	}

	return nil
}
