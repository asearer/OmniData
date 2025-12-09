package formats

import (
	"encoding/json"
	"fmt"
	"io"

	"omnidata/internal/convert"
)

// init registers the JSON format handler in the global Registry
func init() {
	convert.RegisterFormat("json", convert.FormatHandler{
		Name:     "json",
		ReaderFn: readJSON,
		WriterFn: writeJSON,
	})
}

// readJSON reads JSON data from the given reader.
func readJSON(r io.Reader, resource string) (interface{}, error) {
	if r == nil {
		return nil, fmt.Errorf("readJSON requires a valid reader")
	}

	var data interface{}
	if err := json.NewDecoder(r).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode JSON from '%s': %w", resource, err)
	}

	return data, nil
}

// writeJSON writes data to the given writer as pretty-printed JSON.
func writeJSON(w io.Writer, resource string, data interface{}) error {
	if w == nil {
		return fmt.Errorf("writeJSON requires a valid writer")
	}

	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	if err := enc.Encode(data); err != nil {
		return fmt.Errorf("failed to encode JSON to '%s': %w", resource, err)
	}

	return nil
}
