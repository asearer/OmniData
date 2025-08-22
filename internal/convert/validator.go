package convert

import (
	"fmt"
	"os"
)

// ValidateFormats checks if source and target formats are supported
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

// ResolvePaths normalizes STDIN/STDOUT and checks files
func ResolvePaths(opts Options) (string, string, error) {
	in := opts.InputFile
	out := opts.OutputFile

	// Input
	if in == "-" {
		in = "/dev/stdin"
	} else {
		if _, err := os.Stat(in); os.IsNotExist(err) {
			return "", "", fmt.Errorf("input file does not exist: %s", in)
		}
	}

	// Output
	if out == "-" {
		out = "/dev/stdout"
	} else if !opts.DryRun {
		if _, err := os.Stat(out); err == nil {
			return "", "", fmt.Errorf("output file already exists: %s (use --force to overwrite)", out)
		}
	}

	return in, out, nil
}
