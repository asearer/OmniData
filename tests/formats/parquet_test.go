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
	// Test STDIN not supported
	_, err := handler.ReaderFn("")
	if err == nil || err.Error() == "" {
		t.Error("Expected error for Parquet read from STDIN")
	}
	// Test file not found
	_, err = handler.ReaderFn("/tmp/does-not-exist.parquet")
	if err == nil || err.Error() == "" {
		t.Error("Expected error for Parquet file not found")
	}
	// Test dependency error
	f, _ := os.CreateTemp(os.TempDir(), "tmpparquet.parquet")
	defer os.Remove(f.Name())
	_, err = handler.ReaderFn(f.Name())
	if err == nil || err.Error() == "" || !strings.Contains(err.Error(), "Parquet format support") {
		t.Error("Expected Parquet dependency error")
	}
}

func TestParquetWriter_Errors(t *testing.T) {
	handler, ok := convert.GetFormat("parquet")
	if !ok {
		t.Fatal("Parquet handler not registered")
	}
	// STDOUT not supported
	err := handler.WriterFn("", [][]string{})
	if err == nil || err.Error() == "" {
		t.Error("Expected error for Parquet write to STDOUT")
	}
	// Wrong type
	err = handler.WriterFn("foo.parquet", 123)
	if err == nil || err.Error() == "" {
		t.Error("Expected error for Parquet write with wrong type")
	}
	// Correct type, dependency error
	err = handler.WriterFn("foo.parquet", [][]string{{"a"}})
	if err == nil || err.Error() == "" || !strings.Contains(err.Error(), "Parquet format support") {
		t.Error("Expected Parquet dependency error writing file")
	}
}
