# Testing Guide for GoModoro

## Running Tests Locally

### Prerequisites
```bash
# Install Go dependencies
go mod download

# Install test tools (optional but recommended)
make install-tools
```

### Quick Test Commands

```bash
# Run all tests
make test

# Run tests with race detection
make test-race

# Run tests with coverage
make test-coverage

# Run specific test
go test -v -run TestSessionManager

# Run tests in specific file
go test -v session_test.go session.go
```

### Using the Makefile

```bash
# Show all available commands
make help

# Run all checks (format, vet, test)
make check

# Watch for changes and auto-run tests (requires entr)
make watch

# Clean build artifacts and test files
make clean
```

## Test Structure

### Unit Tests

1. **session_test.go**
   - Session manager creation
   - Long break positioning algorithm
   - Session progression and skipping
   - Surprise task limits
   - JSON serialization

2. **timer_test.go**
   - Time formatting
   - Timer state transitions
   - Reset/next/skip functionality
   - Edge cases and error handling

3. **integration_test.go**
   - Full Pomodoro workflow
   - Settings updates
   - Concurrent operations
   - Session completion flow

## CI/CD with GitHub Actions

The project includes a GitHub Actions workflow that:

1. **Tests on every push/PR**
   - Runs all unit tests
   - Checks race conditions
   - Verifies code formatting
   - Runs go vet

2. **Build verification**
   - Builds for Linux
   - Verifies cross-compilation

3. **Code quality**
   - Runs golangci-lint
   - Checks formatting

### Setting up GitHub Actions

1. Create `.github/workflows/` directory in your repo
2. Add the `test.yml` file
3. Push to GitHub
4. Tests will run automatically on push/PR

## Writing New Tests

### Test Naming Convention
```go
// Function tests
func TestFunctionName(t *testing.T) {}

// Method tests  
func TestTypeName_MethodName(t *testing.T) {}

// Integration tests
func TestFeature_Integration(t *testing.T) {}
```

### Table-Driven Tests Example
```go
func TestFormatTime(t *testing.T) {
    tests := []struct {
        name     string
        input    int
        expected string
    }{
        {"zero", 0, "00:00"},
        {"one minute", 60, "01:00"},
        {"one hour", 3600, "60:00"},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := formatTime(tt.input)
            if result != tt.expected {
                t.Errorf("got %s, want %s", result, tt.expected)
            }
        })
    }
}
```

## Coverage Goals

- Aim for >80% code coverage
- Focus on critical paths:
  - Session management logic
  - Timer state transitions
  - Settings updates
  - Break positioning algorithm

## Testing on Different Platforms

### Desktop (Ubuntu)
```bash
go test -v ./...
```

### Android (before deployment)
```bash
# Build and test
make package-android

# Install on device
make install-android

# Manual testing required for UI
```

## Common Test Issues

### 1. Fyne UI Tests
Fyne UI components require a display, so UI tests may fail in headless environments. Use build tags to skip UI tests in CI:

```go
// +build !ci

package main

import "testing"

func TestUIComponent(t *testing.T) {
    // UI-specific tests
}
```

### 2. Random Failures
The surprise task feature uses randomness. Tests account for this by:
- Running multiple iterations
- Checking bounds rather than exact values
- Using deterministic seeds where possible

### 3. Timing Issues
Integration tests may be flaky due to timing. Use:
```go
if testing.Short() {
    t.Skip("Skipping integration test in short mode")
}
```

## Debugging Tests

```bash
# Verbose output
go test -v ./...

# Run specific test with more detail
go test -v -run TestSessionManager

# Debug with delve
dlv test -- -test.run TestSessionManager
```

## Pre-commit Checklist

- [ ] Run `make check` - all tests pass
- [ ] Run `make test-coverage` - coverage >80%
- [ ] Run `make fmt` - code is formatted
- [ ] No new lint warnings
- [ ] Tests added for new features
- [ ] Integration tests pass