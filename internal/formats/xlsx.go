package formats

import (
	"fmt"
	"io"
	"omnidata/internal/convert"

	"github.com/xuri/excelize/v2"
)

// init registers the XLSX format handler in the global Registry
func init() {
	convert.RegisterFormat("xlsx", convert.FormatHandler{
		Name:     "xlsx",
		ReaderFn: readXLSX,
		WriterFn: writeXLSX,
	})
}

// readXLSX reads an XLSX file from the given reader.
// Returns a map of sheet names to [][]string representing rows and columns.
func readXLSX(r io.Reader, resource string) (interface{}, error) {
	if r == nil {
		return nil, fmt.Errorf("readXLSX requires a valid reader")
	}

	f, err := excelize.OpenReader(r)
	if err != nil {
		return nil, fmt.Errorf("failed to open XLSX from reader: %w", err)
	}
	defer f.Close()

	sheets := f.GetSheetList()
	result := make(map[string][][]string)

	for _, sheet := range sheets {
		rows, err := f.GetRows(sheet)
		if err != nil {
			return nil, fmt.Errorf("failed to read rows from sheet '%s': %w", sheet, err)
		}
		result[sheet] = rows
	}

	return result, nil
}

// writeXLSX writes data to an XLSX file to the given writer.
// Expects data as map[string][][]string (sheet name -> rows).
func writeXLSX(w io.Writer, resource string, data interface{}) error {
	if w == nil {
		return fmt.Errorf("writeXLSX requires a valid writer")
	}

	x, ok := data.(map[string][][]string)
	if !ok {
		return fmt.Errorf("invalid data type for XLSX writer, expected map[string][][]string")
	}

	f := excelize.NewFile()
	defer f.Close()

	// Remove default sheet if not used
	if len(f.GetSheetList()) > 0 {
		_ = f.DeleteSheet(f.GetSheetList()[0])
	}

	for sheet, rows := range x {
		index, err := f.NewSheet(sheet)
		if err != nil {
			return fmt.Errorf("failed to create sheet '%s': %w", sheet, err)
		}

		for rIdx, row := range rows {
			for cIdx, cell := range row {
				cellName, _ := excelize.CoordinatesToCellName(cIdx+1, rIdx+1)
				if err := f.SetCellValue(sheet, cellName, cell); err != nil {
					return fmt.Errorf("failed to set cell value at %s: %w", cellName, err)
				}
			}
		}
		f.SetActiveSheet(index)
	}

	// Write to the provided writer
	if err := f.Write(w); err != nil {
		return fmt.Errorf("failed to write XLSX to '%s': %w", resource, err)
	}
	return nil
}
