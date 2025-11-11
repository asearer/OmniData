package output

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"strings"

	"omnidata/internal/inspect"
)

// Formatter defines an output formatter interface
type Formatter interface {
	FormatSchema(schema *inspect.Schema) (string, error)
	FormatDiff(diff *inspect.SchemaDiff, schema1, schema2 *inspect.Schema) (string, error)
}

// MarkdownFormatter formats output as Markdown
type MarkdownFormatter struct{}

func (f *MarkdownFormatter) FormatSchema(schema *inspect.Schema) (string, error) {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("# Schema: %s\n\n", schema.Format))
	sb.WriteString(fmt.Sprintf("- **Rows:** %d\n", schema.RowCount))
	sb.WriteString(fmt.Sprintf("- **Columns:** %d\n\n", schema.ColumnCount))

	sb.WriteString("## Columns\n\n")
	sb.WriteString("| Name | Type | Nullable |\n")
	sb.WriteString("|------|------|----------|\n")

	for _, col := range schema.Columns {
		nullable := "No"
		if col.Nullable {
			nullable = "Yes"
		}
		sb.WriteString(fmt.Sprintf("| %s | %s | %s |\n", col.Name, col.Type, nullable))
	}

	return sb.String(), nil
}

func (f *MarkdownFormatter) FormatDiff(diff *inspect.SchemaDiff, schema1, schema2 *inspect.Schema) (string, error) {
	var sb strings.Builder

	sb.WriteString("# Schema Comparison\n\n")
	sb.WriteString(fmt.Sprintf("## File 1: %s\n", schema1.Format))
	sb.WriteString(fmt.Sprintf("- Rows: %d\n", schema1.RowCount))
	sb.WriteString(fmt.Sprintf("- Columns: %d\n\n", schema1.ColumnCount))

	sb.WriteString(fmt.Sprintf("## File 2: %s\n", schema2.Format))
	sb.WriteString(fmt.Sprintf("- Rows: %d\n", schema2.RowCount))
	sb.WriteString(fmt.Sprintf("- Columns: %d\n\n", schema2.ColumnCount))

	if len(diff.AddedColumns) > 0 {
		sb.WriteString(fmt.Sprintf("## Added Columns (%d)\n\n", len(diff.AddedColumns)))
		for _, col := range diff.AddedColumns {
			nullable := ""
			if col.Nullable {
				nullable = " (nullable)"
			}
			sb.WriteString(fmt.Sprintf("- **%s**: %s%s\n", col.Name, col.Type, nullable))
		}
		sb.WriteString("\n")
	}

	if len(diff.RemovedColumns) > 0 {
		sb.WriteString(fmt.Sprintf("## Removed Columns (%d)\n\n", len(diff.RemovedColumns)))
		for _, col := range diff.RemovedColumns {
			nullable := ""
			if col.Nullable {
				nullable = " (nullable)"
			}
			sb.WriteString(fmt.Sprintf("- **%s**: %s%s\n", col.Name, col.Type, nullable))
		}
		sb.WriteString("\n")
	}

	if len(diff.ChangedColumns) > 0 {
		sb.WriteString(fmt.Sprintf("## Changed Columns (%d)\n\n", len(diff.ChangedColumns)))
		for _, change := range diff.ChangedColumns {
			sb.WriteString(fmt.Sprintf("### %s\n", change.Name))
			if change.OldType != change.NewType {
				sb.WriteString(fmt.Sprintf("- Type: `%s` → `%s`\n", change.OldType, change.NewType))
			}
			if change.OldNullable != change.NewNullable {
				oldVal := "not nullable"
				newVal := "nullable"
				if change.OldNullable {
					oldVal = "nullable"
					newVal = "not nullable"
				}
				sb.WriteString(fmt.Sprintf("- Nullable: %s → %s\n", oldVal, newVal))
			}
			sb.WriteString("\n")
		}
	}

	return sb.String(), nil
}

// HTMLFormatter formats output as HTML
type HTMLFormatter struct{}

func (f *HTMLFormatter) FormatSchema(schema *inspect.Schema) (string, error) {
	tmpl := `<!DOCTYPE html>
<html>
<head>
	<title>Schema: {{.Format}}</title>
	<style>
		body { font-family: Arial, sans-serif; margin: 20px; }
		table { border-collapse: collapse; width: 100%; }
		th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }
		th { background-color: #4CAF50; color: white; }
		tr:nth-child(even) { background-color: #f2f2f2; }
	</style>
</head>
<body>
	<h1>Schema: {{.Format}}</h1>
	<p><strong>Rows:</strong> {{.RowCount}}</p>
	<p><strong>Columns:</strong> {{.ColumnCount}}</p>
	
	<h2>Columns</h2>
	<table>
		<tr>
			<th>Name</th>
			<th>Type</th>
			<th>Nullable</th>
		</tr>
		{{range .Columns}}
		<tr>
			<td>{{.Name}}</td>
			<td>{{.Type}}</td>
			<td>{{if .Nullable}}Yes{{else}}No{{end}}</td>
		</tr>
		{{end}}
	</table>
</body>
</html>`

	t, err := template.New("schema").Parse(tmpl)
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	if err := t.Execute(&sb, schema); err != nil {
		return "", err
	}

	return sb.String(), nil
}

