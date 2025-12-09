package convert

import (
	"io"
	"strings"
)

/*
FormatHandler defines a data format and its associated read/write functions.

Responsibilities:
- Name: the canonical format name (e.g., "csv", "json", "xml", "xlsx").
- ReaderFn: function to read data from a given file path.
- WriterFn: function to write data to a given file path.
*/
type FormatHandler struct {
	Name     string
	ReaderFn func(r io.Reader, resource string) (interface{}, error)
	WriterFn func(w io.Writer, resource string, data interface{}) error
}

/*
Registry holds all registered format handlers.

Key points:
- Maps lowercase format names to handlers.
- Supports dynamic registration of new formats at runtime.
- Accessed via RegisterFormat, GetFormat, and ListFormats.
*/
var Registry = map[string]FormatHandler{}

/*
RegisterFormat registers a new format handler.

Automatically converts the format name to lowercase for case-insensitive retrieval.
Useful for plugin-style extensibility (e.g., adding XML, XLSX, YAML handlers).
*/
func RegisterFormat(name string, handler FormatHandler) {
	Registry[strings.ToLower(name)] = handler
}

/*
ListFormats returns the names of all registered formats.

Useful for displaying supported formats in CLI help or validation.
*/
func ListFormats() []string {
	formats := make([]string, 0, len(Registry))
	for name := range Registry {
		formats = append(formats, name)
	}
	return formats
}

/*
GetFormat retrieves a registered format handler by name.

- The lookup is case-insensitive.
- Returns the handler and a boolean indicating existence.
*/
func GetFormat(name string) (FormatHandler, bool) {
	handler, ok := Registry[strings.ToLower(name)]
	return handler, ok
}
