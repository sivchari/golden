// Package comparator provides advanced comparison logic for golden tests.
package comparator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strings"
)

// Comparator handles advanced comparison logic.
type Comparator struct {
	options Options
}

// Options configures comparison behavior.
type Options struct {
	IgnoreOrder       bool
	IgnoreWhitespace  bool
	CustomCompareFunc func(expected, actual []byte) bool
	IgnoreFields      []string
}

// CompareResult represents the result of a comparison.
type CompareResult struct {
	Equal   bool
	Details string
}

// New creates a new Comparator with default options.
func New() *Comparator {
	return &Comparator{
		options: Options{
			IgnoreOrder:      false,
			IgnoreWhitespace: false,
		},
	}
}

// NewWithOptions creates a new Comparator with custom options.
func NewWithOptions(opts Options) *Comparator {
	return &Comparator{options: opts}
}

// Compare compares two byte arrays with advanced logic.
func (c *Comparator) Compare(expected, actual []byte) *CompareResult {
	// Use custom comparison function if provided
	if c.options.CustomCompareFunc != nil {
		equal := c.options.CustomCompareFunc(expected, actual)

		return &CompareResult{
			Equal:   equal,
			Details: "Custom comparison function",
		}
	}

	// Try JSON comparison first
	if c.isJSON(expected) && c.isJSON(actual) {
		return c.compareJSON(expected, actual)
	}

	// Fall back to text comparison
	return c.compareText(expected, actual)
}

// isJSON checks if data is valid JSON.
func (c *Comparator) isJSON(data []byte) bool {
	data = bytes.TrimSpace(data)
	if len(data) == 0 {
		return false
	}

	first := data[0]

	return first == '{' || first == '['
}

// compareJSON performs semantic JSON comparison.
func (c *Comparator) compareJSON(expected, actual []byte) *CompareResult {
	var expectedObj, actualObj interface{}

	if err := json.Unmarshal(expected, &expectedObj); err != nil {
		return &CompareResult{
			Equal:   false,
			Details: fmt.Sprintf("Failed to parse expected JSON: %v", err),
		}
	}

	if err := json.Unmarshal(actual, &actualObj); err != nil {
		return &CompareResult{
			Equal:   false,
			Details: fmt.Sprintf("Failed to parse actual JSON: %v", err),
		}
	}

	// Normalize both objects
	expectedNorm := c.normalizeValue(expectedObj)
	actualNorm := c.normalizeValue(actualObj)

	equal := c.deepEqual(expectedNorm, actualNorm)

	return &CompareResult{
		Equal:   equal,
		Details: "JSON semantic comparison",
	}
}

// compareText performs text comparison with preprocessing.
func (c *Comparator) compareText(expected, actual []byte) *CompareResult {
	expectedStr := string(expected)
	actualStr := string(actual)

	// Apply text preprocessing
	expectedStr = c.preprocessText(expectedStr)
	actualStr = c.preprocessText(actualStr)

	equal := expectedStr == actualStr

	return &CompareResult{
		Equal:   equal,
		Details: "Text comparison with preprocessing",
	}
}

// normalizeValue normalizes a JSON value for comparison.
func (c *Comparator) normalizeValue(v interface{}) interface{} {
	switch val := v.(type) {
	case map[string]interface{}:
		return c.normalizeObject(val)
	case []interface{}:
		return c.normalizeArray(val)
	case string:
		return c.normalizeString(val)
	default:
		return val
	}
}

// normalizeObject normalizes a JSON object.
func (c *Comparator) normalizeObject(obj map[string]interface{}) map[string]interface{} {
	normalized := make(map[string]interface{})

	for key, value := range obj {
		// Skip ignored fields
		if c.shouldIgnoreField(key) {
			continue
		}

		normalized[key] = c.normalizeValue(value)
	}

	return normalized
}

// normalizeArray normalizes a JSON array.
func (c *Comparator) normalizeArray(arr []interface{}) interface{} {
	normalized := make([]interface{}, len(arr))

	for i, value := range arr {
		normalized[i] = c.normalizeValue(value)
	}

	// Sort array if order should be ignored
	if c.options.IgnoreOrder {
		sort.Slice(normalized, func(i, j int) bool {
			return c.compareValues(normalized[i], normalized[j]) < 0
		})
	}

	return normalized
}

// normalizeString normalizes a string value.
func (c *Comparator) normalizeString(s string) string {
	// Ignore whitespace if configured
	if c.options.IgnoreWhitespace {
		s = strings.TrimSpace(s)
		s = regexp.MustCompile(`\s+`).ReplaceAllString(s, " ")
	}

	return s
}

// shouldIgnoreField checks if a field should be ignored.
func (c *Comparator) shouldIgnoreField(field string) bool {
	for _, ignored := range c.options.IgnoreFields {
		if field == ignored {
			return true
		}
	}

	return false
}

// preprocessText applies text preprocessing options.
func (c *Comparator) preprocessText(s string) string {
	if c.options.IgnoreWhitespace {
		s = strings.TrimSpace(s)
		s = regexp.MustCompile(`\s+`).ReplaceAllString(s, " ")
	}

	return s
}

// deepEqual performs deep equality comparison.
func (c *Comparator) deepEqual(a, b interface{}) bool {
	return reflect.DeepEqual(a, b)
}

// compareValues compares two values for sorting.
func (c *Comparator) compareValues(a, b interface{}) int {
	aStr := fmt.Sprintf("%v", a)
	bStr := fmt.Sprintf("%v", b)

	if aStr < bStr {
		return -1
	} else if aStr > bStr {
		return 1
	}

	return 0
}
