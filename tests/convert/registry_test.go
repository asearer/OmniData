package convert_test

import (
	"testing"

	"omnidata/internal/convert"
	_ "omnidata/internal/formats" // triggers init() for format registration
)

func TestRegistry(t *testing.T) {
	handlers := convert.ListFormats()
	expected := map[string]bool{"csv": true, "json": true, "xml": true, "xlsx": true}

	for _, h := range handlers {
		if !expected[h] {
			t.Errorf("unexpected format registered: %s", h)
		}
		delete(expected, h)
	}

	if len(expected) != 0 {
		t.Errorf("expected formats missing: %v", expected)
	}
}

func TestGetFormat(t *testing.T) {
	handler, ok := convert.GetFormat("csv")
	if !ok || handler.Name != "csv" {
		t.Fatal("failed to get CSV handler")
	}

	handler, ok = convert.GetFormat("JSON")
	if !ok || handler.Name != "json" {
		t.Fatal("GetFormat should be case-insensitive")
	}
}
