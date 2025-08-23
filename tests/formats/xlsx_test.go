package formats_test

import (
	"os"
	"path/filepath"
	"testing"

	"omnidata/internal/convert"
	_ "omnidata/internal/formats" // triggers init() for format registration
)

// TestXLSXReadWrite verifies that the XLSX format handler can correctly write and read Excel files.
func TestXLSXReadWrite(t *testing.T) {
	handler, ok := convert.GetFormat("xlsx")
	if !ok {
		t.Fatal("XLSX handler not registered")
	}

	// Sample XLSX data: map[sheetName][][]string
	data := map[string][][]string{
		"Sheet1": {
			{"Name", "Age"},
			{"Alice", "30"},
			{"Bob", "25"},
		},
	}

	outputXLSX := filepath.Join(os.TempDir(), "test.xlsx")
	defer os.Remove(outputXLSX)

	// Write XLSX using the handler
	if err := handler.WriterFn(outputXLSX, data); err != nil {
		t.Fatalf("failed to write XLSX: %v", err)
	}

	// Read XLSX back
	readData, err := handler.ReaderFn(outputXLSX)
	if err != nil {
		t.Fatalf("failed to read XLSX: %v", err)
	}

	// Type assertion
	readMap, ok := readData.(map[string][][]string)
	if !ok {
		t.Fatal("read XLSX data has incorrect type")
	}

	// Verify number of rows
	if len(readMap["Sheet1"]) != 3 {
		t.Fatalf("expected 3 rows in Sheet1, got %d", len(readMap["Sheet1"]))
	}

	// Verify specific cell value
	if readMap["Sheet1"][1][0] != "Alice" {
		t.Fatalf("unexpected value in row 2, col 1: %s", readMap["Sheet1"][1][0])
	}
}
