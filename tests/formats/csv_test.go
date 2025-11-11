package formats_test

import (
	"os"
	"path/filepath"
	"testing"

	"omnidata/internal/convert"
	_ "omnidata/internal/formats" // triggers init() for format registration
)

// TestCSVReadWrite verifies that CSV format handler can read and write CSV files correctly.
func TestCSVReadWrite(t *testing.T) {
	handler, ok := convert.GetFormat("csv")
	if !ok {
		t.Fatal("CSV handler not registered")
	}

	// Create temporary CSV input file
	inputCSV := filepath.Join(os.TempDir(), "test.csv")
	defer os.Remove(inputCSV)
	content := []byte("a,b\n1,2\n3,4")
	if err := os.WriteFile(inputCSV, content, 0644); err != nil {
		t.Fatalf("failed to write input CSV: %v", err)
	}

	// Read CSV using handler
	data, err := handler.ReaderFn(inputCSV)
	if err != nil {
		t.Fatalf("failed to read CSV: %v", err)
	}

	// Write CSV back to a temporary output file
	outputCSV := filepath.Join(os.TempDir(), "out.csv")
	defer os.Remove(outputCSV)
	if err := handler.WriterFn(outputCSV, data); err != nil {
		t.Fatalf("failed to write CSV: %v", err)
	}

	// Optional: verify output file exists and is non-empty
	info, err := os.Stat(outputCSV)
	if err != nil {
		t.Fatalf("output CSV file missing: %v", err)
	}
	if info.Size() == 0 {
		t.Fatal("output CSV file is empty")
	}
}

func TestCSV_InvalidCases(t *testing.T) {
	handler, ok := convert.GetFormat("csv")
	if !ok {
		t.Fatal("CSV handler not registered")
	}
	// Invalid type for WriterFn
	if err := handler.WriterFn("foo.csv", 12345); err == nil {
		t.Error("expected error for invalid type, got nil")
	}
	// Non-existent file on read
	_, err := handler.ReaderFn("/tmp/no-such-file.csv")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
	// Directory as input
	dir := os.TempDir()
	_, err = handler.ReaderFn(dir)
	if err == nil {
		t.Error("expected error for directory input, got nil")
	}
	// Empty file
	tmp, _ := os.CreateTemp(os.TempDir(), "empty.csv")
	tmp.Close()
	defer os.Remove(tmp.Name())
	_, err = handler.ReaderFn(tmp.Name())
	if err != nil {
		// Accept EOF, but no error (should not panic)
		t.Errorf("unexpected error for empty file: %v", err)
	}
}
