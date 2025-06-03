package logging

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Config represents logging configuration
type Config struct {
	Level      string
	Format     string
	OutputFile string
}

// NewLogger creates a new logger with the given configuration
func NewLogger(config Config) (*logrus.Logger, error) {
	logger := logrus.New()

	// Set log level
	level, err := logrus.ParseLevel(config.Level)
	if err != nil {
		return nil, err
	}
	logger.SetLevel(level)

	// Set log format
	switch config.Format {
	case "json":
		logger.SetFormatter(&logrus.JSONFormatter{})
	case "text":
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	default:
		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}

	// Set output
	if config.OutputFile != "" {
		file, err := os.OpenFile(config.OutputFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return nil, err
		}
		logger.SetOutput(file)
	} else {
		logger.SetOutput(os.Stdout)
	}

	return logger, nil
}

// DefaultLogger creates a logger with default configuration
func DefaultLogger() *logrus.Logger {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	logger.SetOutput(os.Stdout)
	return logger
} 