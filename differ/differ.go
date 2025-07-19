// Package differ provides diff generation and formatting capabilities for golden tests.
package differ

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
)

// Differ handles diff generation and formatting.
type Differ struct {
	options Options
}

// Options configures diff behavior.
type Options struct {
	ContextLines    int
	ColorOutput     bool
	ShowLineNumbers bool
	Algorithm       DiffAlgorithm
}

// DiffAlgorithm specifies the diff algorithm to use.
type DiffAlgorithm int

const (
	// AlgorithmMyers uses the Myers diff algorithm for accurate diffs.
	AlgorithmMyers DiffAlgorithm = iota
	// AlgorithmSimple uses a simple line-by-line comparison.
	AlgorithmSimple
)

// DiffChunk represents a chunk of differences.
type DiffChunk struct {
	Type   ChunkType
	Lines  []string
	StartA int // Start line in A (expected)
	StartB int // Start line in B (actual)
	CountA int // Number of lines in A
	CountB int // Number of lines in B
}

// ChunkType represents the type of diff chunk.
type ChunkType int

const (
	// ChunkEqual represents equal content.
	ChunkEqual ChunkType = iota
	// ChunkDelete represents deleted content.
	ChunkDelete
	// ChunkInsert represents inserted content.
	ChunkInsert
	// ChunkReplace represents replaced content.
	ChunkReplace
)

// Diff represents the complete diff between two texts.
type Diff struct {
	Chunks []DiffChunk
	Equal  bool
}

// New creates a new Differ with default options.
func New() *Differ {
	return &Differ{
		options: Options{
			ContextLines:    3,
			ColorOutput:     true,
			ShowLineNumbers: true,
			Algorithm:       AlgorithmSimple,
		},
	}
}

// NewWithOptions creates a new Differ with custom options.
func NewWithOptions(opts Options) *Differ {
	return &Differ{options: opts}
}

// Diff compares two byte arrays and returns a Diff.
func (d *Differ) Diff(expected, actual []byte) *Diff {
	expectedLines := d.splitLines(expected)
	actualLines := d.splitLines(actual)

	switch d.options.Algorithm {
	case AlgorithmMyers:
		return d.myersDiff(expectedLines, actualLines)
	case AlgorithmSimple:
		return d.simpleDiff(expectedLines, actualLines)
	default:
		return d.simpleDiff(expectedLines, actualLines)
	}
}

// Format formats a diff for display.
func (d *Differ) Format(diff *Diff) string {
	if diff.Equal {
		return ""
	}

	var buf strings.Builder

	for _, chunk := range diff.Chunks {
		switch chunk.Type {
		case ChunkEqual:
			d.formatEqualChunk(&buf, chunk)
		case ChunkDelete:
			d.formatDeleteChunk(&buf, chunk)
		case ChunkInsert:
			d.formatInsertChunk(&buf, chunk)
		case ChunkReplace:
			d.formatReplaceChunk(&buf, chunk)
		}
	}

	return buf.String()
}

// splitLines splits text into lines while preserving line endings.
func (d *Differ) splitLines(data []byte) []string {
	if len(data) == 0 {
		return []string{}
	}

	var lines []string

	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	// Handle case where file doesn't end with newline
	if len(data) > 0 && data[len(data)-1] != '\n' {
		lines = append(lines, "")
	}

	return lines
}

