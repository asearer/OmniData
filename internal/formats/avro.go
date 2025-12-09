package formats

import (
	"fmt"
	"io"

	"omnidata/internal/convert"
)

// init registers the Avro format handler in the global Registry
// Note: Avro support requires additional dependencies
func init() {
	convert.RegisterFormat("avro", convert.FormatHandler{
		Name:     "avro",
		ReaderFn: readAvro,
		WriterFn: writeAvro,
	})
}

// readAvro reads Avro data from the given reader.
func readAvro(r io.Reader, resource string) (interface{}, error) {
	if r == nil {
		return nil, fmt.Errorf("readAvro requires a valid reader")
	}

	// TODO: Implement actual Avro reading
	return nil, fmt.Errorf("Avro format support requires additional dependencies. " +
		"Install with: go get github.com/linkedin/goavro")
}

// writeAvro writes data to an Avro file to the given writer.
func writeAvro(w io.Writer, resource string, data interface{}) error {
	if w == nil {
		return fmt.Errorf("writeAvro requires a valid writer")
	}

	if _, ok := data.([][]string); !ok {
		return fmt.Errorf("invalid data type for Avro writer, expected [][]string")
	}

	// TODO: Implement actual Avro writing
	return fmt.Errorf("Avro format support requires additional dependencies. " +
		"Install with: go get github.com/linkedin/goavro")
}
