package examples

import (
	"fmt"
	"testing"
	"time"

	"github.com/sivchari/golden"
)

// Simple golden testing examples focusing on core Assert() functionality

// Basic usage with different data types.
func TestBasicGoldenUsage(t *testing.T) {
	g := golden.New(t, golden.WithUpdate(true))

	// Test strings
	g.Assert("string_test", "Hello, Golden Test!")

	// Test numbers
	g.Assert("number_test", 42)

	// Test booleans
	g.Assert("boolean_test", true)

	// Test arrays/slices
	g.Assert("array_test", []string{"apple", "banana", "cherry"})
}

// JSON data testing.
func TestJSONGolden(t *testing.T) {
	g := golden.New(t, golden.WithUpdate(true))

	// Test simple JSON objects
	user := map[string]interface{}{
		"id":     123,
		"name":   "John Doe",
		"email":  "john@example.com",
		"active": true,
	}
	g.Assert("user_json", user)

	// Test complex nested JSON
	apiResponse := map[string]interface{}{
		"status": "success",
		"data": map[string]interface{}{
			"users": []map[string]interface{}{
				{"id": 1, "name": "Alice"},
				{"id": 2, "name": "Bob"},
			},
			"total": 2,
		},
		"timestamp": "2023-01-01T00:00:00Z",
	}
	g.Assert("api_response", apiResponse)
}

// Testing with struct types.
func TestStructGolden(t *testing.T) {
	g := golden.New(t, golden.WithUpdate(true))

	type User struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	user := User{
		ID:    1,
		Name:  "Alice",
		Email: "alice@example.com",
	}

	// Struct is automatically converted to JSON
	g.Assert("struct_user", user)
}

// Testing with ignore fields option.
func TestIgnoreFieldsGolden(t *testing.T) {
	g := golden.New(t,
		golden.WithUpdate(true),
		golden.WithIgnoreFields("timestamp", "session_id"),
	)

	// Data with dynamic fields that should be ignored
	dynamicData := map[string]interface{}{
		"user_id":    123,
		"action":     "login",
		"timestamp":  time.Now().Format(time.RFC3339), // This will be ignored
		"session_id": "abc123-def456",                 // This will be ignored
		"success":    true,
	}

	g.Assert("dynamic_data", dynamicData)
}

// Testing array order independence for JSON.
func TestArrayOrderGolden(t *testing.T) {
	g := golden.New(t, golden.WithUpdate(true))

	// Arrays in different order should match (default behavior for JSON)
	data1 := map[string]interface{}{
		"tags": []string{"go", "testing", "golden"},
	}

	data2 := map[string]interface{}{
		"tags": []string{"testing", "golden", "go"}, // Different order
	}

	g.Assert("ordered_data", data1)
	// This would match data1 because array order is ignored by default for JSON
	g.Assert("ordered_data", data2)
}

// Simple benchmark for core functionality.
func BenchmarkGoldenAssert(b *testing.B) {
	g := golden.New(&testing.T{}, golden.WithUpdate(true))

	testData := map[string]interface{}{
		"id":      123,
		"message": "benchmark test",
		"data":    []int{1, 2, 3, 4, 5},
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		g.Assert(fmt.Sprintf("bench_%d", i), testData)
	}
}
