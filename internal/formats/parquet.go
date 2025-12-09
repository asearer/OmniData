package formats

import (
	"fmt"
	"io"

	"omnidata/internal/convert"
)

// init registers the Parquet format handler in the global Registry
// Note: Parquet support requires additional dependencies
func init() {
	convert.RegisterFormat("parquet", convert.FormatHandler{
		Name:     "parquet",
		ReaderFn: readParquet,
		WriterFn: writeParquet,
	})
}

// readParquet reads Parquet data from the given reader.
func readParquet(r io.Reader, resource string) (interface{}, error) {
	if r == nil {
		return nil, fmt.Errorf("readParquet requires a valid reader")
	}

	// TODO: Implement actual Parquet reading
	return nil, fmt.Errorf("Parquet format support requires additional dependencies. " +
		"Install with: go get github.com/xitongsys/parquet-go")
}

// writeParquet writes data to a Parquet file to the given writer.
func writeParquet(w io.Writer, resource string, data interface{}) error {
	if w == nil {
		return fmt.Errorf("writeParquet requires a valid writer")
	}

	if _, ok := data.([][]string); !ok {
		return fmt.Errorf("invalid data type for Parquet writer, expected [][]string")
	}

	// TODO: Implement actual Parquet writing
	return fmt.Errorf("Parquet format support requires additional dependencies. " +
		"Install with: go get github.com/xitongsys/parquet-go")
}
