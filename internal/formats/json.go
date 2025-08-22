package formats

import (
	"encoding/json"
	"os"

	"omnidata/internal/convert"
)

func init() {
	convert.RegisterFormat("json", convert.FormatHandler{
		Name:     "json",
		ReaderFn: readJSON,
		WriterFn: writeJSON,
	})
}

func readJSON(path string) (interface{}, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var data interface{}
	if err := json.NewDecoder(f).Decode(&data); err != nil {
		return nil, err
	}
	return data, nil
}

func writeJSON(path string, data interface{}) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(data)
}
