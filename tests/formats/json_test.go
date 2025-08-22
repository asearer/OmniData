package formats_test

import (
	"os"
	"path/filepath"
	"testing"

	"omnidata/internal/convert"
	_ "omnidata/internal/formats" // triggers init() for format registration
)

func TestJSONReadWrite(t *testing.T) {
	handler, ok := convert.GetFormat("json")
	if !ok {
		t.Fatal("JSON handler not registered")
	}

	data := []map[string]interface{}{
		{"name": "Alice", "age": 30},
		{"name": "Bob", "age": 25},
	}

	outputJSON := filepath.Join(os.TempDir(), "test.json")
	defer os.Remove(outputJSON)

	if err := handler.WriterFn(outputJSON, data); err != nil {
		t.Fatalf("failed to write JSON: %v", err)
	}

	readData, err := handler.ReaderFn(outputJSON)
	if err != nil {
		t.Fatalf("failed to read JSON: %v", err)
	}

	if len(readData.([]interface{})) != 2 {
		t.Fatal("unexpected number of objects in JSON")
	}
}
