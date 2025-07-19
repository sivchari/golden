package examples

import (
	"testing"

	"github.com/sivchari/golden"
)

// Basic usage example - works with any type!
func TestBasicUsage(t *testing.T) {
	g := golden.New(t, golden.WithUpdate(true))

	// String output
	g.Assert("string_output", "Hello, Golden Test World!")

	// Integer output
	g.Assert("number_output", 42)

	// Boolean output
	g.Assert("bool_output", true)
}

// JSON example - automatic formatting!
func TestJSONOutput(t *testing.T) {
	g := golden.New(t, golden.WithUpdate(true))

	// Just pass the struct/map - it's automatically formatted as JSON
	data := map[string]interface{}{
		"name":    "Golden Test",
		"version": "1.0.0",
		"features": []string{
			"simple API",
			"beautiful diff",
			"high performance",
		},
	}

	// One line - no manual JSON marshaling needed!
	g.Assert("json_output", data)
}

// Smart array order handling - works automatically!
func TestSmartComparison(t *testing.T) {
	g := golden.New(t, golden.WithUpdate(true))
	// Smart default: array order is ignored for JSON automatically!

	data := map[string]interface{}{
		"name":    "Important data that matters",
		"tags":    []string{"c", "a", "b"}, // Order ignored automatically for JSON!
		"numbers": []int{3, 1, 2},          // Order ignored here too!
		"users": []map[string]interface{}{
			{"name": "Bob", "id": 2},
			{"name": "Alice", "id": 1}, // Object order ignored too!
		},
	}

	g.Assert("smart_comparison", data)
}

// Most users only need these two options!
func TestEssentialOptions(t *testing.T) {
	// Option 1: Update mode (create/update golden files)
	g := golden.New(t, golden.WithUpdate(true))
	g.Assert("update_example", "This creates/updates the golden file")

	// Option 2: Custom directory (default is "testdata")
	g2 := golden.New(t, golden.WithUpdate(true), golden.WithDir("custom_golden"))
	g2.Assert("custom_dir", "Files go in custom_golden/ instead")
}

// Power user options: fine-grained control.
func TestAdvancedOptions(t *testing.T) {
	// Ignore specific fields that change between runs
	g := golden.New(t,
		golden.WithUpdate(true),
		golden.WithIgnoreFields("session_id", "request_id", "timestamp", "created_at"),
		golden.WithIgnoreOrder(false), // Care about array order
	)

	apiResponse := map[string]interface{}{
		"user_id":    123,
		"session_id": "abc123xyz",            // Ignored!
		"request_id": "req-789",              // Ignored!
		"timestamp":  "2023-01-01T12:00:00Z", // Ignored!
		"user_data": map[string]interface{}{
			"name":       "John",
			"email":      "john@example.com",
			"created_at": "2023-01-01T10:00:00Z", // Ignored!
		},
		"items": []string{"first", "second"}, // Order matters now
	}

	g.Assert("advanced_options", apiResponse)
}
