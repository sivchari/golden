// Package manager handles golden file operations and naming strategies.
package manager

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Manager handles golden file operations.
type Manager struct {
	baseDir  string
	testFile string
	testFunc string

	// File naming strategy
	naming NamingStrategy

	// Thread safety
	mu    sync.RWMutex
	locks map[string]*sync.RWMutex
}

// NamingStrategy defines how golden files are named.
type NamingStrategy interface {
	GenerateFilename(testFile, testFunc, goldenName string) string
	ParseFilename(filename string) (testFile, testFunc, goldenName string, err error)
}

// New creates a new Manager.
func New(baseDir, testFile, testFunc string) *Manager {
	return &Manager{
		baseDir:  baseDir,
		testFile: testFile,
		testFunc: testFunc,
		naming:   &DefaultNaming{},
		locks:    make(map[string]*sync.RWMutex),
	}
}

// GetFilename generates the full path for a golden file.
func (m *Manager) GetFilename(goldenName string) string {
	filename := m.naming.GenerateFilename(m.testFile, m.testFunc, goldenName)

	return filepath.Join(m.baseDir, filename)
}

// ReadFile reads a golden file.
func (m *Manager) ReadFile(filename string) ([]byte, error) {
	unlock := m.lockFile(filename, false)
	defer unlock()

	data, err := os.ReadFile(filename) //nolint:gosec // G304: File reading is necessary for golden file functionality
	if err != nil {
		return nil, fmt.Errorf("failed to read golden file %s: %w", filename, err)
	}

	return data, nil
}

// WriteFile writes data to a golden file.
func (m *Manager) WriteFile(filename string, data []byte) error {
	unlock := m.lockFile(filename, true)
	defer unlock()

	// Ensure directory exists
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", dir, err)
	}

	// Write to temporary file first for atomic operation
	tmpFile := filename + ".tmp"
	if err := os.WriteFile(tmpFile, data, 0o600); err != nil {
		return fmt.Errorf("failed to write temporary file %s: %w", tmpFile, err)
	}

	// Atomically move temporary file to final location
	if err := os.Rename(tmpFile, filename); err != nil {
		_ = os.Remove(tmpFile) // Clean up on failure, ignore error

		return fmt.Errorf("failed to rename %s to %s: %w", tmpFile, filename, err)
	}

	return nil
}

// lockFile provides thread-safe file operations.
func (m *Manager) lockFile(filename string, exclusive bool) func() {
	m.mu.Lock()

	lock, exists := m.locks[filename]
	if !exists {
		lock = &sync.RWMutex{}
		m.locks[filename] = lock
	}
	m.mu.Unlock()

	if exclusive {
		lock.Lock()

		return func() { lock.Unlock() }
	}

	lock.RLock()

	return func() { lock.RUnlock() }
}

// DefaultNaming implements the default naming strategy
// Format: TestFunction_goldenName.golden.
type DefaultNaming struct{}

// GenerateFilename generates a filename using the default strategy.
func (dn *DefaultNaming) GenerateFilename(testFile, testFunc, goldenName string) string {
	// Remove .go extension from test file
	baseFile := strings.TrimSuffix(testFile, ".go")

	// Generate filename: TestFile_TestFunction_goldenName.golden
	return fmt.Sprintf("%s_%s_%s.golden", baseFile, testFunc, goldenName)
}

// ParseFilename parses a filename to extract components.
func (dn *DefaultNaming) ParseFilename(filename string) (testFile, testFunc, goldenName string, err error) {
	// Remove .golden extension
	base := strings.TrimSuffix(filename, ".golden")

	// Split by underscore
	parts := strings.Split(base, "_")
	if len(parts) < 3 {
		return "", "", "", fmt.Errorf("invalid filename format: %s", filename)
	}

	// Last part is golden name, everything else is test file and function
	goldenName = parts[len(parts)-1]
	testFunc = parts[len(parts)-2]
	testFile = strings.Join(parts[:len(parts)-2], "_") + ".go"

	return testFile, testFunc, goldenName, nil
}
