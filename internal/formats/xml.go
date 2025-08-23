package formats

import (
	"encoding/xml"
	"fmt"
	"os"

	"omnidata/internal/convert"
)

// init registers the XML format handler in the global Registry
func init() {
	convert.RegisterFormat("xml", convert.FormatHandler{
		Name:     "xml",
		ReaderFn: readXML,
		WriterFn: writeXML,
	})
}

// readXML reads XML data from the given path.
// If path is empty, reads from os.Stdin (for "-").
// Returns parsed data as interface{}.
func readXML(path string) (interface{}, error) {
	var f *os.File
	var err error

	if path == "" {
		// STDIN support
		f = os.Stdin
	} else {
		f, err = os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("failed to open XML file '%s': %w", path, err)
		}
		defer f.Close()
	}

	var data interface{}
	if err := xml.NewDecoder(f).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode XML from '%s': %w", path, err)
	}

	return data, nil
}

// writeXML writes data as XML to the given path.
// If path is empty, writes to os.Stdout (for "-").
func writeXML(path string, data interface{}) error {
	var f *os.File
	var err error

	if path == "" {
		// STDOUT support
		f = os.Stdout
	} else {
		f, err = os.Create(path)
		if err != nil {
			return fmt.Errorf("failed to create XML file '%s': %w", path, err)
		}
		defer f.Close()
	}

	enc := xml.NewEncoder(f)
	enc.Indent("", "  ")
	if err := enc.Encode(data); err != nil {
		return fmt.Errorf("failed to encode XML to '%s': %w", path, err)
	}

	return nil
}
