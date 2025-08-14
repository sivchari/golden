package golden

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Options configures Golden test behavior.
type Options struct {
	// Basic settings
	Update bool   // Update mode to create/update golden files
	Dir    string // Golden file directory (defaults to "testdata")

	// Advanced settings
	IgnoreOrder   bool                               // Array order handling (default: true for JSON)
	IgnoreFields  []string                           // Specific JSON fields to ignore
	CustomCompare func(expected, actual []byte) bool // Custom comparison function

	// Internal settings
	contextLines int       // Lines of context in diff
	bufferSize   int       // Buffer size for file operations
	maxFileSize  int64     // Safety limit
	input        io.Reader // For testing
	output       io.Writer // For testing
}

// Option is a functional option for Golden.
type Option func(*Options)

// Basic options

// WithUpdate enables update mode to create/update golden files.
func WithUpdate(update bool) Option {
	return func(o *Options) {
		o.Update = update
	}
}

// WithDir sets custom golden file directory within testdata (default: "testdata").
// If the provided directory doesn't start with "testdata", it will be placed under testdata.
func WithDir(dir string) Option {
	return func(o *Options) {
		// Always ensure the directory is under testdata to avoid Go build issues
		if dir == "" {
			o.Dir = "testdata"
		} else if filepath.IsAbs(dir) {
			// If it's an absolute path (for testing), use it as-is
			o.Dir = dir
		} else if strings.HasPrefix(dir, "testdata") {
			o.Dir = dir
		} else {
			o.Dir = filepath.Join("testdata", dir)
		}
	}
}

// Advanced options

// WithIgnoreFields ignores specific JSON fields during comparison
// Example: WithIgnoreFields("created_at", "updated_at", "id").
func WithIgnoreFields(fields ...string) Option {
	return func(o *Options) {
		o.IgnoreFields = fields
	}
}

// WithIgnoreOrder controls array order sensitivity (default: true for JSON).
func WithIgnoreOrder(ignore bool) Option {
	return func(o *Options) {
		o.IgnoreOrder = ignore
	}
}

// WithCustomCompare sets custom comparison function for special cases.
func WithCustomCompare(fn func(expected, actual []byte) bool) Option {
	return func(o *Options) {
		o.CustomCompare = fn
	}
}

// defaultOptions returns default configuration.
func defaultOptions() *Options {
	return &Options{
		// Default values
		Dir:    "testdata",
		Update: isUpdateModeFromEnv(), // Check GOLDEN_UPDATE environment variable

		// JSON comparison defaults
		IgnoreOrder: true, // Ignore array order for JSON

		// Internal settings
		contextLines: 3,                // Context lines in diff
		bufferSize:   8192,             // File buffer size
		maxFileSize:  50 * 1024 * 1024, // 50MB safety limit
		input:        os.Stdin,
		output:       os.Stdout,
	}
}

// isUpdateModeFromEnv checks if update mode is enabled via GOLDEN_UPDATE environment variable.
func isUpdateModeFromEnv() bool {
	env := os.Getenv("GOLDEN_UPDATE")
	if env == "" {
		return false
	}

	env = strings.ToLower(strings.TrimSpace(env))

	return env == "true"
}
