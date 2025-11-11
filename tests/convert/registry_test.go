package convert_test

import (
	"testing"

	"omnidata/internal/convert"
	_ "omnidata/internal/formats" // triggers init() to register all formats
)

// TestRegistry verifies that all expected formats are registered in the global Registry.
func TestRegistry(t *testing.T) {
	handlers := convert.ListFormats()

	// Expected formats
	expected := map[string]bool{"csv": true, "json": true, "xml": true, "xlsx": true}

	// Check that all expected formats are present
	for _, h := range handlers {
		delete(expected, h) // remove any found format from expected
	}

	if len(expected) != 0 {
		t.Errorf("expected formats missing: %v", expected)
	}
}

// TestGetFormat ensures GetFormat works correctly (case-insensitive lookup)
func TestGetFormat(t *testing.T) {
	// Lowercase lookup
	handler, ok := convert.GetFormat("csv")
	if !ok || handler.Name != "csv" {
		t.Fatal("failed to get CSV handler")
	}

	// Uppercase lookup should also succeed
	handler, ok = convert.GetFormat("JSON")
	if !ok || handler.Name != "json" {
		t.Fatal("GetFormat should be case-insensitive")
	}
}