func (f *HTMLFormatter) FormatDiff(diff *inspect.SchemaDiff, schema1, schema2 *inspect.Schema) (string, error) {
	// Simple HTML diff output
	var sb strings.Builder
	sb.WriteString("<!DOCTYPE html>\n<html>\n<head>\n<title>Schema Comparison</title>\n")
	sb.WriteString("<style>body { font-family: Arial, sans-serif; margin: 20px; }</style>\n")
	sb.WriteString("</head>\n<body>\n")
	sb.WriteString("<h1>Schema Comparison</h1>\n")
	sb.WriteString(fmt.Sprintf("<h2>File 1: %s</h2>\n", schema1.Format))
	sb.WriteString(fmt.Sprintf("<p>Rows: %d, Columns: %d</p>\n", schema1.RowCount, schema1.ColumnCount))
	sb.WriteString(fmt.Sprintf("<h2>File 2: %s</h2>\n", schema2.Format))
	sb.WriteString(fmt.Sprintf("<p>Rows: %d, Columns: %d</p>\n", schema2.RowCount, schema2.ColumnCount))

	if len(diff.AddedColumns) > 0 {
		sb.WriteString(fmt.Sprintf("<h3>Added Columns (%d)</h3>\n<ul>\n", len(diff.AddedColumns)))
		for _, col := range diff.AddedColumns {
			sb.WriteString(fmt.Sprintf("<li>%s: %s</li>\n", col.Name, col.Type))
		}
		sb.WriteString("</ul>\n")
	}

	if len(diff.RemovedColumns) > 0 {
		sb.WriteString(fmt.Sprintf("<h3>Removed Columns (%d)</h3>\n<ul>\n", len(diff.RemovedColumns)))
		for _, col := range diff.RemovedColumns {
			sb.WriteString(fmt.Sprintf("<li>%s: %s</li>\n", col.Name, col.Type))
		}
		sb.WriteString("</ul>\n")
	}

	if len(diff.ChangedColumns) > 0 {
		sb.WriteString(fmt.Sprintf("<h3>Changed Columns (%d)</h3>\n<ul>\n", len(diff.ChangedColumns)))
		for _, change := range diff.ChangedColumns {
			sb.WriteString(fmt.Sprintf("<li>%s: %s → %s</li>\n", change.Name, change.OldType, change.NewType))
		}
		sb.WriteString("</ul>\n")
	}

	sb.WriteString("</body>\n</html>\n")
	return sb.String(), nil
}

// JSONFormatter formats output as JSON
type JSONFormatter struct{}

func (f *JSONFormatter) FormatSchema(schema *inspect.Schema) (string, error) {
	data := map[string]interface{}{
		"format":      schema.Format,
		"rowCount":    schema.RowCount,
		"columnCount": schema.ColumnCount,
		"columns":     schema.Columns,
	}

	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}

func (f *JSONFormatter) FormatDiff(diff *inspect.SchemaDiff, schema1, schema2 *inspect.Schema) (string, error) {
	data := map[string]interface{}{
		"file1": map[string]interface{}{
			"format":      schema1.Format,
			"rowCount":    schema1.RowCount,
			"columnCount": schema1.ColumnCount,
		},
		"file2": map[string]interface{}{
			"format":      schema2.Format,
			"rowCount":    schema2.RowCount,
			"columnCount": schema2.ColumnCount,
		},
		"addedColumns":   diff.AddedColumns,
		"removedColumns": diff.RemovedColumns,
		"changedColumns": diff.ChangedColumns,
		"sameColumns":    diff.SameColumns,
	}

	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", err
	}

	return string(jsonBytes), nil
}

// GetFormatter returns a formatter for the given output type
func GetFormatter(outputType string) (Formatter, error) {
	switch strings.ToLower(outputType) {
	case "markdown", "md":
		return &MarkdownFormatter{}, nil
	case "html":
		return &HTMLFormatter{}, nil
	case "json":
		return &JSONFormatter{}, nil
	default:
		return nil, fmt.Errorf("unsupported output format: %s (supported: markdown, html, json)", outputType)
	}
}

// WriteOutput writes formatted output to a file or stdout
func WriteOutput(content string, outputPath string) error {
	if outputPath == "" || outputPath == "-" {
		fmt.Print(content)
		return nil
	}

	return os.WriteFile(outputPath, []byte(content), 0644)
}

