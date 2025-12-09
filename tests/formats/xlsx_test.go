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
	fOut, err := os.Create(outputXLSX)
	if err != nil {
		t.Fatalf("failed to create output XLSX: %v", err)
	}
	defer fOut.Close()

	if err := handler.WriterFn(fOut, outputXLSX, data); err != nil {
		t.Fatalf("failed to write XLSX: %v", err)
	}

	// Read XLSX back
	fIn, err := os.Open(outputXLSX)
	if err != nil {
		t.Fatalf("failed to open input XLSX: %v", err)
	}
	defer fIn.Close()

	readData, err := handler.ReaderFn(fIn, outputXLSX)
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

func TestXLSX_InvalidCases(t *testing.T) {
	handler, ok := convert.GetFormat("xlsx")
	if !ok {
		t.Fatal("XLSX handler not registered")
	}
	// Nil reader
	_, err := handler.ReaderFn(nil, "foo.xlsx")
	if err == nil {
		t.Error("expected error for nil reader, got nil")
	}

	// WriterFn wrong type
	if err := handler.WriterFn(os.Stdout, "foo.xlsx", 12345); err == nil {
		t.Error("expected error for WriterFn wrong type, got nil")
	}

	// Write empty data
	f, _ := os.CreateTemp(os.TempDir(), "empty-out.xlsx")
	defer f.Close()
	defer os.Remove(f.Name())
	if err := handler.WriterFn(f, f.Name(), map[string][][]string{}); err != nil {
		// Acceptable if fails gracefully, but not panic
	}
}
