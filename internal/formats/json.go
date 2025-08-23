package formats

import (
	"encoding/json"
	"fmt"
	"os"

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

// readJSON reads JSON data from the given path.
// If path is empty, reads from os.Stdin (for "-").
// Returns data as interface{} (map[string]interface{} or []interface{}).
func readJSON(path string) (interface{}, error) {
	var f *os.File
	var err error

	if path == "" {
		// STDIN support
		f = os.Stdin
	} else {
		f, err = os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("failed to open JSON file '%s': %w", path, err)
		}
		defer f.Close()
	}

	var data interface{}
	if err := json.NewDecoder(f).Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode JSON from '%s': %w", path, err)
	}

	return data, nil
}

// writeJSON writes data to the given path as pretty-printed JSON.
// If path is empty, writes to os.Stdout (for "-").
func writeJSON(path string, data interface{}) error {
	var f *os.File
	var err error

	if path == "" {
		// STDOUT support
		f = os.Stdout
	} else {
		f, err = os.Create(path)
		if err != nil {
			return fmt.Errorf("failed to create JSON file '%s': %w", path, err)
		}
		defer f.Close()
	}

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(data); err != nil {
		return fmt.Errorf("failed to encode JSON to '%s': %w", path, err)
	}

	return nil
}
