package logger

import (
	"testing"
	"fmt"
	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
)

// Mocking the hclog.Logger to intercept log outputs
type MockLogger struct {
	Level hclog.Level
	Color hclog.ColorOption
}

func (m *MockLogger) Log(level hclog.Level, msg string, args ...interface{}) {
	// No-op for testing, you can expand this to capture logs if needed
}

func (m *MockLogger) IsEnabledFor(level hclog.Level) bool {
	return level >= m.Level
}

func (m *MockLogger) Debug(msg string, args ...interface{}) {
	if m.IsEnabledFor(hclog.Debug) {
		fmt.Printf(msg, args...)
	}
}

func (m *MockLogger) Info(msg string, args ...interface{}) {
	if m.IsEnabledFor(hclog.Info) {
		fmt.Printf(msg, args...)
	}
}

func (m *MockLogger) Warn(msg string, args ...interface{}) {
	if m.IsEnabledFor(hclog.Warn) {
		fmt.Printf(msg, args...)
	}
}

func (m *MockLogger) Error(msg string, args ...interface{}) {
	if m.IsEnabledFor(hclog.Error) {
		fmt.Printf(msg, args...)
	}
}

func (m *MockLogger) Trace(msg string, args ...interface{}) {
	if m.IsEnabledFor(hclog.Trace) {
		fmt.Printf(msg, args...)
	}
}

// Unit Test to check InitLogging behavior
func TestInitLogging(t *testing.T) {
	// Capturing output from fmt.Println and fmt.Printf
	// Redirect stdout to capture the logs
	t.Run("Test Debug Logging", func(t *testing.T) {
		// Initialize logging with debug set to true and colorization
		InitLogging(true, true, false)

		// Check if the global LogLevel is set to "DEBUG"
		assert.Equal(t, "DEBUG", LogLevel, "Log level should be DEBUG")

		// Check if the colorization option is set correctly (this would need custom verification in a real test)
		assert.NotNil(t, Logger, "Logger should not be nil")

		// For mocking Logger's output, you can assert the logs printed based on fmt.Printf (or a custom logger)
		// Example assertion for fmt.Println call in InitLogging()
		// check if colorization is applied properly
	})

	t.Run("Test Info Logging", func(t *testing.T) {
		// Initialize logging with debug set to false and no colorization
		InitLogging(false, false, false)

		// Assert that the global LogLevel is "INFO"
		assert.Equal(t, "INFO", LogLevel, "Log level should be INFO")

		// Check if the colorization option is turned off
		assert.NotNil(t, Logger, "Logger should not be nil")

		// Example output check
		// Capture the print output (fmt.Println or fmt.Printf used in InitLogging)
	})

	t.Run("Test Colorization", func(t *testing.T) {
		// Initialize logging with colorization enabled
		InitLogging(true, true, false)

		// Assertions for the color option behavior can be added here
		// Custom logger behavior for color or checks for the colorization part
	})

	t.Run("Test Without Colorization", func(t *testing.T) {
		// Initialize logging with colorization disabled
		InitLogging(false, false, false)

		// Assertions for the color option behavior can be added here
		// Ensure colorization is turned off
	})

	// You can mock out the hclog.Logger to check actual log levels and messages if needed
}
