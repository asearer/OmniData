package inspect

import (
	"fmt"
	"os"

	"omnidata/internal/convert"
)

// DiffOptions holds configuration for the diff command
type DiffOptions struct {
	File1   string
	File2   string
	Format1 string
	Format2 string
}

// SchemaDiff represents differences between two schemas
type SchemaDiff struct {
	AddedColumns   []ColumnInfo
	RemovedColumns []ColumnInfo
	ChangedColumns []ColumnChange
	SameColumns    []ColumnInfo
}

// ColumnChange represents a change in a column
type ColumnChange struct {
	Name        string
	OldType     string
	NewType     string
	OldNullable bool
	NewNullable bool
}

// RunDiff compares two data files and shows schema differences
func RunDiff(opts DiffOptions) error {
	// Get format handlers
	handler1, ok := convert.GetFormat(opts.Format1)
	if !ok {
		return fmt.Errorf("unsupported format: %s", opts.Format1)
	}

	handler2, ok := convert.GetFormat(opts.Format2)
	if !ok {
		return fmt.Errorf("unsupported format: %s", opts.Format2)
	}

	// Resolve input paths
	path1 := opts.File1
	if path1 == "-" {
		return fmt.Errorf("STDIN not supported for diff command")
	}
	if _, err := os.Stat(path1); err != nil {
		return fmt.Errorf("file1 does not exist: %s", path1)
	}

	path2 := opts.File2
	if path2 == "-" {
		return fmt.Errorf("STDIN not supported for diff command")
	}
	if _, err := os.Stat(path2); err != nil {
		return fmt.Errorf("file2 does not exist: %s", path2)
	}

	// Read data from both files
	data1, err := handler1.ReaderFn(path1)
	if err != nil {
		return fmt.Errorf("failed to read file1: %w", err)
	}

	data2, err := handler2.ReaderFn(path2)
	if err != nil {
		return fmt.Errorf("failed to read file2: %w", err)
	}

	// Infer schemas
	schema1, err := InferSchema(data1, opts.Format1)
	if err != nil {
		return fmt.Errorf("failed to infer schema for file1: %w", err)
	}

	schema2, err := InferSchema(data2, opts.Format2)
	if err != nil {
		return fmt.Errorf("failed to infer schema for file2: %w", err)
	}

	// Compare schemas
	diff := CompareSchemas(schema1, schema2)

	// Display results
	displayDiff(schema1, schema2, diff, opts)

	return nil
}

// CompareSchemas compares two schemas and returns the differences
func CompareSchemas(schema1, schema2 *Schema) *SchemaDiff {
	diff := &SchemaDiff{
		AddedColumns:   make([]ColumnInfo, 0),
		RemovedColumns: make([]ColumnInfo, 0),
		ChangedColumns: make([]ColumnChange, 0),
		SameColumns:    make([]ColumnInfo, 0),
	}

	// Create maps for quick lookup
	cols1 := make(map[string]ColumnInfo)
	for _, col := range schema1.Columns {
		cols1[col.Name] = col
	}

	cols2 := make(map[string]ColumnInfo)
	for _, col := range schema2.Columns {
		cols2[col.Name] = col
	}

	// Find added, removed, changed, and same columns
	for name, col1 := range cols1 {
		if col2, exists := cols2[name]; exists {
			// Column exists in both schemas
			if col1.Type != col2.Type || col1.Nullable != col2.Nullable {
				diff.ChangedColumns = append(diff.ChangedColumns, ColumnChange{
					Name:        name,
					OldType:     col1.Type,
					NewType:     col2.Type,
					OldNullable: col1.Nullable,
					NewNullable: col2.Nullable,
				})
			} else {
				diff.SameColumns = append(diff.SameColumns, col1)
			}
		} else {
			// Column removed
			diff.RemovedColumns = append(diff.RemovedColumns, col1)
		}
	}

	// Find added columns
	for name, col2 := range cols2 {
		if _, exists := cols1[name]; !exists {
			diff.AddedColumns = append(diff.AddedColumns, col2)
		}
	}

	return diff
}

func displayDiff(schema1, schema2 *Schema, diff *SchemaDiff, opts DiffOptions) {
	fmt.Printf("\n🔍 Schema Comparison\n")
	fmt.Printf("═══════════════════════════════════════\n")
	fmt.Printf("File 1: %s (%s)\n", opts.File1, schema1.Format)
	fmt.Printf("  Rows:    %d\n", schema1.RowCount)
	fmt.Printf("  Columns: %d\n", schema1.ColumnCount)
	fmt.Printf("\n")
	fmt.Printf("File 2: %s (%s)\n", opts.File2, schema2.Format)
	fmt.Printf("  Rows:    %d\n", schema2.RowCount)
	fmt.Printf("  Columns: %d\n", schema2.ColumnCount)
	fmt.Printf("\n")

	// Show differences
	if len(diff.AddedColumns) > 0 {
		fmt.Printf("➕ Added Columns (%d)\n", len(diff.AddedColumns))
		fmt.Printf("═══════════════════════════════════════\n")
		for _, col := range diff.AddedColumns {
			nullable := ""
			if col.Nullable {
				nullable = " (nullable)"
			}
			fmt.Printf("  + %s: %s%s\n", col.Name, col.Type, nullable)
		}
		fmt.Printf("\n")
	}

	if len(diff.RemovedColumns) > 0 {
		fmt.Printf("➖ Removed Columns (%d)\n", len(diff.RemovedColumns))
		fmt.Printf("═══════════════════════════════════════\n")
		for _, col := range diff.RemovedColumns {
			nullable := ""
			if col.Nullable {
				nullable = " (nullable)"
			}
			fmt.Printf("  - %s: %s%s\n", col.Name, col.Type, nullable)
		}
		fmt.Printf("\n")
	}

	if len(diff.ChangedColumns) > 0 {
		fmt.Printf("🔄 Changed Columns (%d)\n", len(diff.ChangedColumns))
		fmt.Printf("═══════════════════════════════════════\n")
		for _, change := range diff.ChangedColumns {
			fmt.Printf("  ~ %s\n", change.Name)
			if change.OldType != change.NewType {
				fmt.Printf("      Type: %s → %s\n", change.OldType, change.NewType)
			}
			if change.OldNullable != change.NewNullable {
				oldVal := "not nullable"
				newVal := "nullable"
				if change.OldNullable {
					oldVal = "nullable"
					newVal = "not nullable"
				}
				fmt.Printf("      Nullable: %s → %s\n", oldVal, newVal)
			}
		}
		fmt.Printf("\n")
	}

	if len(diff.SameColumns) > 0 {
		fmt.Printf("✓ Unchanged Columns (%d)\n", len(diff.SameColumns))
		fmt.Printf("═══════════════════════════════════════\n")
		for _, col := range diff.SameColumns {
			nullable := ""
			if col.Nullable {
				nullable = " (nullable)"
			}
			fmt.Printf("  ✓ %s: %s%s\n", col.Name, col.Type, nullable)
		}
		fmt.Printf("\n")
	}

	// Summary
	totalChanges := len(diff.AddedColumns) + len(diff.RemovedColumns) + len(diff.ChangedColumns)
	if totalChanges == 0 {
		fmt.Printf("✨ Schemas are identical!\n")
	} else {
		fmt.Printf("📊 Summary: %d change(s) detected\n", totalChanges)
	}
}
