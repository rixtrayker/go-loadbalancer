package logging

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Logger wraps the underlying logging implementation
type Logger struct {
	logger *logrus.Logger
}

// NewLogger creates a new logger instance
func NewLogger() *Logger {
	log := logrus.New()
	log.SetOutput(os.Stdout)
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// Set log level from environment variable or default to info
	level := os.Getenv("LOG_LEVEL")
	switch level {
	case "debug":
		log.SetLevel(logrus.DebugLevel)
	case "warn":
		log.SetLevel(logrus.WarnLevel)
	case "error":
		log.SetLevel(logrus.ErrorLevel)
	default:
		log.SetLevel(logrus.InfoLevel)
	}

	return &Logger{
		logger: log,
	}
}

// Info logs an info message
func (l *Logger) Info(msg string, args ...interface{}) {
	if len(args) == 0 {
		l.logger.Info(msg)
		return
	}

	fields := makeFields(args...)
	l.logger.WithFields(fields).Info(msg)
}

// Debug logs a debug message
func (l *Logger) Debug(msg string, args ...interface{}) {
	if len(args) == 0 {
		l.logger.Debug(msg)
		return
	}

	fields := makeFields(args...)
	l.logger.WithFields(fields).Debug(msg)
}

// Warn logs a warning message
func (l *Logger) Warn(msg string, args ...interface{}) {
	if len(args) == 0 {
		l.logger.Warn(msg)
		return
	}

	fields := makeFields(args...)
	l.logger.WithFields(fields).Warn(msg)
}

// Error logs an error message
func (l *Logger) Error(msg string, args ...interface{}) {
	if len(args) == 0 {
		l.logger.Error(msg)
		return
	}

	fields := makeFields(args...)
	l.logger.WithFields(fields).Error(msg)
}

// makeFields converts a list of key-value pairs to logrus fields
func makeFields(args ...interface{}) logrus.Fields {
	fields := logrus.Fields{}
	for i := 0; i < len(args); i += 2 {
		if i+1 < len(args) {
			fields[args[i].(string)] = args[i+1]
		}
	}
	return fields
}
