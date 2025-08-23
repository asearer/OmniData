package formats_test

import (
	"os"
	"path/filepath"
	"testing"

	"omnidata/internal/convert"
	_ "omnidata/internal/formats" // triggers init() for format registration
)

// TestJSONReadWrite verifies that the JSON format handler can correctly write and read JSON files.
func TestJSONReadWrite(t *testing.T) {
	handler, ok := convert.GetFormat("json")
	if !ok {
		t.Fatal("JSON handler not registered")
	}

	// Test data: slice of maps
	data := []map[string]interface{}{
		{"name": "Alice", "age": 30},
		{"name": "Bob", "age": 25},
	}

	// Create temporary JSON output file
	outputJSON := filepath.Join(os.TempDir(), "test.json")
	defer os.Remove(outputJSON)

	// Write JSON using handler
	if err := handler.WriterFn(outputJSON, data); err != nil {
		t.Fatalf("failed to write JSON: %v", err)
	}

	// Read JSON back
	readData, err := handler.ReaderFn(outputJSON)
	if err != nil {
		t.Fatalf("failed to read JSON: %v", err)
	}

	// Verify number of objects
	arr, ok := readData.([]interface{})
	if !ok {
		t.Fatal("read data is not a slice")
	}
	if len(arr) != 2 {
		t.Fatalf("unexpected number of objects in JSON: got %d, want 2", len(arr))
	}
}
