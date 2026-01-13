# Testing Guidelines

Quality is central to GitSynq. We strive for high test coverage and reliable tests.

## 1. Unit Tests

- **Location:** Place tests in the same package as the code they test (e.g., `pkg/utils/utils_test.go`).
- **Table-Driven Tests:** Use table-driven tests for functions with multiple inputs and outputs.
- **Mocks:** Use interfaces to mock external dependencies like the SSH client or the filesystem.

## 2. Integration Tests

Integration tests verify that the different components of GitSynq work together.

- **Mock SSH:** We use a mock SSH server for testing push and pull operations without requiring a real remote machine.
- **Git Mocking:** We use temporary directories to create mock Git repositories for testing bundle operations.

## 3. Running Tests

```bash
# Run all tests
make test

# Run a specific test
go test ./internal/bundle -run TestCreateFull

# Run tests with coverage
go test ./... -cover
```

## 4. Writing a New Test

1. Define the test cases in a struct.
2. Iterate over the test cases using `t.Run`.
3. Use `reflect.DeepEqual` for comparing complex structs or slices.
4. Clean up any temporary files or directories after the test.

## Example

```go
func TestMyFunction(t *testing.T) {
    tests := []struct {
        name string
        input string
        want string
    }{
        {"test 1", "in1", "out1"},
        {"test 2", "in2", "out2"},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := MyFunction(tt.input)
            if got != tt.want {
                t.Errorf("got %v, want %v", got, tt.want)
            }
        })
    }
}
```
