package formats

import (
	"fmt"
	"os"

	"github.com/xuri/excelize/v2"
	"omnidata/internal/convert"
)

func init() {
	convert.RegisterFormat("xlsx", convert.FormatHandler{
		Name:     "xlsx",
		ReaderFn: readXLSX,
		WriterFn: writeXLSX,
	})
}

func readXLSX(path string) (interface{}, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	sheets := f.GetSheetList()
	result := make(map[string][][]string)

	for _, sheet := range sheets {
		rows, err := f.GetRows(sheet)
		if err != nil {
			return nil, err
		}
		result[sheet] = rows
	}
	return result, nil
}

func writeXLSX(path string, data interface{}) error {
	x, ok := data.(map[string][][]string)
	if !ok {
		return fmt.Errorf("invalid data type for XLSX writer")
	}

	f := excelize.NewFile()
	defer f.Close()

	// Remove default sheet if not used
	if len(f.GetSheetList()) > 0 {
		f.DeleteSheet(f.GetSheetList()[0])
	}

	for sheet, rows := range x {
		index, err := f.NewSheet(sheet)
		if err != nil {
			return err
		}

		for rIdx, row := range rows {
			for cIdx, cell := range row {
				cellName, _ := excelize.CoordinatesToCellName(cIdx+1, rIdx+1)
				if err := f.SetCellValue(sheet, cellName, cell); err != nil {
					return err
				}
			}
		}
		// Set first sheet as active
		f.SetActiveSheet(index)
	}

	if err := f.SaveAs(path); err != nil {
		return err
	}
	return nil
}
