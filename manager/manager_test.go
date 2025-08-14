package manager

import (
	"testing"
)

func TestNamingStrategy(t *testing.T) {
	t.Parallel()

	naming := &DefaultNaming{}

	// Test GenerateFilename
	tests := []struct {
		testFile   string
		testFunc   string
		goldenName string
		expected   string
	}{
		{"test.go", "TestBasic", "output", "test_TestBasic_output.golden.go"},
	}

	for _, tt := range tests {
		result := naming.GenerateFilename(tt.testFile, tt.testFunc, tt.goldenName)
		if result != tt.expected {
			t.Errorf("GenerateFilename(%s, %s, %s) = %s, want %s",
				tt.testFile, tt.testFunc, tt.goldenName, result, tt.expected)
		}
	}

	// Test ParseFilename
	testFile, testFunc, goldenName, err := naming.ParseFilename("test_TestBasic_output.golden.go")
	if err != nil {
		t.Fatalf("ParseFilename() error = %v", err)
	}

	if testFile != "test.go" || testFunc != "TestBasic" || goldenName != "output" {
		t.Errorf("ParseFilename() = (%s, %s, %s), want (test.go, TestBasic, output)",
			testFile, testFunc, goldenName)
	}
}
