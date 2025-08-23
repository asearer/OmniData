package convert

import (
	"fmt"
	"os"
	"runtime"
)

// ValidateFormats checks if the source and target formats are supported.
//
// Returns an error if:
// - Either format is not registered in the Registry.
// - Source and target formats are the same.
func ValidateFormats(from, to string) error {
	if _, ok := Registry[from]; !ok {
		return fmt.Errorf("unsupported source format: %s", from)
	}
	if _, ok := Registry[to]; !ok {
		return fmt.Errorf("unsupported target format: %s", to)
	}
	if from == to {
		return fmt.Errorf("source and target formats are the same")
	}
	return nil
}

// ResolvePaths normalizes input/output paths and checks file validity.
//
// Supports:
// - "-" for STDIN/STDOUT (cross-platform).
// - Ensures input file exists and is not a directory.
// - Prevents accidental overwrite of output file unless DryRun is true.
//
// Returns normalized input and output paths (empty string indicates STDIN/STDOUT).
func ResolvePaths(opts Options) (string, string, error) {
	inputPath := opts.InputFile
	outputPath := opts.OutputFile

	// ---------------------------
	// Handle Input
	// ---------------------------
	if inputPath == "-" {
		// Cross-platform STDIN placeholder
		inputPath = ""
	} else {
		info, err := os.Stat(inputPath)
		if err != nil {
			if os.IsNotExist(err) {
				return "", "", fmt.Errorf("input file does not exist: %s", inputPath)
			}
			return "", "", fmt.Errorf("cannot access input file '%s': %w", inputPath, err)
		}
		if info.IsDir() {
			return "", "", fmt.Errorf("input path is a directory: %s", inputPath)
		}
	}

	// ---------------------------
	// Handle Output
	// ---------------------------
	if outputPath == "-" {
		// Cross-platform STDOUT placeholder
		outputPath = ""
	} else if !opts.DryRun {
		if _, err := os.Stat(outputPath); err == nil {
			// File exists: prevent accidental overwrite
			return "", "", fmt.Errorf("output file already exists: %s (use --force to overwrite)", outputPath)
		} else if !os.IsNotExist(err) {
			// Unexpected error accessing file
			return "", "", fmt.Errorf("cannot access output file '%s': %w", outputPath, err)
		}
	}

	// On Windows, normalize STDIN/STDOUT for handlers
	if runtime.GOOS == "windows" {
		if inputPath == "" {
			inputPath = "CON"
		}
		if outputPath == "" {
			outputPath = "CON"
		}
	}

	return inputPath, outputPath, nil
}