// simpleDiff implements a simple line-by-line diff algorithm.
func (d *Differ) simpleDiff(expected, actual []string) *Diff {
	diff := &Diff{Equal: true}

	maxLen := len(expected)
	if len(actual) > maxLen {
		maxLen = len(actual)
	}

	i := 0
	for i < maxLen {
		switch {
		case i >= len(expected):
			// Extra lines in actual
			chunk := DiffChunk{
				Type:   ChunkInsert,
				Lines:  []string{actual[i]},
				StartB: i,
				CountB: 1,
			}
			diff.Chunks = append(diff.Chunks, chunk)
			diff.Equal = false
		case i >= len(actual):
			// Missing lines in actual
			chunk := DiffChunk{
				Type:   ChunkDelete,
				Lines:  []string{expected[i]},
				StartA: i,
				CountA: 1,
			}
			diff.Chunks = append(diff.Chunks, chunk)
			diff.Equal = false
		case expected[i] == actual[i]:
			// Equal lines
			chunk := DiffChunk{
				Type:   ChunkEqual,
				Lines:  []string{expected[i]},
				StartA: i,
				StartB: i,
				CountA: 1,
				CountB: 1,
			}
			diff.Chunks = append(diff.Chunks, chunk)
		default:
			// Different lines
			chunk := DiffChunk{
				Type:   ChunkReplace,
				Lines:  []string{expected[i], actual[i]},
				StartA: i,
				StartB: i,
				CountA: 1,
				CountB: 1,
			}
			diff.Chunks = append(diff.Chunks, chunk)
			diff.Equal = false
		}

		i++
	}

	return diff
}

// myersDiff implements Myers diff algorithm (simplified version).
func (d *Differ) myersDiff(expected, actual []string) *Diff {
	// For now, fall back to simple diff
	// TODO: Implement full Myers algorithm
	return d.simpleDiff(expected, actual)
}

// formatEqualChunk formats equal lines.
func (d *Differ) formatEqualChunk(buf *strings.Builder, chunk DiffChunk) {
	for i, line := range chunk.Lines {
		lineNum := chunk.StartA + i + 1
		if d.options.ShowLineNumbers {
			fmt.Fprintf(buf, " %4d  %s\n", lineNum, line)
		} else {
			fmt.Fprintf(buf, "  %s\n", line)
		}
	}
}

// formatDeleteChunk formats deleted lines.
func (d *Differ) formatDeleteChunk(buf *strings.Builder, chunk DiffChunk) {
	for i, line := range chunk.Lines {
		lineNum := chunk.StartA + i + 1
		d.writeDeleteLine(buf, line, lineNum)
	}
}

// writeDeleteLine writes a single delete line with appropriate formatting.
func (d *Differ) writeDeleteLine(buf *strings.Builder, line string, lineNum int) {
	switch {
	case d.options.ColorOutput && d.options.ShowLineNumbers:
		fmt.Fprintf(buf, "\033[31m-%4d  %s\033[0m\n", lineNum, line)
	case d.options.ColorOutput:
		fmt.Fprintf(buf, "\033[31m- %s\033[0m\n", line)
	case d.options.ShowLineNumbers:
		fmt.Fprintf(buf, "-%4d  %s\n", lineNum, line)
	default:
		fmt.Fprintf(buf, "- %s\n", line)
	}
}

// formatInsertChunk formats inserted lines.
func (d *Differ) formatInsertChunk(buf *strings.Builder, chunk DiffChunk) {
	for i, line := range chunk.Lines {
		lineNum := chunk.StartB + i + 1
		d.writeInsertLine(buf, line, lineNum)
	}
}

// writeInsertLine writes a single insert line with appropriate formatting.
func (d *Differ) writeInsertLine(buf *strings.Builder, line string, lineNum int) {
	switch {
	case d.options.ColorOutput && d.options.ShowLineNumbers:
		fmt.Fprintf(buf, "\033[32m+%4d  %s\033[0m\n", lineNum, line)
	case d.options.ColorOutput:
		fmt.Fprintf(buf, "\033[32m+ %s\033[0m\n", line)
	case d.options.ShowLineNumbers:
		fmt.Fprintf(buf, "+%4d  %s\n", lineNum, line)
	default:
		fmt.Fprintf(buf, "+ %s\n", line)
	}
}

// formatReplaceChunk formats replaced lines.
func (d *Differ) formatReplaceChunk(buf *strings.Builder, chunk DiffChunk) {
	// Show as delete followed by insert
	expectedLine := chunk.Lines[0]
	actualLine := chunk.Lines[1]
	lineNum := chunk.StartA + 1

	d.writeDeleteLine(buf, expectedLine, lineNum)
	d.writeInsertLine(buf, actualLine, lineNum)
}
