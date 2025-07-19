// Package golden provides a modern, simple golden test library for Go.
package golden

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/sivchari/golden/comparator"
	"github.com/sivchari/golden/differ"
	"github.com/sivchari/golden/manager"
)

// Golden is the main structure for golden testing.
type Golden struct {
	t          *testing.T
	options    *Options
	manager    *manager.Manager
	comparator *comparator.Comparator
	differ     *differ.Differ
}

// New creates a new Golden instance.
func New(t *testing.T, opts ...Option) *Golden {
	t.Helper()

	options := defaultOptions()
	for _, opt := range opts {
		opt(options)
	}

	// Get test file and function name
	testFile, testFunc := getTestInfo()

	mgr := manager.New(options.Dir, testFile, testFunc)

	// Create comparator with smart options
	compOpts := comparator.Options{
		IgnoreOrder:       options.IgnoreOrder,
		IgnoreFields:      options.IgnoreFields,
		CustomCompareFunc: options.CustomCompare,
	}
	comp := comparator.NewWithOptions(compOpts)

	// Create differ with optimized options
	diffOpts := differ.Options{
		ContextLines:    options.contextLines,
		ColorOutput:     options.colorOutput,
		ShowLineNumbers: true,
		Algorithm:       differ.AlgorithmSimple,
	}
	diff := differ.NewWithOptions(diffOpts)

	return &Golden{
		t:          t,
		options:    options,
		manager:    mgr,
		comparator: comp,
		differ:     diff,
	}
}

// Assert compares any value with the golden file (main API)
// Automatically detects the type and formats appropriately with beautiful diff output.
func (g *Golden) Assert(name string, actual interface{}) {
	// Convert actual value to formatted bytes
	actualBytes := g.formatValue(actual)
	g.assertBytes(name, actualBytes)
}

// formatValue converts any value to a well-formatted byte representation.
func (g *Golden) formatValue(value interface{}) []byte {
	switch v := value.(type) {
	case []byte:
		// If it's already bytes, check if it's JSON
		if g.isJSON(v) {
			return g.formatJSON(v)
		}

		return v
	case string:
		// If it's a string, check if it's JSON
		data := []byte(v)
		if g.isJSON(data) {
			return g.formatJSON(data)
		}

		return data
	case nil:
		return []byte("null")
	default:
		// Try to marshal as JSON (works for structs, maps, slices, etc.)
		if jsonBytes, err := json.MarshalIndent(v, "", "  "); err == nil {
			return jsonBytes
		}
		// Fall back to Go's default string representation
		return []byte(fmt.Sprintf("%+v", v))
	}
}

// isJSON checks if data appears to be JSON.
func (g *Golden) isJSON(data []byte) bool {
	if len(data) == 0 {
		return false
	}

	data = bytes.TrimSpace(data)
	if len(data) == 0 {
		return false
	}

	first := data[0]

	return first == '{' || first == '['
}

// formatJSON ensures JSON is consistently formatted.
func (g *Golden) formatJSON(jsonData []byte) []byte {
	var parsed interface{}
	if err := json.Unmarshal(jsonData, &parsed); err != nil {
		return jsonData // Return as-is if not valid JSON
	}

	formatted, err := json.MarshalIndent(parsed, "", "  ")
	if err != nil {
		return jsonData // Return as-is if formatting fails
	}

	return formatted
}

// assertBytes is the internal implementation.
func (g *Golden) assertBytes(name string, actual []byte) {
	filename := g.manager.GetFilename(name)

	if g.options.Update {
		if err := g.manager.WriteFile(filename, actual); err != nil {
			g.t.Fatalf("Failed to write golden file %s: %v", filename, err)
		}

		return
	}

	expected, err := g.manager.ReadFile(filename)
	if err != nil {
		// If file doesn't exist and we're not in update mode, suggest update mode
		if os.IsNotExist(err) {
			g.t.Fatalf("Golden file %s does not exist. Run with update mode to create it.", filename)
		}

		g.t.Fatalf("Failed to read golden file %s: %v", filename, err)
	}

	// Use advanced comparison
	result := g.comparator.Compare(expected, actual)
	if !result.Equal {
		// Generate beautiful diff output
		diff := g.differ.Diff(expected, actual)
		diffOutput := g.differ.Format(diff)

		// Create beautiful error message with diff
		errorMsg := g.formatDiffError(filename, diffOutput)
		g.t.Fatalf("%s", errorMsg)
	}
}

// formatDiffError creates a beautiful error message with diff.
func (g *Golden) formatDiffError(filename, diffOutput string) string {
	var buf strings.Builder

	// Header with emojis and colors
	if g.options.colorOutput {
		buf.WriteString("🔍 \033[1;31mGolden test failed\033[0m\n")
		buf.WriteString(fmt.Sprintf("📁 File: \033[1;36m%s\033[0m\n", filename))
		buf.WriteString("\n")
		buf.WriteString("🔄 \033[1;33mDifferences found:\033[0m\n")
		buf.WriteString(strings.Repeat("─", 80))
		buf.WriteString("\n")
	} else {
		buf.WriteString("Golden test failed\n")
		buf.WriteString(fmt.Sprintf("File: %s\n", filename))
		buf.WriteString("\n")
		buf.WriteString("Differences found:\n")
		buf.WriteString(strings.Repeat("-", 80))
		buf.WriteString("\n")
	}

	// Add the diff output
	buf.WriteString(diffOutput)

	// Footer
	if g.options.colorOutput {
		buf.WriteString(strings.Repeat("─", 80))
		buf.WriteString("\n")
		buf.WriteString("💡 \033[1;32mTip: Run with update mode to accept changes\033[0m\n")
	} else {
		buf.WriteString(strings.Repeat("-", 80))
		buf.WriteString("\n")
		buf.WriteString("Tip: Run with update mode to accept changes\n")
	}

	return buf.String()
}

// getTestInfo extracts test file and function information from runtime.
func getTestInfo() (string, string) {
	pc := make([]uintptr, 10)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])

	for {
		frame, more := frames.Next()
		if strings.Contains(frame.Function, "Test") {
			file := filepath.Base(frame.File)
			funcParts := strings.Split(frame.Function, ".")
			funcName := funcParts[len(funcParts)-1]

			return file, funcName
		}

		if !more {
			break
		}
	}

	return "unknown_test.go", "UnknownTest"
}
