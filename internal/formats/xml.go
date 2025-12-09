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

// Node represents a generic XML element to allow round-tripping arbitrary XML
type Node struct {
	XMLName xml.Name
	Attrs   []xml.Attr `xml:",any,attr"`
	Content []byte     `xml:",innerxml"`
	Nodes   []Node     `xml:",any"`
}

// readXML reads XML data from the given path.
func readXML(path string) (interface{}, error) {
	var f *os.File
	var err error

	if path == "" {
		f = os.Stdin
	} else {
		f, err = os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("failed to open XML file '%s': %w", path, err)
		}
		defer f.Close()
	}

	// Try to decode as a generic Node to preserve structure
	var node Node
	if err := xml.NewDecoder(f).Decode(&node); err != nil {
		return nil, fmt.Errorf("failed to decode XML from '%s': %w", path, err)
	}

	return node, nil
}

// writeXML writes data back to XML.
func writeXML(path string, data interface{}) error {
	var f *os.File
	var err error

	if path == "" {
		f = os.Stdout
	} else {
		f, err = os.Create(path)
		if err != nil {
			return fmt.Errorf("failed to create XML file '%s': %w", path, err)
		}
		defer f.Close()
	}

	if _, ok := data.(Node); !ok {
		return fmt.Errorf("data is not a valid XML Node")
	}

	enc := xml.NewEncoder(f)
	enc.Indent("", "  ")
	if err := enc.Encode(data); err != nil {
		return fmt.Errorf("failed to encode XML to '%s': %w", path, err)
	}
	return nil
}
