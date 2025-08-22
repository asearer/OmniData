package formats

import (
	"encoding/csv"
	"fmt"
	"os"

	"omnidata/internal/convert"
)

func init() {
	convert.RegisterFormat("csv", convert.FormatHandler{
		Name:     "csv",
		ReaderFn: readCSV,
		WriterFn: writeCSV,
	})
}

func readCSV(path string) (interface{}, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	return records, nil
}

func writeCSV(path string, data interface{}) error {
	records, ok := data.([][]string)
	if !ok {
		return fmt.Errorf("invalid data for CSV writer")
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	return w.WriteAll(records)
}
