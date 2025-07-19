package golden

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNew(t *testing.T) {
	g := New(t)
	if g == nil {
		t.Fatal("New() returned nil")
	}

	if g.t != t {
		t.Fatal("Golden instance does not have correct testing.T")
	}

	if g.options == nil {
		t.Fatal("Golden instance has nil options")
	}

	if g.manager == nil {
		t.Fatal("Golden instance has nil manager")
	}
}

func TestNewWithOptions(t *testing.T) {
	g := New(t,
		WithDir("custom"),
		WithUpdate(true),
		WithIgnoreOrder(false),
	)

	if g.options.Dir != "custom" {
		t.Errorf("Expected Dir=custom, got %s", g.options.Dir)
	}

	if !g.options.Update {
		t.Error("Expected Update=true")
	}

	if g.options.IgnoreOrder {
		t.Error("Expected IgnoreOrder=false")
	}
}

func TestAssert(t *testing.T) {
	// Setup temporary directory
	tmpDir := t.TempDir()

	g := New(t, WithDir(tmpDir), WithUpdate(true))

	// Test string
	g.Assert("string_test", "Hello, World!")

	// Test number
	g.Assert("number_test", 42)

	// Test JSON object
	data := map[string]interface{}{
		"name":    "test",
		"version": "1.0.0",
	}
	g.Assert("json_test", data)

	// Verify files were created
	stringFile := g.manager.GetFilename("string_test")
	if _, err := os.Stat(stringFile); os.IsNotExist(err) {
		t.Fatalf("Golden file was not created: %s", stringFile)
	}
}

func TestAssertComparison(t *testing.T) {
	// Setup temporary directory
	tmpDir := t.TempDir()

	// Create golden file manually
	goldenContent := "Expected output"

	filename := filepath.Join(tmpDir, "golden_test_TestAssertComparison_comparison.golden")
	if err := os.MkdirAll(filepath.Dir(filename), 0o750); err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}

	if err := os.WriteFile(filename, []byte(goldenContent), 0o600); err != nil {
		t.Fatalf("Failed to write golden file: %v", err)
	}

	g := New(t, WithDir(tmpDir))

	// Test matching content - this should pass
	g.Assert("comparison", goldenContent)
}

func TestAssertStruct(t *testing.T) {
	tmpDir := t.TempDir()
	g := New(t, WithDir(tmpDir), WithUpdate(true))

	type testStruct struct {
		Name  string
		Value int
	}

	expected := testStruct{Name: "test", Value: 42}
	g.Assert("struct_test", expected)

	// Verify file was created
	filename := g.manager.GetFilename("struct_test")
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Fatalf("Golden file was not created: %s", filename)
	}
}

func TestIgnoreFields(t *testing.T) {
	tmpDir := t.TempDir()

	g := New(t, WithDir(tmpDir), WithUpdate(true), WithIgnoreFields("timestamp", "id"))

	// Create golden file
	original := map[string]interface{}{
		"name":      "test",
		"timestamp": "2023-01-01T00:00:00Z",
		"id":        "123",
		"value":     42,
	}
	g.Assert("ignore_fields", original)

	// Test with different ignored fields (should pass)
	g2 := New(t, WithDir(tmpDir), WithIgnoreFields("timestamp", "id"))
	modified := map[string]interface{}{
		"name":      "test",
		"timestamp": "2024-01-01T00:00:00Z", // Different timestamp (ignored)
		"id":        "456",                  // Different ID (ignored)
		"value":     42,                     // Same value (compared)
	}
	g2.Assert("ignore_fields", modified)
}
