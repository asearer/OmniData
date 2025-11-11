package formats

import (
	"fmt"
	"os"

	"omnidata/internal/convert"
)

// init registers the Avro format handler in the global Registry
// Note: Avro support requires additional dependencies
func init() {
	convert.RegisterFormat("avro", convert.FormatHandler{
		Name:     "avro",
		ReaderFn: readAvro,
		WriterFn: writeAvro,
	})
}

// readAvro reads Avro data from the given path.
// Returns [][]string containing all rows (similar to CSV).
func readAvro(path string) (interface{}, error) {
	if path == "" {
		return nil, fmt.Errorf("Avro read from STDIN is not supported")
	}

	// Check if file exists
	if _, err := os.Stat(path); err != nil {
		return nil, fmt.Errorf("failed to access Avro file '%s': %w", path, err)
	}

	// TODO: Implement actual Avro reading
	// This requires an Avro library like github.com/linkedin/goavro
	// For now, return an error indicating the feature needs implementation
	return nil, fmt.Errorf("Avro format support requires additional dependencies. " +
		"Install with: go get github.com/linkedin/goavro")
}

// writeAvro writes data to an Avro file at the given path.
// Expects data as [][]string (rows and columns).
func writeAvro(path string, data interface{}) error {
	if path == "" {
		return fmt.Errorf("Avro write to STDOUT is not supported")
	}

	if _, ok := data.([][]string); !ok {
		return fmt.Errorf("invalid data type for Avro writer, expected [][]string")
	}

	// TODO: Implement actual Avro writing
	// This requires an Avro library
	return fmt.Errorf("Avro format support requires additional dependencies. " +
		"Install with: go get github.com/linkedin/goavro")
}
