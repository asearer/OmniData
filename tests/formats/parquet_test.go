package formats_test

import (
	"os"
	"strings"
	"testing"

	"omnidata/internal/convert"
	_ "omnidata/internal/formats"
)

func TestParquetReader_Errors(t *testing.T) {
	handler, ok := convert.GetFormat("parquet")
	if !ok {
		t.Fatal("Parquet handler not registered")
	}
	// Test nil reader
	_, err := handler.ReaderFn(nil, "")
	if err == nil || err.Error() == "" {
		t.Error("Expected error for nil reader")
	}

	// Test dependency error
	f, _ := os.CreateTemp(os.TempDir(), "tmpparquet.parquet")
	defer os.Remove(f.Name())
	defer f.Close()
	_, err = handler.ReaderFn(f, f.Name())
	if err == nil || err.Error() == "" || !strings.Contains(err.Error(), "Parquet format support") {
		t.Error("Expected Parquet dependency error")
	}
}

func TestParquetWriter_Errors(t *testing.T) {
	handler, ok := convert.GetFormat("parquet")
	if !ok {
		t.Fatal("Parquet handler not registered")
	}
	// Nil writer
	err := handler.WriterFn(nil, "", [][]string{})
	if err == nil {
		t.Error("Expected error for nil writer")
	}
	// Wrong type
	err = handler.WriterFn(os.Stdout, "foo.parquet", 123)
	if err == nil || err.Error() == "" {
		t.Error("Expected error for Parquet write with wrong type")
	}
	// Correct type, dependency error
	err = handler.WriterFn(os.Stdout, "foo.parquet", [][]string{{"a"}})
	if err == nil || err.Error() == "" || !strings.Contains(err.Error(), "Parquet format support") {
		t.Error("Expected Parquet dependency error writing file")
	}
}
