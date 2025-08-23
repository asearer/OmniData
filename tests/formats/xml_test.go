package formats_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"omnidata/internal/convert"
	_ "omnidata/internal/formats" // triggers init() for format registration
)

// TestXMLReadWrite verifies that the XML format handler can correctly read and write XML files.
func TestXMLReadWrite(t *testing.T) {
	handler, ok := convert.GetFormat("xml")
	if !ok {
		t.Fatal("XML handler not registered")
	}

	// Sample XML content
	xmlContent := `<people><person><name>Alice</name><age>30</age></person></people>`

	// Create temporary input file
	inputXML := filepath.Join(os.TempDir(), "test.xml")
	defer os.Remove(inputXML)
	if err := os.WriteFile(inputXML, []byte(xmlContent), 0644); err != nil {
		t.Fatalf("failed to write input XML: %v", err)
	}

	// Read XML using the handler
	data, err := handler.ReaderFn(inputXML)
	if err != nil {
		t.Fatalf("failed to read XML: %v", err)
	}

	// Write XML back to a new file
	outputXML := filepath.Join(os.TempDir(), "out.xml")
	defer os.Remove(outputXML)
	if err := handler.WriterFn(outputXML, data); err != nil {
		t.Fatalf("failed to write XML: %v", err)
	}

	// Read written file and check for expected content
	written, err := os.ReadFile(outputXML)
	if err != nil {
		t.Fatalf("failed to read output XML: %v", err)
	}
	if !strings.Contains(string(written), "Alice") {
		t.Fatal("output XML does not contain expected content")
	}
}
