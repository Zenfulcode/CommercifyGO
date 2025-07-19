package logger

import (
	"log"
	"os"
)

// Logger defines the interface for logging
type Logger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Fatal(msg string, args ...interface{})
}

// SimpleLogger implements a simple logger using the standard log package
type SimpleLogger struct {
	debugLogger *log.Logger
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
	fatalLogger *log.Logger
}

// NewLogger creates a new SimpleLogger
func NewLogger() *SimpleLogger {
	return &SimpleLogger{
		debugLogger: log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile),
		infoLogger:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime),
		warnLogger:  log.New(os.Stdout, "WARN: ", log.Ldate|log.Ltime),
		errorLogger: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		fatalLogger: log.New(os.Stderr, "FATAL: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

// Debug logs a debug message
func (l *SimpleLogger) Debug(msg string, args ...interface{}) {
	if len(args) > 0 {
		l.debugLogger.Printf(msg, args...)
	} else {
		l.debugLogger.Println(msg)
	}
}

// Info logs an info message
func (l *SimpleLogger) Info(msg string, args ...interface{}) {
	if len(args) > 0 {
		l.infoLogger.Printf(msg, args...)
	} else {
		l.infoLogger.Println(msg)
	}
}

// Warn logs a warning message
func (l *SimpleLogger) Warn(msg string, args ...interface{}) {
	if len(args) > 0 {
		l.warnLogger.Printf(msg, args...)
	} else {
		l.warnLogger.Println(msg)
	}
}

// Error logs an error message
func (l *SimpleLogger) Error(msg string, args ...interface{}) {
	if len(args) > 0 {
		l.errorLogger.Printf(msg, args...)
	} else {
		l.errorLogger.Println(msg)
	}
}

// Fatal logs a fatal message and exits
func (l *SimpleLogger) Fatal(msg string, args ...interface{}) {
	if len(args) > 0 {
		l.fatalLogger.Fatalf(msg, args...)
	} else {
		l.fatalLogger.Fatalln(msg)
	}
}
