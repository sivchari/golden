package manager

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestNew(t *testing.T) {
	m := New("testdata", "test.go", "TestExample")
	if m == nil {
		t.Fatal("New() returned nil")
	}

	if m.baseDir != "testdata" {
		t.Errorf("Expected baseDir=testdata, got %s", m.baseDir)
	}

	if m.testFile != "test.go" {
		t.Errorf("Expected testFile=test.go, got %s", m.testFile)
	}

	if m.testFunc != "TestExample" {
		t.Errorf("Expected testFunc=TestExample, got %s", m.testFunc)
	}
}

func TestGetFilename(t *testing.T) {
	m := New("testdata", "example_test.go", "TestBasic")
	filename := m.GetFilename("output")

	expected := filepath.Join("testdata", "example_test_TestBasic_output.golden")
	if filename != expected {
		t.Errorf("Expected filename=%s, got %s", expected, filename)
	}
}

func TestWriteAndReadFile(t *testing.T) {
	tmpDir := t.TempDir()
	m := New(tmpDir, "test.go", "TestWrite")

	filename := m.GetFilename("write_test")
	content := []byte("Test content for writing")

	// Write file
	if err := m.WriteFile(filename, content); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Fatalf("File was not created: %s", filename)
	}

	// Read file
	readContent, err := m.ReadFile(filename)
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}

	// Compare content
	if !bytes.Equal(content, readContent) {
		t.Fatalf("Content mismatch. Expected: %s, Got: %s", content, readContent)
	}
}

func TestWriteFileAtomic(t *testing.T) {
	tmpDir := t.TempDir()
	m := New(tmpDir, "test.go", "TestAtomic")

	filename := m.GetFilename("atomic_test")
	content := []byte("Atomic write test")

	// Write file
	if err := m.WriteFile(filename, content); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	// Verify no temporary file remains
	tmpFile := filename + ".tmp"
	if _, err := os.Stat(tmpFile); !os.IsNotExist(err) {
		t.Errorf("Temporary file should not exist: %s", tmpFile)
	}

	// Verify actual file exists with correct content
	readContent, err := m.ReadFile(filename)
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}

	if !bytes.Equal(content, readContent) {
		t.Fatalf("Content mismatch. Expected: %s, Got: %s", content, readContent)
	}
}

func TestDefaultNamingGenerateFilename(t *testing.T) {
	naming := &DefaultNaming{}

	tests := []struct {
		testFile   string
		testFunc   string
		goldenName string
		expected   string
	}{
		{"test.go", "TestBasic", "output", "test_TestBasic_output.golden"},
		{"example_test.go", "TestComplex", "result", "example_test_TestComplex_result.golden"},
		{"my_test.go", "TestWithNumbers123", "data", "my_test_TestWithNumbers123_data.golden"},
	}

	for _, tt := range tests {
		result := naming.GenerateFilename(tt.testFile, tt.testFunc, tt.goldenName)
		if result != tt.expected {
			t.Errorf("GenerateFilename(%s, %s, %s) = %s, want %s",
				tt.testFile, tt.testFunc, tt.goldenName, result, tt.expected)
		}
	}
}

func TestDefaultNamingParseFilename(t *testing.T) {
	naming := &DefaultNaming{}

	tests := []struct {
		filename       string
		expectedFile   string
		expectedFunc   string
		expectedGolden string
		expectError    bool
	}{
		{"test_TestBasic_output.golden", "test.go", "TestBasic", "output", false},
		{"example_test_TestComplex_result.golden", "example_test.go", "TestComplex", "result", false},
		{"my_test_TestWithNumbers123_data.golden", "my_test.go", "TestWithNumbers123", "data", false},
		{"invalid.golden", "", "", "", true},
		{"too_short.golden", "", "", "", true},
	}

	for _, tt := range tests {
		testFile, testFunc, goldenName, err := naming.ParseFilename(tt.filename)

		if tt.expectError {
			if err == nil {
				t.Errorf("ParseFilename(%s) expected error, got nil", tt.filename)
			}

			continue
		}

		if err != nil {
			t.Errorf("ParseFilename(%s) unexpected error: %v", tt.filename, err)

			continue
		}

		if testFile != tt.expectedFile {
			t.Errorf("ParseFilename(%s) testFile = %s, want %s", tt.filename, testFile, tt.expectedFile)
		}

		if testFunc != tt.expectedFunc {
			t.Errorf("ParseFilename(%s) testFunc = %s, want %s", tt.filename, testFunc, tt.expectedFunc)
		}

		if goldenName != tt.expectedGolden {
			t.Errorf("ParseFilename(%s) goldenName = %s, want %s", tt.filename, goldenName, tt.expectedGolden)
		}
	}
}

func TestConcurrentWriteRead(t *testing.T) {
	tmpDir := t.TempDir()
	m := New(tmpDir, "concurrent_test.go", "TestConcurrent")

	const numGoroutines = 10

	const numOperations = 100

	// Create channels for synchronization
	done := make(chan bool, numGoroutines)

	// Start multiple goroutines writing different files
	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer func() { done <- true }()

			for j := 0; j < numOperations; j++ {
				filename := m.GetFilename(fmt.Sprintf("concurrent_%d_%d", id, j))
				content := []byte(fmt.Sprintf("Content from goroutine %d, operation %d", id, j))

				if err := m.WriteFile(filename, content); err != nil {
					t.Errorf("WriteFile failed: %v", err)

					return
				}

				readContent, err := m.ReadFile(filename)
				if err != nil {
					t.Errorf("ReadFile failed: %v", err)

					return
				}

				if !bytes.Equal(content, readContent) {
					t.Errorf("Content mismatch in goroutine %d, operation %d", id, j)

					return
				}
			}
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < numGoroutines; i++ {
		<-done
	}
}
