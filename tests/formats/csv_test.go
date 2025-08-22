package formats_test

import (
	"os"
	"path/filepath"
	"testing"

	"omnidata/internal/convert"
	_ "omnidata/internal/formats" // triggers init() for format registration
)

func TestCSVReadWrite(t *testing.T) {
	handler, ok := convert.GetFormat("csv")
	if !ok {
		t.Fatal("CSV handler not registered")
	}

	inputCSV := filepath.Join(os.TempDir(), "test.csv")
	defer os.Remove(inputCSV)
	content := []byte("a,b\n1,2\n3,4")
	if err := os.WriteFile(inputCSV, content, 0644); err != nil {
		t.Fatalf("failed to write input CSV: %v", err)
	}

	data, err := handler.ReaderFn(inputCSV)
	if err != nil {
		t.Fatalf("failed to read CSV: %v", err)
	}

	outputCSV := filepath.Join(os.TempDir(), "out.csv")
	defer os.Remove(outputCSV)
	if err := handler.WriterFn(outputCSV, data); err != nil {
		t.Fatalf("failed to write CSV: %v", err)
	}
}
