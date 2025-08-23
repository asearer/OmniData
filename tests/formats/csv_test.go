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
