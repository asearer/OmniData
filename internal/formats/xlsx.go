package formats

import (
	"fmt"
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

// readXLSX reads an XLSX file from the given path.
// Returns a map of sheet names to [][]string representing rows and columns.
// If path is empty, this currently returns an error (STDIN not supported for XLSX).
func readXLSX(path string) (interface{}, error) {
	if path == "" {
		return nil, fmt.Errorf("XLSX read from STDIN is not supported")
	}

	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open XLSX file '%s': %w", path, err)
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

// writeXLSX writes data to an XLSX file at the given path.
// Expects data as map[string][][]string (sheet name -> rows).
// STDOUT writing is not supported because XLSX is binary.
func writeXLSX(path string, data interface{}) error {
	if path == "" {
		return fmt.Errorf("XLSX write to STDOUT is not supported")
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

	if err := f.SaveAs(path); err != nil {
		return fmt.Errorf("failed to save XLSX file '%s': %w", path, err)
	}
	return nil
}
