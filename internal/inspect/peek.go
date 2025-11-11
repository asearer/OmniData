package inspect

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"omnidata/internal/convert"
)

// PeekOptions holds configuration for the peek command
type PeekOptions struct {
	InputFile string
	Format    string
	Rows      int
	ShowStats bool
}

// PeekResult holds the result of peeking at data
type PeekResult struct {
	Schema  *Schema
	Preview []map[string]string // First N rows as key-value pairs
}

// RunPeek executes the peek command
func RunPeek(opts PeekOptions) error {
	// Get format handler
	handler, ok := convert.GetFormat(opts.Format)
	if !ok {
		return fmt.Errorf("unsupported format: %s", opts.Format)
	}

	// Resolve input path
	inputPath := opts.InputFile
	if inputPath == "-" {
		inputPath = ""
	} else {
		if _, err := os.Stat(inputPath); err != nil {
			return fmt.Errorf("input file does not exist: %s", inputPath)
		}
	}

	// Read data
	data, err := handler.ReaderFn(inputPath)
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}

	// Infer schema
	schema, err := InferSchema(data, opts.Format)
	if err != nil {
		return fmt.Errorf("failed to infer schema: %w", err)
	}

	// Get preview rows
	preview := getPreview(data, opts.Format, opts.Rows)

	// Display results
	displayPeek(schema, preview, opts)

	return nil
}

func getPreview(data interface{}, format string, maxRows int) []map[string]string {
	preview := make([]map[string]string, 0)

	switch format {
	case "csv":
		if records, ok := data.([][]string); ok && len(records) > 0 {
			headers := records[0]
			for i := 1; i < len(records) && i <= maxRows+1; i++ {
				row := make(map[string]string)
				for j, header := range headers {
					if j < len(records[i]) {
						row[header] = records[i][j]
					} else {
						row[header] = ""
					}
				}
				preview = append(preview, row)
			}
		}
	case "json":
		if arr, ok := data.([]interface{}); ok {
			for i := 0; i < len(arr) && i < maxRows; i++ {
				if obj, ok := arr[i].(map[string]interface{}); ok {
					row := make(map[string]string)
					for k, v := range obj {
						row[k] = fmt.Sprintf("%v", v)
					}
					preview = append(preview, row)
				}
			}
		} else if obj, ok := data.(map[string]interface{}); ok {
			row := make(map[string]string)
			for k, v := range obj {
				row[k] = fmt.Sprintf("%v", v)
			}
			preview = append(preview, row)
		}
	case "xlsx":
		if sheets, ok := data.(map[string][][]string); ok {
			for _, rows := range sheets {
				if len(rows) > 0 {
					headers := rows[0]
					for i := 1; i < len(rows) && i <= maxRows+1; i++ {
						row := make(map[string]string)
						for j, header := range headers {
							if j < len(rows[i]) {
								row[header] = rows[i][j]
							} else {
								row[header] = ""
							}
						}
						preview = append(preview, row)
					}
				}
				break // Only show first sheet
			}
		}
	}

	return preview
}

func displayPeek(schema *Schema, preview []map[string]string, opts PeekOptions) {
	fmt.Printf("\n📊 Schema Information\n")
	fmt.Printf("═══════════════════════════════════════\n")
	fmt.Printf("Format:      %s\n", schema.Format)
	fmt.Printf("Rows:        %d\n", schema.RowCount)
	fmt.Printf("Columns:     %d\n", schema.ColumnCount)
	fmt.Printf("\n")

	if opts.ShowStats {
		fmt.Printf("📈 Column Statistics\n")
		fmt.Printf("═══════════════════════════════════════\n")
		for _, col := range schema.Columns {
			fmt.Printf("  %s\n", col.Name)
			fmt.Printf("    Type:       %s\n", col.Type)
			fmt.Printf("    Nullable:   %v\n", col.Nullable)
			if col.MinLength >= 0 {
				fmt.Printf("    Length:     %d - %d\n", col.MinLength, col.MaxLength)
			}
			if len(col.SampleValues) > 0 {
				fmt.Printf("    Samples:    %s\n", strings.Join(col.SampleValues, ", "))
			}
			fmt.Printf("\n")
		}
	} else {
		fmt.Printf("📋 Columns\n")
		fmt.Printf("═══════════════════════════════════════\n")
		for i, col := range schema.Columns {
			nullable := ""
			if col.Nullable {
				nullable = " (nullable)"
			}
			fmt.Printf("  %d. %s: %s%s\n", i+1, col.Name, col.Type, nullable)
		}
		fmt.Printf("\n")
	}

	if len(preview) > 0 {
		fmt.Printf("👀 Preview (first %d rows)\n", len(preview))
		fmt.Printf("═══════════════════════════════════════\n")

		// Get all column names from first row
		allCols := make([]string, 0)
		if len(preview) > 0 {
			for k := range preview[0] {
				allCols = append(allCols, k)
			}
		}

		// Print header
		fmt.Printf("| ")
		for _, col := range allCols {
			val := col
			if len(val) > 20 {
				val = val[:17] + "..."
			}
			fmt.Printf("%-20s | ", val)
		}
		fmt.Printf("\n")
		fmt.Printf("| ")
		for range allCols {
			fmt.Printf("%-20s | ", strings.Repeat("-", 20))
		}
		fmt.Printf("\n")

		// Print rows
		for _, row := range preview {
			fmt.Printf("| ")
			for _, col := range allCols {
				val := row[col]
				if len(val) > 20 {
					val = val[:17] + "..."
				}
				fmt.Printf("%-20s | ", val)
			}
			fmt.Printf("\n")
		}
		fmt.Printf("\n")
	}
}

// FormatValue formats a value for display
func FormatValue(v interface{}) string {
	if v == nil {
		return "<null>"
	}

	switch val := v.(type) {
	case string:
		return val
	case int, int8, int16, int32, int64:
		return strconv.FormatInt(reflect.ValueOf(val).Int(), 10)
	case uint, uint8, uint16, uint32, uint64:
		return strconv.FormatUint(reflect.ValueOf(val).Uint(), 10)
	case float32, float64:
		return strconv.FormatFloat(reflect.ValueOf(val).Float(), 'f', -1, 64)
	case bool:
		return strconv.FormatBool(val)
	default:
		return fmt.Sprintf("%v", v)
	}
}
