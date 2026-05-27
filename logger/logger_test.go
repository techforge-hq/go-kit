package logger

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSlogAdapter_JSONFormat(t *testing.T) {
	var buf bytes.Buffer
	logger := newSlogAdapter(Config{
		Level:  LevelInfo,
		Format: FormatJSON,
	}, &buf)

	logger.Info("test message", "key", "value", "count", 42)

	output := buf.String()

	// Verify it's valid JSON
	var logEntry map[string]any
	err := json.Unmarshal([]byte(output), &logEntry)
	assert.NoError(t, err, "Expected valid JSON, got error: %v\nOutput: %s", err, output)

	// Verify fields
	assert.Equal(t, "test message", logEntry["msg"], "Expected msg='test message'")
	assert.Equal(t, "value", logEntry["key"], "Expected key='value'")
	assert.Equal(t, float64(42), logEntry["count"], "Expected count=42")
}

func TestSlogAdapter_TextFormat(t *testing.T) {
	var buf bytes.Buffer
	logger := newSlogAdapter(Config{
		Level:  LevelDebug,
		Format: FormatText,
	}, &buf)

	logger.Debug("debug message", "user", "john")

	output := buf.String()

	assert.Contains(t, output, "debug message", "Expected output to contain 'debug message'")
	assert.Contains(t, output, "user=john", "Expected output to contain 'user=john'")
}

func TestSlogAdapter_Levels(t *testing.T) {
	tests := []struct {
		name      string
		level     Level
		logFunc   func(Logger)
		shouldLog bool
	}{
		{
			name:  "Debug level logs debug messages",
			level: LevelDebug,
			logFunc: func(l Logger) {
				l.Debug("debug msg")
			},
			shouldLog: true,
		},
		{
			name:  "Info level does not log debug messages",
			level: LevelInfo,
			logFunc: func(l Logger) {
				l.Debug("debug msg")
			},
			shouldLog: false,
		},
		{
			name:  "Info level logs info messages",
			level: LevelInfo,
			logFunc: func(l Logger) {
				l.Info("info msg")
			},
			shouldLog: true,
		},
		{
			name:  "Error level only logs errors",
			level: LevelError,
			logFunc: func(l Logger) {
				l.Info("info msg")
			},
			shouldLog: false,
		},
		{
			name:  "Error level logs error messages",
			level: LevelError,
			logFunc: func(l Logger) {
				l.Error("error msg")
			},
			shouldLog: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := newSlogAdapter(Config{
				Level:  tt.level,
				Format: FormatText,
			}, &buf)

			tt.logFunc(logger)

			output := buf.String()
			hasOutput := len(output) > 0

			assert.Equal(t, tt.shouldLog, hasOutput, "Expected shouldLog=%v, but got output: %q", tt.shouldLog, output)
		})
	}
}

func TestSlogAdapter_With(t *testing.T) {
	var buf bytes.Buffer
	baseLogger := newSlogAdapter(Config{
		Level:  LevelInfo,
		Format: FormatJSON,
	}, &buf)

	// Create child logger with additional fields
	childLogger := baseLogger.With("request_id", "abc123", "user_id", 42)
	childLogger.Info("processing request")

	output := buf.String()

	var logEntry map[string]interface{}
	err := json.Unmarshal([]byte(output), &logEntry)
	assert.NoError(t, err, "Expected valid JSON, got error: %v", err)

	assert.Equal(t, "abc123", logEntry["request_id"], "Expected request_id='abc123'")
	assert.Equal(t, float64(42), logEntry["user_id"], "Expected user_id=42")
}

func TestSlogAdapter_AllLevels(t *testing.T) {
	var buf bytes.Buffer
	logger := newSlogAdapter(Config{
		Level:  LevelDebug,
		Format: FormatText,
	}, &buf)

	// Test all log levels
	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warn message")
	logger.Error("error message")

	output := buf.String()

	// Verify all messages are present
	expectedMessages := []string{"debug message", "info message", "warn message", "error message"}
	for _, msg := range expectedMessages {
		assert.Contains(t, output, msg, "Expected output to contain %q", msg)
	}
}

func TestNoopLogger(_ *testing.T) {
	logger := NewNoop()

	// These should not panic or produce output
	logger.Debug("debug")
	logger.Info("info")
	logger.Warn("warn")
	logger.Error("error")

	child := logger.With("key", "value")
	child.Info("child log")

	// If we get here without panic, the test passes
}

func TestNewDevelopment(t *testing.T) {
	logger := NewDevelopment()
	assert.NotNil(t, logger, "Expected NewDevelopment to return a logger")

	// Verify it's a SlogAdapter
	assert.IsType(t, SlogAdapter{}, logger, "Expected NewDevelopment to return SlogAdapter")
}

func TestNewProduction(t *testing.T) {
	logger := NewProduction()
	assert.NotNil(t, logger, "Expected NewProduction to return a logger")

	// Verify it's a SlogAdapter
	assert.IsType(t, SlogAdapter{}, logger, "Expected NewProduction to return SlogAdapter")
}

func TestConvertLevel(t *testing.T) {
	tests := []struct {
		name     string
		level    Level
		expected string // We can't directly compare slog.Level, so we'll use string representation
	}{
		{"Debug level", LevelDebug, "DEBUG"},
		{"Info level", LevelInfo, "INFO"},
		{"Warn level", LevelWarn, "WARN"},
		{"Error level", LevelError, "ERROR"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			slogLevel := convertLevel(tt.level)
			assert.Equal(t, tt.expected, slogLevel.String(), "Expected %s, got %s", tt.expected, slogLevel.String())
		})
	}
}
