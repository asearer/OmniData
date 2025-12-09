package formats

import (
	"encoding/xml"
	"fmt"
	"io"

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

// readXML reads XML data from the given reader.
func readXML(r io.Reader, resource string) (interface{}, error) {
	if r == nil {
		return nil, fmt.Errorf("readXML requires a valid reader")
	}

	// Try to decode as a generic Node to preserve structure
	var node Node
	if err := xml.NewDecoder(r).Decode(&node); err != nil {
		return nil, fmt.Errorf("failed to decode XML from '%s': %w", resource, err)
	}

	return node, nil
}

// writeXML writes data back to XML.
func writeXML(w io.Writer, resource string, data interface{}) error {
	if w == nil {
		return fmt.Errorf("writeXML requires a valid writer")
	}

	if _, ok := data.(Node); !ok {
		return fmt.Errorf("data is not a valid XML Node")
	}

	enc := xml.NewEncoder(w)
	enc.Indent("", "  ")
	if err := enc.Encode(data); err != nil {
		return fmt.Errorf("failed to encode XML to '%s': %w", resource, err)
	}
	return nil
}
