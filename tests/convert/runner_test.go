package convert_test

import (
	"os"
	"path/filepath"
	"testing"

	"omnidata/internal/convert"
	_ "omnidata/internal/formats" // triggers init() for format registration
)

// tempFile creates a temporary file with the given content and returns its path.
func tempFile(t *testing.T, content []byte) string {
	f, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	if _, err := f.Write(content); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	f.Close()
	return f.Name()
}

// TestRunDryRun tests that the dry-run mode executes without writing any output.
func TestRunDryRun(t *testing.T) {
	input := tempFile(t, []byte("name,age\nAlice,30"))
	defer os.Remove(input)

	opts := convert.Options{
		InputFile:  input,
		OutputFile: "out.json", // should not be created
		From:       "csv",
		To:         "json",
		DryRun:     true,
	}

	if err := convert.Run(opts); err != nil {
		t.Fatalf("dry-run failed: %v", err)
	}

	// Ensure file was not created
	if _, err := os.Stat(opts.OutputFile); err == nil {
		t.Fatal("output file should not be created in dry-run mode")
	}
}

// TestRunCSVtoJSON tests a full CSV -> JSON conversion.
func TestRunCSVtoJSON(t *testing.T) {
	inputCSV := tempFile(t, []byte("name,age\nAlice,30"))
	defer os.Remove(inputCSV)

	outputJSON := filepath.Join(os.TempDir(), "test_out.json")
	defer os.Remove(outputJSON)

	opts := convert.Options{
		InputFile:  inputCSV,
		OutputFile: outputJSON,
		From:       "csv",
		To:         "json",
		DryRun:     false,
	}

	if err := convert.Run(opts); err != nil {
		t.Fatalf("conversion failed: %v", err)
	}

	// Verify output file exists and is not empty
	data, err := os.ReadFile(outputJSON)
	if err != nil {
		t.Fatalf("failed to read output: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("output JSON is empty")
	}
}

func TestRun_Errors(t *testing.T) {
	input := tempFile(t, []byte("name,age\nAlice,30"))
	defer os.Remove(input)

	// File not found
	opts := convert.Options{
		InputFile:  "/tmp/nope.csv",
		OutputFile: "nope.json",
		From:       "csv",
		To:         "json",
	}
	err := convert.Run(opts)
	if err == nil {
		t.Error("expected error for missing input file, got nil")
	}
	// Unsupported from format
	opts.InputFile = input
	opts.From = "NOTREAL"
	err = convert.Run(opts)
	if err == nil {
		t.Error("expected error for unknown input format, got nil")
	}
	// Unsupported to format
	opts.From = "csv"
	opts.To = "UNK"
	err = convert.Run(opts)
	if err == nil {
		t.Error("expected error for unknown output format, got nil")
	}
	// Output file exists (should refuse to overwrite)
	f, _ := os.CreateTemp(os.TempDir(), "alreadythere.json")
	f.Close()
	defer os.Remove(f.Name())
	opts.To = "json"
	opts.OutputFile = f.Name()
	err = convert.Run(opts)
	if err == nil {
		t.Error("expected error for existing output file, got nil")
	}
	// STDIN/STDOUT for unsupported type (xlsx)
	opts.InputFile = ""
	opts.OutputFile = ""
	opts.From = "xlsx"
	opts.To = "csv"
	err = convert.Run(opts)
	if err == nil {
		t.Error("expected error for reading xlsx from STDIN, got nil")
	}
	opts.From = "csv"
	opts.To = "xlsx"
	err = convert.Run(opts)
	if err == nil {
		t.Error("expected error for writing xlsx to STDOUT, got nil")
	}
	// Both formats unknown
	opts.From = "?"
	opts.To = "?"
	err = convert.Run(opts)
	if err == nil {
		t.Error("expected error for both formats unknown, got nil")
	}
}
