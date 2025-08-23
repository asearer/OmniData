package convert_test

import (
	"testing"

	"omnidata/internal/convert"
	_ "omnidata/internal/formats" // triggers init() to register all formats
)

// TestValidateFormats verifies that ValidateFormats correctly accepts supported formats
// and rejects unsupported formats or when source and target are the same.
func TestValidateFormats(t *testing.T) {
	// Valid conversion: CSV -> JSON
	if err := convert.ValidateFormats("csv", "json"); err != nil {
		t.Fatalf("valid formats flagged as invalid: %v", err)
	}

	// Invalid: source and target are the same
	if err := convert.ValidateFormats("csv", "csv"); err == nil {
		t.Fatal("expected error when source and target formats are the same")
	}

	// Invalid: unsupported target format
	if err := convert.ValidateFormats("csv", "unknown"); err == nil {
		t.Fatal("expected error for unsupported target format")
	}

	// Invalid: unsupported source format
	if err := convert.ValidateFormats("unknown", "json"); err == nil {
		t.Fatal("expected error for unsupported source format")
	}
}
