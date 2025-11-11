package stream

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

// StreamingReader provides streaming read capabilities for large files
type StreamingReader interface {
	ReadRow() (map[string]string, error)
	Close() error
}

// CSVStreamingReader reads CSV files row by row
type CSVStreamingReader struct {
	file   *os.File
	reader *csv.Reader
	header []string
}

// NewCSVStreamingReader creates a new streaming CSV reader
func NewCSVStreamingReader(path string) (*CSVStreamingReader, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file: %w", err)
	}

	reader := csv.NewReader(bufio.NewReader(file))

	// Read header
	header, err := reader.Read()
	if err != nil {
		file.Close()
		return nil, fmt.Errorf("failed to read CSV header: %w", err)
	}

	return &CSVStreamingReader{
		file:   file,
		reader: reader,
		header: header,
	}, nil
}

// ReadRow reads the next row from the CSV file
func (r *CSVStreamingReader) ReadRow() (map[string]string, error) {
	record, err := r.reader.Read()
	if err == io.EOF {
		return nil, io.EOF
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read CSV row: %w", err)
	}

	row := make(map[string]string)
	for i, value := range record {
		if i < len(r.header) {
			row[r.header[i]] = value
		}
	}

	return row, nil
}

// Close closes the underlying file
func (r *CSVStreamingReader) Close() error {
	return r.file.Close()
}

// JSONStreamingReader reads JSON files with streaming support
// For JSON arrays, reads one object at a time
type JSONStreamingReader struct {
	file    *os.File
	decoder *json.Decoder
	first   bool
}

// NewJSONStreamingReader creates a new streaming JSON reader
func NewJSONStreamingReader(path string) (*JSONStreamingReader, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open JSON file: %w", err)
	}

	decoder := json.NewDecoder(file)

	// Check if it's an array
	token, err := decoder.Token()
	if err != nil {
		file.Close()
		return nil, fmt.Errorf("failed to read JSON: %w", err)
	}

	// If it's not a delimiter, it's a single object
	if delim, ok := token.(json.Delim); !ok || delim != '[' {
		file.Close()
		return nil, fmt.Errorf("JSON streaming only supports arrays")
	}

	return &JSONStreamingReader{
		file:    file,
		decoder: decoder,
		first:   true,
	}, nil
}

// ReadRow reads the next object from the JSON array
func (r *JSONStreamingReader) ReadRow() (map[string]string, error) {
	// Check if we've reached the end
	if !r.decoder.More() {
		// Consume closing bracket
		_, _ = r.decoder.Token()
		return nil, io.EOF
	}

	var obj map[string]interface{}
	if err := r.decoder.Decode(&obj); err != nil {
		return nil, fmt.Errorf("failed to decode JSON object: %w", err)
	}

	// Convert to map[string]string
	row := make(map[string]string)
	for k, v := range obj {
		row[k] = fmt.Sprintf("%v", v)
	}

	return row, nil
}

// Close closes the underlying file
func (r *JSONStreamingReader) Close() error {
	return r.file.Close()
}

// StreamingWriter provides streaming write capabilities
type StreamingWriter interface {
	WriteRow(row map[string]string) error
	Close() error
}

// CSVStreamingWriter writes CSV files row by row
type CSVStreamingWriter struct {
	file          *os.File
	writer        *csv.Writer
	header        []string
	headerWritten bool
}

// NewCSVStreamingWriter creates a new streaming CSV writer
func NewCSVStreamingWriter(path string, header []string) (*CSVStreamingWriter, error) {
	file, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("failed to create CSV file: %w", err)
	}

	writer := csv.NewWriter(file)

	return &CSVStreamingWriter{
		file:          file,
		writer:        writer,
		header:        header,
		headerWritten: false,
	}, nil
}

// WriteRow writes a row to the CSV file
func (w *CSVStreamingWriter) WriteRow(row map[string]string) error {
	if !w.headerWritten {
		if err := w.writer.Write(w.header); err != nil {
			return fmt.Errorf("failed to write CSV header: %w", err)
		}
		w.headerWritten = true
	}

	// Convert row map to slice in header order
	record := make([]string, len(w.header))
	for i, col := range w.header {
		record[i] = row[col]
	}

	if err := w.writer.Write(record); err != nil {
		return fmt.Errorf("failed to write CSV row: %w", err)
	}

	return nil
}

// Close closes the writer and flushes data
func (w *CSVStreamingWriter) Close() error {
	w.writer.Flush()
	if err := w.writer.Error(); err != nil {
		w.file.Close()
		return fmt.Errorf("CSV writer error: %w", err)
	}
	return w.file.Close()
}
