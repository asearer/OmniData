package formats

import (
	"fmt"
	"io"

	"omnidata/internal/convert"

	"gopkg.in/yaml.v3"
)

// init registers the YAML format handler in the global Registry
func init() {
	convert.RegisterFormat("yaml", convert.FormatHandler{
		Name:     "yaml",
		ReaderFn: readYAML,
		WriterFn: writeYAML,
	})
}

// readYAML reads YAML data from the given reader.
func readYAML(r io.Reader, resource string) (interface{}, error) {
	if r == nil {
		return nil, fmt.Errorf("readYAML requires a valid reader")
	}

	// We'll decode into a generic interface{}.
	// Depending on structure it could be map[string]interface{} or []interface{}
	var data interface{}
	decoder := yaml.NewDecoder(r)
	if err := decoder.Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to decode YAML from '%s': %w", resource, err)
	}

	return data, nil
}

// writeYAML writes data as YAML to the given writer.
func writeYAML(w io.Writer, resource string, data interface{}) error {
	if w == nil {
		return fmt.Errorf("writeYAML requires a valid writer")
	}

	encoder := yaml.NewEncoder(w)
	// encoder.Close() is important for flushing any buffered data,
	// though for YAML it mostly closes the stream structure.
	defer encoder.Close()

	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("failed to encode YAML to '%s': %w", resource, err)
	}

	return nil
}
