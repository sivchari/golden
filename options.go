package golden

import (
	"io"
	"os"
)

// Options configures Golden test behavior with intelligent defaults.
type Options struct {
	// Essential settings (what users actually need)
	Update bool   // Update mode - the only setting most users need
	Dir    string // Golden file directory (defaults to "testdata")

	// Advanced settings (smart defaults, but customizable)
	IgnoreOrder   bool                               // Smart array order handling (default: true for JSON)
	IgnoreFields  []string                           // Specific JSON fields to ignore
	CustomCompare func(expected, actual []byte) bool // For advanced users only

	// Internal (automatically optimized)
	colorOutput  bool      // Auto-detected from terminal
	contextLines int       // Optimized for readability
	bufferSize   int       // Performance optimized
	maxFileSize  int64     // Safety limit
	input        io.Reader // For testing
	output       io.Writer // For testing
}

// Option is a functional option for Golden.
type Option func(*Options)

// Essential options (what 95% of users need)

// WithUpdate enables update mode to create/update golden files.
func WithUpdate(update bool) Option {
	return func(o *Options) {
		o.Update = update
	}
}

// WithDir sets custom golden file directory (default: "testdata").
func WithDir(dir string) Option {
	return func(o *Options) {
		o.Dir = dir
	}
}

// Advanced options (for power users)

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

// defaultOptions returns intelligent defaults optimized for the best experience.
func defaultOptions() *Options {
	return &Options{
		// Essential defaults
		Dir:    "testdata",
		Update: false,

		// Smart defaults for better experience
		IgnoreOrder: true, // Most JSON comparisons don't care about array order

		// Optimized internal settings
		colorOutput:  isTerminal(),     // Auto-detect color support
		contextLines: 3,                // Good balance of context
		bufferSize:   8192,             // Optimal for most file sizes
		maxFileSize:  50 * 1024 * 1024, // 50MB safety limit
		input:        os.Stdin,
		output:       os.Stdout,
	}
}

// isTerminal detects if output supports colors.
func isTerminal() bool {
	// Simple heuristic - check if stdout is a terminal
	if fileInfo, err := os.Stdout.Stat(); err == nil {
		return (fileInfo.Mode() & os.ModeCharDevice) != 0
	}

	return false
}
