package data

/*
Data package for OmniData.

Responsibilities:
- Implement actual data handling logic (reading, converting, writing).
- Provide helper functions for supported formats (CSV, JSON).
- Expose Convert(inputPath, outputPath, from, to) as the main entry point.
- Extensible for new formats with minimal changes.
*/

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// ConverterFunc defines the signature for a conversion function
type ConverterFunc func(input io.Reader, output io.Writer) error

// registry maps "fromFormat -> toFormat -> function"
var converters = map[string]map[string]ConverterFunc{}

// RegisterConverter adds a new conversion function
func RegisterConverter(fromFormat, toFormat string, fn ConverterFunc) {
	fromFormat = strings.ToLower(fromFormat)
	toFormat = strings.ToLower(toFormat)

	if converters[fromFormat] == nil {
		converters[fromFormat] = make(map[string]ConverterFunc)
	}

	converters[fromFormat][toFormat] = fn
}

// Convert performs a conversion from inputPath -> outputPath
func Convert(inputPath, outputPath, from, to string) error {
	from = strings.ToLower(from)
	to = strings.ToLower(to)

	fn, ok := converters[from][to]
	if !ok {
		return fmt.Errorf("unsupported conversion type: %s -> %s", from, to)
	}

	inFile, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("failed to open input file: %w", err)
	}
	defer inFile.Close()

	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outFile.Close()

	return fn(inFile, outFile)
}
