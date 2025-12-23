package golden

import (
	"io"
	"os"
	"strings"
)

// Options configures Golden test behavior.
type Options struct {
	// Basic settings
	Update bool // Update mode to create/update golden files

	// Advanced settings
	IgnoreOrder   bool                               // Array order handling (default: true for JSON)
	IgnoreFields  []string                           // Specific JSON fields to ignore
	CustomCompare func(expected, actual []byte) bool // Custom comparison function

	// Path settings
	BaseDir string // Base directory for golden files (default: "testdata")

	// Internal settings
	contextLines int       // Lines of context in diff
	bufferSize   int       // Buffer size for file operations
	maxFileSize  int64     // Safety limit
	input        io.Reader // For testing
	output       io.Writer // For testing
}

// Option is a functional option for Golden.
type Option func(*Options)

// WithUpdate enables update mode to create/update golden files.
func WithUpdate(update bool) Option {
	return func(o *Options) {
		o.Update = update
	}
}

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

// WithBaseDir sets a custom base directory for golden files.
// Default is "testdata".
func WithBaseDir(dir string) Option {
	return func(o *Options) {
		o.BaseDir = dir
	}
}

// defaultOptions returns default configuration.
func defaultOptions() *Options {
	return &Options{
		// Default values
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
