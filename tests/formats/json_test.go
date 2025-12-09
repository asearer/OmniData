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
	fOut, err := os.Create(outputJSON)
	if err != nil {
		t.Fatalf("failed to create output file: %v", err)
	}
	defer fOut.Close()

	if err := handler.WriterFn(fOut, outputJSON, data); err != nil {
		t.Fatalf("failed to write JSON: %v", err)
	}

	// Read JSON back
	fIn, err := os.Open(outputJSON)
	if err != nil {
		t.Fatalf("failed to open input file: %v", err)
	}
	defer fIn.Close()

	readData, err := handler.ReaderFn(fIn, outputJSON)
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

func TestJSON_InvalidCases(t *testing.T) {
	handler, ok := convert.GetFormat("json")
	if !ok {
		t.Fatal("JSON handler not registered")
	}
	// Invalid type for WriterFn
	if err := handler.WriterFn(os.Stdout, "foo.json", make(chan int)); err == nil {
		t.Error("expected error for invalid type, got nil")
	}
	// Nil reader
	_, err := handler.ReaderFn(nil, "foo.json")
	if err == nil {
		t.Error("expected error for nil reader, got nil")
	}

	// Bad/malformed JSON
	tmp, _ := os.CreateTemp(os.TempDir(), "bad.json")
	tmp.Write([]byte("not-json"))
	tmp.Seek(0, 0)
	defer tmp.Close()
	defer os.Remove(tmp.Name())

	_, err = handler.ReaderFn(tmp, tmp.Name())
	if err == nil {
		t.Error("expected error for malformed JSON, got nil")
	}
	// Empty file (should error, not panic)
	empty, _ := os.CreateTemp(os.TempDir(), "empty.json")
	defer empty.Close()
	defer os.Remove(empty.Name())

	_, err = handler.ReaderFn(empty, empty.Name())
	if err == nil {
		t.Error("expected error for empty file, got nil")
	}
}
