package inspect

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// ColumnInfo holds information about a single column
type ColumnInfo struct {
	Name         string
	Type         string
	Nullable     bool
	MinLength    int
	MaxLength    int
	SampleValues []string
}

// Schema represents the structure of a dataset
type Schema struct {
	Format      string
	RowCount    int
	ColumnCount int
	Columns     []ColumnInfo
}

// InferSchema analyzes data and returns schema information
func InferSchema(data interface{}, format string) (*Schema, error) {
	switch format {
	case "csv":
		return inferCSVSchema(data)
	case "json":
		return inferJSONSchema(data)
	case "xml":
		return inferXMLSchema(data)
	case "xlsx":
		return inferXLSXSchema(data)
	default:
		return nil, fmt.Errorf("unsupported format for schema inference: %s", format)
	}
}

func inferCSVSchema(data interface{}) (*Schema, error) {
	records, ok := data.([][]string)
	if !ok {
		return nil, fmt.Errorf("invalid CSV data type")
	}

	if len(records) == 0 {
		return &Schema{Format: "csv", RowCount: 0, ColumnCount: 0}, nil
	}

	headers := records[0]
	columnCount := len(headers)
	rowCount := len(records) - 1 // Exclude header

	columns := make([]ColumnInfo, columnCount)

	// Initialize columns
	for i, header := range headers {
		columns[i] = ColumnInfo{
			Name:         header,
			Type:         "string",
			Nullable:     false,
			MinLength:    -1,
			MaxLength:    -1,
			SampleValues: make([]string, 0, 5),
		}
	}

	// Analyze data rows (limit to first 1000 for performance)
	maxRows := 1000
	if rowCount > maxRows {
		maxRows = rowCount
	}

	for rowIdx := 1; rowIdx <= maxRows && rowIdx < len(records); rowIdx++ {
		row := records[rowIdx]
		for colIdx := 0; colIdx < columnCount && colIdx < len(row); colIdx++ {
			value := row[colIdx]

			// Update nullable
			if value == "" {
				columns[colIdx].Nullable = true
			}

			// Update length stats
			length := len(value)
			if columns[colIdx].MinLength == -1 || length < columns[colIdx].MinLength {
				columns[colIdx].MinLength = length
			}
			if length > columns[colIdx].MaxLength {
				columns[colIdx].MaxLength = length
			}

			// Try to infer type
			if columns[colIdx].Type == "string" {
				if value != "" {
					if _, err := strconv.ParseFloat(value, 64); err == nil {
						columns[colIdx].Type = "number"
					} else if strings.ToLower(value) == "true" || strings.ToLower(value) == "false" {
						columns[colIdx].Type = "boolean"
					}
				}
			}

			// Collect sample values
			if len(columns[colIdx].SampleValues) < 5 && value != "" {
				columns[colIdx].SampleValues = append(columns[colIdx].SampleValues, value)
			}
		}
	}

	return &Schema{
		Format:      "csv",
		RowCount:    rowCount,
		ColumnCount: columnCount,
		Columns:     columns,
	}, nil
}

func inferJSONSchema(data interface{}) (*Schema, error) {
	schema := &Schema{Format: "json"}

	switch v := data.(type) {
	case []interface{}:
		// Array of objects
		if len(v) == 0 {
			return schema, nil
		}

		// Analyze first object to get structure
		if firstObj, ok := v[0].(map[string]interface{}); ok {
			columns := make([]ColumnInfo, 0)
			for key, val := range firstObj {
				col := ColumnInfo{
					Name:         key,
					Type:         inferJSONType(val),
					Nullable:     val == nil,
					SampleValues: make([]string, 0, 5),
				}
				columns = append(columns, col)
			}

			schema.RowCount = len(v)
			schema.ColumnCount = len(columns)
			schema.Columns = columns
		} else {
			// Array of primitives
			schema.RowCount = len(v)
			schema.ColumnCount = 1
			schema.Columns = []ColumnInfo{
				{
					Name:         "value",
					Type:         inferJSONType(v[0]),
					Nullable:     false,
					SampleValues: make([]string, 0, 5),
				},
			}
		}
	case map[string]interface{}:
		// Single object
		columns := make([]ColumnInfo, 0)
		for key, val := range v {
			col := ColumnInfo{
				Name:         key,
				Type:         inferJSONType(val),
				Nullable:     val == nil,
				SampleValues: make([]string, 0, 5),
			}
			columns = append(columns, col)
		}
		schema.RowCount = 1
		schema.ColumnCount = len(columns)
		schema.Columns = columns
	default:
		return nil, fmt.Errorf("unsupported JSON structure")
	}

	return schema, nil
}

func inferJSONType(val interface{}) string {
	if val == nil {
		return "null"
	}

	switch reflect.TypeOf(val).Kind() {
	case reflect.String:
		return "string"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return "number"
	case reflect.Bool:
		return "boolean"
	case reflect.Slice, reflect.Array:
		return "array"
	case reflect.Map:
		return "object"
	default:
		return "unknown"
	}
}

func inferXMLSchema(data interface{}) (*Schema, error) {
	// XML schema inference is similar to JSON
	return inferJSONSchema(data)
}

func inferXLSXSchema(data interface{}) (*Schema, error) {
	sheets, ok := data.(map[string][][]string)
	if !ok {
		return nil, fmt.Errorf("invalid XLSX data type")
	}

	if len(sheets) == 0 {
		return &Schema{Format: "xlsx", RowCount: 0, ColumnCount: 0}, nil
	}

	// Use first sheet for schema
	var firstSheet string
	var firstData [][]string
	for sheet, rows := range sheets {
		firstSheet = sheet
		firstData = rows
		break
	}

	if len(firstData) == 0 {
		return &Schema{Format: "xlsx", RowCount: 0, ColumnCount: 0}, nil
	}

	// Treat XLSX similar to CSV
	records := firstData
	headers := records[0]
	columnCount := len(headers)
	rowCount := len(records) - 1

	columns := make([]ColumnInfo, columnCount)
	for i, header := range headers {
		columns[i] = ColumnInfo{
			Name:         fmt.Sprintf("%s.%s", firstSheet, header),
			Type:         "string",
			Nullable:     false,
			MinLength:    -1,
			MaxLength:    -1,
			SampleValues: make([]string, 0, 5),
		}
	}

	// Analyze data rows
	maxRows := 1000
	if rowCount > maxRows {
		maxRows = rowCount
	}

	for rowIdx := 1; rowIdx <= maxRows && rowIdx < len(records); rowIdx++ {
		row := records[rowIdx]
		for colIdx := 0; colIdx < columnCount && colIdx < len(row); colIdx++ {
			value := row[colIdx]

			if value == "" {
				columns[colIdx].Nullable = true
			}

			length := len(value)
			if columns[colIdx].MinLength == -1 || length < columns[colIdx].MinLength {
				columns[colIdx].MinLength = length
			}
			if length > columns[colIdx].MaxLength {
				columns[colIdx].MaxLength = length
			}

			if len(columns[colIdx].SampleValues) < 5 && value != "" {
				columns[colIdx].SampleValues = append(columns[colIdx].SampleValues, value)
			}
		}
	}

	return &Schema{
		Format:      "xlsx",
		RowCount:    rowCount,
		ColumnCount: columnCount,
		Columns:     columns,
	}, nil
}
