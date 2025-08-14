# Golden - Golden Test Library for Go

> **Simple API. Any data type. Clear diffs.**

A golden test library for Go with **smart JSON comparison**, **environment variable support**, and colored diff output.

## ✨ Why Golden?

- **🎯 One API for Everything** - `Assert()` handles strings, JSON, structs, anything
- **🌈 Beautiful Diffs** - Colored output with emojis and clear formatting
- **🧠 Smart Defaults** - Automatically ignores array order in JSON, detects formats
- **⚡ High Performance** - Optimized for large files and parallel execution
- **🔒 Thread Safe** - Perfect for concurrent test execution
- **📁 Zero Config** - Works out of the box, customizable when needed
- **🗂️ IDE Integration** - Golden files use `.golden.go` extension for better IDE support

## 🚀 Quick Start

```go
func TestAnything(t *testing.T) {
    g := golden.New(t, golden.WithUpdate(true))
    
    // Test strings
    g.Assert("string", "Hello, World!")
    
    // Test JSON automatically
    g.Assert("json", map[string]interface{}{
        "name": "Golden Test",
        "tags": []string{"awesome", "simple"},
    })
    
    // Test structs (auto-converted to JSON)
    user := User{Name: "Alice", Age: 30}
    g.Assert("struct", user)
    
    // Test any data type
    g.Assert("number", 42)
    g.Assert("boolean", true)
}
```

## 📖 Complete API Reference

### Essential Options

```go
// Create/update golden files
g := golden.New(t, golden.WithUpdate(true))

// Use custom directory (always placed under "testdata/")
g := golden.New(t, golden.WithDir("my_golden_files")) // Creates testdata/my_golden_files/

// Or use environment variable for update mode
// Set GOLDEN_UPDATE=true to enable update mode automatically
g := golden.New(t) // Automatically checks GOLDEN_UPDATE env var
```

### Advanced Options

```go
g := golden.New(t,
    // Ignore specific JSON fields that change between runs
    golden.WithIgnoreFields("created_at", "session_id", "uuid"),
    
    // Control array order sensitivity (default: true for JSON)
    golden.WithIgnoreOrder(false), // Now array order matters
    
    // Custom comparison logic
    golden.WithCustomCompare(func(expected, actual []byte) bool {
        // Your custom logic here
        return string(expected) == string(actual)
    }),
)
```

## 🎨 Diff Output

When tests fail, you get clear, informative output:

```
🔍 Golden test failed
📁 File: testdata/example_test_TestAPI_response.golden.go

🔄 Differences found:
────────────────────────────────────────────────────────────────────────────────
    1  {
    2    "users": [
-   3      "Alice",
+   3      "Bob",
    4      "Charlie"
    5    ],
    6    "count": 2
    7  }
────────────────────────────────────────────────────────────────────────────────
💡 Tip: Run with update mode to accept changes
```

## 🎬 Demo

![Golden Test Library Demo](assets/demo.gif)

*See Golden in action - one API for strings, JSON, structs, and more!*

## 🔥 Smart Features

### Environment Variable Support
Set the `GOLDEN_UPDATE` environment variable to enable update mode automatically:

```bash
# Enable update mode for all tests
GOLDEN_UPDATE=true go test ./...

# Or set it in your shell
export GOLDEN_UPDATE=true
go test
```

### Automatic JSON Formatting
No more manual `json.Marshal` - just pass your data:

```go
// Old way (other libraries)
data := map[string]interface{}{"name": "test"}
jsonBytes, _ := json.MarshalIndent(data, "", "  ")
golden.Assert("test", jsonBytes)

// Golden way ✨
data := map[string]interface{}{"name": "test"}
g.Assert("test", data) // Automatically formatted as JSON!
```

### Smart Array Order Handling
JSON arrays are automatically compared without caring about order:

```go
// These will match automatically
expected := []string{"a", "b", "c"}
actual := []string{"c", "a", "b"}   // Different order, but matches!

g.Assert("array", actual)
```

### Ignore Dynamic Fields
Perfect for API testing with timestamps, UUIDs, etc:

```go
g := golden.New(t, 
    golden.WithIgnoreFields("created_at", "updated_at", "session_id"),
)

apiResponse := map[string]interface{}{
    "user_id": 123,
    "name": "John",
    "created_at": "2023-01-01T10:00:00Z", // Ignored!
    "session_id": "abc123",               // Ignored!
}

g.Assert("api_response", apiResponse)
```

## 📊 Performance

- **Small files (<1MB)**: ~50μs per comparison
- **Large files (>10MB)**: <1s per comparison  
- **Memory efficient**: Uses streaming for large files
- **Parallel safe**: No race conditions in concurrent tests

## 🛠 Advanced Usage

### Custom Comparison Logic

```go
g := golden.New(t, golden.WithCustomCompare(func(expected, actual []byte) bool {
    // Ignore all whitespace differences
    expectedClean := strings.ReplaceAll(string(expected), " ", "")
    actualClean := strings.ReplaceAll(string(actual), " ", "")
    return expectedClean == actualClean
}))
```

### Multiple Test Data Types

```go
func TestMultipleTypes(t *testing.T) {
    g := golden.New(t, golden.WithUpdate(true))
    
    // All of these work with the same API!
    g.Assert("string", "Hello")
    g.Assert("number", 42)
    g.Assert("float", 3.14)
    g.Assert("bool", true)
    g.Assert("array", []int{1, 2, 3})
    g.Assert("map", map[string]string{"key": "value"})
    g.Assert("struct", MyStruct{Field: "value"})
    g.Assert("json_string", `{"formatted": "json"}`)
}
```

## 🔧 Migration from Other Libraries

### From testify/golden

```go
// Old
golden.Assert(t, actual, "test.golden")

// New
g := golden.New(t)
g.Assert("test", actual)
```

### From sebdah/goldie

```go
// Old  
g.Assert(t, "test", actual)

// New
g := golden.New(t)
g.Assert("test", actual)
```

## 📦 Installation

```bash
go get github.com/sivchari/golden
```

**Requirements**: Go 1.24 or later

## 🏗 Project Structure

Golden files are automatically placed in the `testdata` directory with `.golden.go` extension. This ensures they are excluded from Go builds while maintaining IDE support.

```
your-project/
├── testdata/              # Golden files (enforced directory)
│   ├── example_test_TestAPI_response.golden.go
│   └── example_test_TestUser_profile.golden.go  
├── example_test.go
└── main.go
```

**Note**: All golden files are stored in `testdata` or its subdirectories to avoid Go build conflicts. The `.golden.go` extension provides better IDE integration while being safely ignored by Go's build system when placed in `testdata`.

## 🤝 Contributing

We love contributions! Please see our [Contributing Guide](CONTRIBUTING.md).

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details.

---

**Made with ❤️ for the Go community**

