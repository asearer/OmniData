package convert_test

import (
	"os"
	"path/filepath"
	"testing"

	"omnidata/internal/convert"
	_ "omnidata/internal/formats" // triggers init() for format registration
)

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

func TestRunDryRun(t *testing.T) {
	input := tempFile(t, []byte("name,age\nAlice,30"))
	defer os.Remove(input)

	opts := convert.Options{
		InputFile:  input,
		OutputFile: "out.json",
		From:       "csv",
		To:         "json",
		DryRun:     true,
	}

	if err := convert.Run(opts); err != nil {
		t.Fatalf("dry-run failed: %v", err)
	}
}

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

	data, err := os.ReadFile(outputJSON)
	if err != nil {
		t.Fatalf("failed to read output: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("output JSON is empty")
	}
}
