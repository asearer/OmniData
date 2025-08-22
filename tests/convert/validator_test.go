package convert_test

import (
	"testing"

	"omnidata/internal/convert"
	_ "omnidata/internal/formats" // triggers init() for format registration
)

func TestValidateFormats(t *testing.T) {
	if err := convert.ValidateFormats("csv", "json"); err != nil {
		t.Fatalf("valid formats flagged as invalid: %v", err)
	}

	if err := convert.ValidateFormats("csv", "csv"); err == nil {
		t.Fatal("expected error when source and target are the same")
	}

	if err := convert.ValidateFormats("csv", "unknown"); err == nil {
		t.Fatal("expected error for unsupported target")
	}

	if err := convert.ValidateFormats("unknown", "json"); err == nil {
		t.Fatal("expected error for unsupported source")
	}
}
