package mock

import (
	"fmt"
	"sync"

	"github.com/zenfulcode/commercify/internal/infrastructure/logger"
)

// Logger is a mock implementation of logger.Logger for testing
type Logger struct {
	mu          sync.RWMutex
	logs        []LogEntry
	ShouldFatal bool // If true, Fatal will actually panic for testing
}

// LogEntry represents a single log entry
type LogEntry struct {
	Level   string
	Message string
	Args    []interface{}
}

// NewLogger creates a new mock logger
func NewLogger() logger.Logger {
	return &Logger{
		logs: make([]LogEntry, 0),
	}
}

// Debug logs a debug message
func (m *Logger) Debug(msg string, args ...interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.logs = append(m.logs, LogEntry{
		Level:   "DEBUG",
		Message: msg,
		Args:    args,
	})
}

// Info logs an info message
func (m *Logger) Info(msg string, args ...interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.logs = append(m.logs, LogEntry{
		Level:   "INFO",
		Message: msg,
		Args:    args,
	})
}

// Warn logs a warning message
func (m *Logger) Warn(msg string, args ...interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.logs = append(m.logs, LogEntry{
		Level:   "WARN",
		Message: msg,
		Args:    args,
	})
}

// Error logs an error message
func (m *Logger) Error(msg string, args ...interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.logs = append(m.logs, LogEntry{
		Level:   "ERROR",
		Message: msg,
		Args:    args,
	})
}

// Fatal logs a fatal message and optionally panics
func (m *Logger) Fatal(msg string, args ...interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.logs = append(m.logs, LogEntry{
		Level:   "FATAL",
		Message: msg,
		Args:    args,
	})

	if m.ShouldFatal {
		panic(fmt.Sprintf(msg, args...))
	}
}

// GetLogs returns all logged entries (for testing purposes)
func (m *Logger) GetLogs() []LogEntry {
	m.mu.RLock()
	defer m.mu.RUnlock()

	logs := make([]LogEntry, len(m.logs))
	copy(logs, m.logs)
	return logs
}

// GetLogsByLevel returns logs filtered by level (for testing purposes)
func (m *Logger) GetLogsByLevel(level string) []LogEntry {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var filtered []LogEntry
	for _, log := range m.logs {
		if log.Level == level {
			filtered = append(filtered, log)
		}
	}
	return filtered
}

// HasLogWithMessage checks if there's a log entry with the specified message (for testing purposes)
func (m *Logger) HasLogWithMessage(message string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, log := range m.logs {
		if log.Message == message {
			return true
		}
	}
	return false
}

// HasLogWithLevel checks if there's a log entry with the specified level (for testing purposes)
func (m *Logger) HasLogWithLevel(level string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, log := range m.logs {
		if log.Level == level {
			return true
		}
	}
	return false
}

// Clear clears all logged entries (for testing purposes)
func (m *Logger) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.logs = make([]LogEntry, 0)
}

// SetShouldFatal sets whether Fatal should actually panic (for testing purposes)
func (m *Logger) SetShouldFatal(shouldFatal bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.ShouldFatal = shouldFatal
}

// GetLastLog returns the last logged entry, or nil if no logs (for testing purposes)
func (m *Logger) GetLastLog() *LogEntry {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if len(m.logs) == 0 {
		return nil
	}

	log := m.logs[len(m.logs)-1]
	return &log
}

// LogCount returns the total number of log entries (for testing purposes)
func (m *Logger) LogCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return len(m.logs)
}
