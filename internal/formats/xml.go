package formats

import (
	"encoding/xml"
	"fmt"
	"os"

	"omnidata/internal/convert"
)

func init() {
	convert.RegisterFormat("xml", convert.FormatHandler{
		Name:     "xml",
		ReaderFn: readXML,
		WriterFn: writeXML,
	})
}

func readXML(path string) (interface{}, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var data interface{}
	if err := xml.NewDecoder(f).Decode(&data); err != nil {
		return nil, err
	}
	return data, nil
}

func writeXML(path string, data interface{}) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := xml.NewEncoder(f)
	enc.Indent("", "  ")
	if err := enc.Encode(data); err != nil {
		return fmt.Errorf("failed to encode XML: %w", err)
	}
	return nil
}
