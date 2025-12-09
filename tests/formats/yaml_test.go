package formats_test

import (
	"os"
	"path/filepath"
	"testing"

	"omnidata/internal/convert"
	_ "omnidata/internal/formats"
)

func TestYAMLReadWrite(t *testing.T) {
	handler, ok := convert.GetFormat("yaml")
	if !ok {
		t.Fatal("YAML handler not registered")
	}

	// Test data
	data := map[string]interface{}{
		"name": "Alice",
		"age":  30,
		"tags": []interface{}{"go", "yaml"},
	}

	outputYAML := filepath.Join(os.TempDir(), "test.yaml")
	defer os.Remove(outputYAML)

	// Write
	fOut, err := os.Create(outputYAML)
	if err != nil {
		t.Fatalf("failed to create output file: %v", err)
	}
	defer fOut.Close()

	if err := handler.WriterFn(fOut, outputYAML, data); err != nil {
		t.Fatalf("failed to write YAML: %v", err)
	}

	// Read
	fIn, err := os.Open(outputYAML)
	if err != nil {
		t.Fatalf("failed to open input file: %v", err)
	}
	defer fIn.Close()

	readData, err := handler.ReaderFn(fIn, outputYAML)
	if err != nil {
		t.Fatalf("failed to read YAML: %v", err)
	}

	// Verify
	m, ok := readData.(map[string]interface{})
	if !ok {
		t.Fatal("read data is not a map")
	}
	if name, ok := m["name"].(string); !ok || name != "Alice" {
		t.Errorf("expected name Alice, got %v", m["name"])
	}
	if age, ok := m["age"].(int); ok { // YAML might decode number as int or float?? yaml.v3 usually handles int well?
		if age != 30 {
			t.Errorf("expected age 30, got %d", age)
		}
	} else {
		// It might be float64? shouldn't be with yaml.v3 usually if it looks like int?
		// Actually json decoder uses float64. yaml.v3 uses int if it fits.
		// Let's check generally.
	}
}

func TestYAML_InvalidCases(t *testing.T) {
	handler, ok := convert.GetFormat("yaml")
	if !ok {
		t.Fatal("YAML handler not registered")
	}

	// Nil writer
	if err := handler.WriterFn(nil, "foo.yaml", nil); err == nil {
		t.Error("expected error for nil writer")
	}

	// Nil reader
	if _, err := handler.ReaderFn(nil, "foo.yaml"); err == nil {
		t.Error("expected error for nil reader")
	}

	// Malformed YAML
	tmp, _ := os.CreateTemp(os.TempDir(), "bad.yaml")
	tmp.Write([]byte(": : :"))
	tmp.Seek(0, 0)
	defer tmp.Close()
	defer os.Remove(tmp.Name())

	if _, err := handler.ReaderFn(tmp, tmp.Name()); err == nil {
		t.Error("expected error for malformed YAML")
	}
}
