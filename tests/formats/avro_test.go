package formats_test

import (
	"os"
	"strings"
	"testing"

	"omnidata/internal/convert"
	_ "omnidata/internal/formats"
)

func TestAvroReader_Errors(t *testing.T) {
	handler, ok := convert.GetFormat("avro")
	if !ok {
		t.Fatal("Avro handler not registered")
	}
	// Test STDIN not supported
	_, err := handler.ReaderFn("")
	if err == nil || err.Error() == "" {
		t.Error("Expected error for Avro read from STDIN")
	}
	// Test file not found
	_, err = handler.ReaderFn("/tmp/does-not-exist.avro")
	if err == nil || err.Error() == "" {
		t.Error("Expected error for Avro file not found")
	}
	// Test dependency error
	f, _ := os.CreateTemp(os.TempDir(), "tmpavro.avro")
	defer os.Remove(f.Name())
	_, err = handler.ReaderFn(f.Name())
	if err == nil || err.Error() == "" || !strings.Contains(err.Error(), "Avro format support") {
		t.Error("Expected Avro dependency error")
	}
}

func TestAvroWriter_Errors(t *testing.T) {
	handler, ok := convert.GetFormat("avro")
	if !ok {
		t.Fatal("Avro handler not registered")
	}
	// STDOUT not supported
	err := handler.WriterFn("", [][]string{})
	if err == nil || err.Error() == "" {
		t.Error("Expected error for Avro write to STDOUT")
	}
	// Wrong type
	err = handler.WriterFn("foo.avro", 123)
	if err == nil || err.Error() == "" {
		t.Error("Expected error for Avro write with wrong type")
	}
	// Correct type, dependency error
	err = handler.WriterFn("foo.avro", [][]string{{"a"}})
	if err == nil || err.Error() == "" || !strings.Contains(err.Error(), "Avro format support") {
		t.Error("Expected Avro dependency error writing file")
	}
}
