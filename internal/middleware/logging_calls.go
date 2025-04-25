package middleware

import (
	"log"
	"os"
	"strings"
)

// LogLevel represents the logging level
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
)

var logLevel LogLevel

func init() {
	level := strings.ToUpper(os.Getenv("LOG_LEVEL"))
	switch level {
	case "DEBUG":
		logLevel = DEBUG
	case "INFO":
		logLevel = INFO
	case "WARN":
		logLevel = WARN
	case "ERROR":
		logLevel = ERROR
	default:
		logLevel = INFO // Default log level
	}
}

// Logger is a utility for logging messages with levels
type Logger struct{}

// Debug logs a debug message
func (l *Logger) Debug(v ...interface{}) {
	if logLevel <= DEBUG {
		log.Println("[DEBUG]", v)
	}
}

// Info logs an info message
func (l *Logger) Info(v ...interface{}) {
	if logLevel <= INFO {
		log.Println("[INFO]", v)
	}
}

// Warn logs a warning message
func (l *Logger) Warn(v ...interface{}) {
	if logLevel <= WARN {
		log.Println("[WARN]", v)
	}
}

// Error logs an error message
func (l *Logger) Error(v ...interface{}) {
	if logLevel <= ERROR {
		log.Println("[ERROR]", v)
	}
}

// NewLogger creates a new instance of Logger
func NewLogger() *Logger {
	return &Logger{}
}
