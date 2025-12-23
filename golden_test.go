package golden

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGoldenFileCreationAndComparison(t *testing.T) {
	t.Parallel()

	// Create golden file
	g := New(t, WithUpdate(true))
	testData := "test content"
	g.Assert("test_file", testData)

	// Compare with existing golden file (should pass)
	g = New(t, WithUpdate(false))
	g.Assert("test_file", testData)
}

func TestGoldenJSONFormatting(t *testing.T) {
	t.Parallel()

	g := New(t, WithUpdate(true))

	// Test struct as JSON
	type TestData struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	data := TestData{Name: "test", Value: 42}
	g.Assert("json_test", data)

	// Verify comparison works
	g = New(t, WithUpdate(false))
	g.Assert("json_test", data)
}

func TestGoldenIgnoreFields(t *testing.T) {
	t.Parallel()

	// Create golden file with ignored fields
	g := New(t, WithUpdate(true), WithIgnoreFields("timestamp"))
	original := map[string]interface{}{
		"user":      "john",
		"timestamp": "2024-01-01T10:00:00Z",
	}
	g.Assert("ignore_test", original)

	// Test with different timestamp (should pass because timestamp is ignored)
	g = New(t, WithUpdate(false), WithIgnoreFields("timestamp"))
	modified := map[string]interface{}{
		"user":      "john",
		"timestamp": "2024-12-31T23:59:59Z",
	}
	g.Assert("ignore_test", modified)
}

func TestGoldenEnvironmentVariable(t *testing.T) {
	// Test GOLDEN_UPDATE environment variable
	t.Setenv("GOLDEN_UPDATE", "true")

	g := New(t)
	g.Assert("env_test", "test data")

	// Verify file was created
	expectedPath := filepath.Join("testdata", "golden_test_TestGoldenEnvironmentVariable_env_test.golden.go")
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Fatalf("Golden file was not created when GOLDEN_UPDATE=true")
	}
}

func TestGoldenWithBaseDir(t *testing.T) {
	t.Parallel()

	// Create a temporary directory for the custom base dir
	customDir := t.TempDir()

	// Create golden file in custom directory
	g := New(t, WithUpdate(true), WithBaseDir(customDir))
	testData := "custom dir test content"
	g.Assert("basedir_test", testData)

	// Verify file was created in custom directory
	expectedPath := filepath.Join(customDir, "golden_test_TestGoldenWithBaseDir_basedir_test.golden.go")
	if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
		t.Fatalf("Golden file was not created in custom base directory: %s", expectedPath)
	}

	// Compare with existing golden file (should pass)
	g = New(t, WithUpdate(false), WithBaseDir(customDir))
	g.Assert("basedir_test", testData)
}
