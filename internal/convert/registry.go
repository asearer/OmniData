package convert

import "strings"

// FormatHandler defines a format and its read/write functions
type FormatHandler struct {
	Name     string
	ReaderFn func(path string) (interface{}, error)
	WriterFn func(path string, data interface{}) error
}

// Registry holds all registered formats
var Registry = map[string]FormatHandler{}

// RegisterFormat registers a new format handler
func RegisterFormat(name string, handler FormatHandler) {
	Registry[strings.ToLower(name)] = handler
}

// ListFormats returns the names of all registered formats
func ListFormats() []string {
	formats := make([]string, 0, len(Registry))
	for name := range Registry {
		formats = append(formats, name)
	}
	return formats
}

// GetFormat retrieves a registered format handler by name (case-insensitive)
func GetFormat(name string) (FormatHandler, bool) {
	handler, ok := Registry[strings.ToLower(name)]
	return handler, ok
}
