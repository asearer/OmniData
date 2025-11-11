package formats

import (
	"database/sql"
	"fmt"
	"strings"

	"omnidata/internal/convert"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

// init registers the SQL format handler in the global Registry
func init() {
	convert.RegisterFormat("sql", convert.FormatHandler{
		Name:     "sql",
		ReaderFn: readSQL,
		WriterFn: writeSQL,
	})
}

// SQLConnection holds database connection parameters
type SQLConnection struct {
	Driver string // mysql, postgres, sqlite3
	DSN    string // Data Source Name
	Query  string // SQL query to execute
	Table  string // Table name (alternative to query)
}

// readSQL reads data from a SQL database.
// The path parameter should be a connection string in format:
// "driver://user:password@host:port/database?query=SELECT * FROM table"
// or for SQLite: "sqlite3:///path/to/database.db?table=table_name"
func readSQL(path string) (interface{}, error) {
	if path == "" {
		return nil, fmt.Errorf("SQL read from STDIN is not supported")
	}

	conn, err := parseSQLPath(path)
	if err != nil {
		return nil, fmt.Errorf("failed to parse SQL connection string: %w", err)
	}

	// Open database connection
	db, err := sql.Open(conn.Driver, conn.DSN)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Build query
	query := conn.Query
	if query == "" && conn.Table != "" {
		query = fmt.Sprintf("SELECT * FROM %s", conn.Table)
	}
	if query == "" {
		return nil, fmt.Errorf("no query or table specified")
	}

	// Execute query
	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get columns: %w", err)
	}

	// Read all rows
	result := make([][]string, 0)

	// Add header row
	result = append(result, columns)

	// Read data rows
	for rows.Next() {
		// Create slice of interfaces to hold column values
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		// Scan row into values
		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		// Convert values to strings
		row := make([]string, len(columns))
		for i, val := range values {
			if val == nil {
				row[i] = ""
			} else {
				row[i] = fmt.Sprintf("%v", val)
			}
		}
		result = append(result, row)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return result, nil
}

// writeSQL writes data to a SQL database table.
// The path parameter should be a connection string with table name:
// "driver://user:password@host:port/database?table=table_name"
func writeSQL(path string, data interface{}) error {
	if path == "" {
		return fmt.Errorf("SQL write to STDOUT is not supported")
	}

	records, ok := data.([][]string)
	if !ok {
		return fmt.Errorf("invalid data type for SQL writer, expected [][]string")
	}

	if len(records) == 0 {
		return fmt.Errorf("no data to write")
	}

	conn, err := parseSQLPath(path)
	if err != nil {
		return fmt.Errorf("failed to parse SQL connection string: %w", err)
	}

	if conn.Table == "" {
		return fmt.Errorf("table name is required for SQL write")
	}

	// Open database connection
	db, err := sql.Open(conn.Driver, conn.DSN)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	// Get headers (first row)
	headers := records[0]
	placeholders := make([]string, len(headers))
	for i := range placeholders {
		placeholders[i] = "?"
	}

	// Build INSERT statement
	insertSQL := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		conn.Table,
		strings.Join(headers, ", "),
		strings.Join(placeholders, ", "))

	// Prepare statement
	stmt, err := db.Prepare(insertSQL)
	if err != nil {
		return fmt.Errorf("failed to prepare INSERT statement: %w", err)
	}
	defer stmt.Close()

	// Insert rows
	for i := 1; i < len(records); i++ {
		row := records[i]
		values := make([]interface{}, len(headers))
		for j := 0; j < len(headers) && j < len(row); j++ {
			if row[j] == "" {
				values[j] = nil
			} else {
				values[j] = row[j]
			}
		}
		// Pad with nil if row is shorter than headers
		for j := len(row); j < len(headers); j++ {
			values[j] = nil
		}

		if _, err := stmt.Exec(values...); err != nil {
			return fmt.Errorf("failed to insert row %d: %w", i, err)
		}
	}

	return nil
}

// parseSQLPath parses a SQL connection string
// Format: "driver://dsn?query=SELECT * FROM table&table=table_name"
func parseSQLPath(path string) (*SQLConnection, error) {
	conn := &SQLConnection{}

	// Split driver and DSN
	parts := strings.SplitN(path, "://", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid SQL connection string format")
	}

	conn.Driver = parts[0]
	dsnAndParams := parts[1]

	// Split DSN and query parameters
	paramParts := strings.SplitN(dsnAndParams, "?", 2)
	conn.DSN = paramParts[0]

	// Parse query parameters
	if len(paramParts) > 1 {
		params := strings.Split(paramParts[1], "&")
		for _, param := range params {
			kv := strings.SplitN(param, "=", 2)
			if len(kv) == 2 {
				key := strings.ToLower(kv[0])
				value := kv[1]
				switch key {
				case "query":
					conn.Query = value
				case "table":
					conn.Table = value
				}
			}
		}
	}

	return conn, nil
}
