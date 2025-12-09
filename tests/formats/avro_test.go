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
	// Test nil reader
	_, err := handler.ReaderFn(nil, "")
	if err == nil || err.Error() == "" {
		t.Error("Expected error for nil reader")
	}

	// Test dependency error
	f, _ := os.CreateTemp(os.TempDir(), "tmpavro.avro")
	defer os.Remove(f.Name())
	defer f.Close()
	_, err = handler.ReaderFn(f, f.Name())
	if err == nil || err.Error() == "" || !strings.Contains(err.Error(), "Avro format support") {
		t.Error("Expected Avro dependency error")
	}
}

func TestAvroWriter_Errors(t *testing.T) {
	handler, ok := convert.GetFormat("avro")
	if !ok {
		t.Fatal("Avro handler not registered")
	}
	// Nil writer
	err := handler.WriterFn(nil, "", [][]string{})
	if err == nil {
		t.Error("Expected error for nil writer")
	}
	// Wrong type
	err = handler.WriterFn(os.Stdout, "foo.avro", 123)
	if err == nil || err.Error() == "" {
		t.Error("Expected error for Avro write with wrong type")
	}
	// Correct type, dependency error
	err = handler.WriterFn(os.Stdout, "foo.avro", [][]string{{"a"}})
	if err == nil || err.Error() == "" || !strings.Contains(err.Error(), "Avro format support") {
		t.Error("Expected Avro dependency error writing file")
	}
}
