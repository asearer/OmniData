package formats_test

import (
	"os"
	"strings"
	"testing"

	"omnidata/internal/convert"
	_ "omnidata/internal/formats" // triggers init() for format registration
)

// normalizeXML removes whitespace and newlines for comparison
func normalizeXML(s string) string {
	s = strings.ReplaceAll(s, "\n", "")
	s = strings.ReplaceAll(s, "\r", "")
	s = strings.ReplaceAll(s, " ", "")
	return s
}

// TestXMLReadWrite verifies that the XML format handler can correctly read and write XML files.
func TestXMLReadWrite(t *testing.T) {
	handler, ok := convert.GetFormat("xml")
	if !ok {
		t.Fatal("XML handler not registered")
	}

	// Sample XML content
	xmlContent := `<people><person><name>Alice</name><age>30</age></person></people>`

	// Create temporary input file
	inputFile, err := os.CreateTemp("", "test-*.xml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer func() {
		if cerr := inputFile.Close(); cerr != nil {
			t.Errorf("failed to close input file: %v", cerr)
		}
		if err := os.Remove(inputFile.Name()); err != nil {
			t.Errorf("failed to remove input file: %v", err)
		}
	}()
	if _, err := inputFile.Write([]byte(xmlContent)); err != nil {
		t.Fatalf("failed to write input XML: %v", err)
	}

	// Read XML using the handler
	if _, err := inputFile.Seek(0, 0); err != nil {
		t.Fatalf("failed to seek: %v", err)
	}
	data, err := handler.ReaderFn(inputFile, inputFile.Name())
	if err != nil {
		t.Fatalf("failed to read XML: %v", err)
	}

	// Create temporary output file
	outputFile, err := os.CreateTemp("", "out-*.xml")
	if err != nil {
		t.Fatalf("failed to create output file: %v", err)
	}
	defer func() {
		if cerr := outputFile.Close(); cerr != nil {
			t.Errorf("failed to close output file: %v", cerr)
		}
		if err := os.Remove(outputFile.Name()); err != nil {
			t.Errorf("failed to remove output file: %v", err)
		}
	}()

	// Write XML back to output file
	if err := handler.WriterFn(outputFile, outputFile.Name(), data); err != nil {
		t.Fatalf("failed to write XML: %v", err)
	}

	// Read written file and check for expected content
	written, err := os.ReadFile(outputFile.Name())
	if err != nil {
		t.Fatalf("failed to read output XML: %v", err)
	}
	if !strings.Contains(normalizeXML(string(written)), "Alice") {
		t.Fatal("output XML does not contain expected content")
	}
}

// TestXML_InvalidCases ensures XML handler correctly errors on invalid inputs.
func TestXML_InvalidCases(t *testing.T) {
	handler, ok := convert.GetFormat("xml")
	if !ok {
		t.Fatal("XML handler not registered")
	}

	// Invalid type for WriterFn
	if err := handler.WriterFn(os.Stdout, "foo.xml", make(chan int)); err == nil {
		t.Error("expected error for invalid type, got nil")
	}

	// Nil reader
	_, err := handler.ReaderFn(nil, "foo.xml")
	if err == nil {
		t.Error("expected error for nil reader, got nil")
	}

	// Bad/malformed XML
	tmp, err := os.CreateTemp("", "bad-*.xml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	if _, err := tmp.Write([]byte("<not>")); err != nil {
		t.Fatalf("failed to write malformed XML: %v", err)
	}
	if cerr := tmp.Close(); cerr != nil {
		t.Errorf("failed to close tmp file: %v", cerr)
	}
	defer func() {
		if err := os.Remove(tmp.Name()); err != nil {
			t.Errorf("failed to remove tmp file: %v", err)
		}
	}()
	fBad, _ := os.Open(tmp.Name())
	defer fBad.Close()
	_, err = handler.ReaderFn(fBad, tmp.Name())
	if err == nil {
		t.Error("expected error for malformed XML, got nil")
	}

	// Empty file (should error, not panic)
	empty, err := os.CreateTemp("", "empty-*.xml")
	if err != nil {
		t.Fatalf("failed to create empty temp file: %v", err)
	}
	if cerr := empty.Close(); cerr != nil {
		t.Errorf("failed to close empty file: %v", cerr)
	}
	defer func() {
		if err := os.Remove(empty.Name()); err != nil {
			t.Errorf("failed to remove empty file: %v", err)
		}
	}()
	fEmpty, _ := os.Open(empty.Name())
	defer fEmpty.Close()
	_, err = handler.ReaderFn(fEmpty, empty.Name())
	if err == nil {
		t.Error("expected error for empty file, got nil")
	}
}
