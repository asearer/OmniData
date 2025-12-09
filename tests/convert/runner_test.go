package convert_test

import (
	"os"
	"path/filepath"
	"strings"
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

// TestRunGzip tests transparent Gzip handling.
func TestRunGzip(t *testing.T) {
	// Create a simple CSV file
	inputCSV := tempFile(t, []byte("name,age\nAlice,30"))
	defer os.Remove(inputCSV)

	// Output file with .gz extension
	outputGZ := filepath.Join(os.TempDir(), "test_out.csv.gz")
	defer os.Remove(outputGZ)

	opts := convert.Options{
		InputFile:  inputCSV,
		OutputFile: outputGZ,
		From:       "csv",
		To:         "csv",
	}

	if err := convert.Run(opts); err != nil {
		t.Fatalf("conversion to gzip failed: %v", err)
	}

	// Verify output exists and is not empty
	stat, err := os.Stat(outputGZ)
	if err != nil {
		t.Fatalf("output gz file missing: %v", err)
	}
	if stat.Size() == 0 {
		t.Fatal("output gz file is empty")
	}

	// Try reading it back from .gz to .csv
	// This verifies decompression works too
	outputCSVRec := filepath.Join(os.TempDir(), "test_rec.csv")
	defer os.Remove(outputCSVRec)

	opts2 := convert.Options{
		InputFile:  outputGZ,
		OutputFile: outputCSVRec,
		From:       "csv",
		To:         "csv",
	}

	if err := convert.Run(opts2); err != nil {
		t.Fatalf("conversion from gzip failed: %v", err)
	}

	// Verify recovered CSV matches original (roughly)
	data, err := os.ReadFile(outputCSVRec)
	if err != nil {
		t.Fatalf("failed to read recovered csv: %v", err)
	}
	content := string(data)
	if !strings.Contains(content, "Alice") || !strings.Contains(content, "30") {
		t.Errorf("recovered csv content mismatch. Got: %s", content)
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
