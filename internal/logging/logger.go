package logging

import (
	"context"
	"io"
	"os"

	"github.com/rixtrayker/go-loadbalancer/configs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.opentelemetry.io/otel/trace"
)

// Logger provides structured logging capabilities
type Logger struct {
	logger *zap.Logger
	sugar  *zap.SugaredLogger
}

// NewLogger creates a new logger instance with default configuration
func NewLogger() *Logger {
	// Create default production config
	config := zap.NewProductionConfig()
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// Create logger
	logger, _ := config.Build(zap.AddCallerSkip(1))
	sugar := logger.Sugar()

	return &Logger{
		logger: logger,
		sugar:  sugar,
	}
}

// Configure sets up the logger based on the provided configuration
func (l *Logger) Configure(config configs.LoggingConfig) error {
	// Configure log level
	var level zapcore.Level
	switch config.Level {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}

	// Configure encoder
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "timestamp",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Configure output
	var output zapcore.WriteSyncer
	switch config.Output {
	case "stdout":
		output = zapcore.AddSync(os.Stdout)
	case "stderr":
		output = zapcore.AddSync(os.Stderr)
	default:
		output = zapcore.AddSync(os.Stdout)
	}

	// Configure format
	var encoder zapcore.Encoder
	switch config.Format {
	case "json":
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	case "console":
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	default:
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	// Create core
	core := zapcore.NewCore(encoder, output, zap.NewAtomicLevelAt(level))

	// Create logger
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	l.logger = logger
	l.sugar = logger.Sugar()

	return nil
}

// WithContext returns a logger with trace context information if available
func (l *Logger) WithContext(ctx context.Context) *Logger {
	if ctx == nil {
		return l
	}

	// Extract trace information if available
	span := trace.SpanFromContext(ctx)
	if !span.SpanContext().IsValid() {
		return l
	}

	// Create a new logger with trace information
	spanCtx := span.SpanContext()
	logger := l.logger.With(
		zap.String("trace_id", spanCtx.TraceID().String()),
		zap.String("span_id", spanCtx.SpanID().String()),
	)

	return &Logger{
		logger: logger,
		sugar:  logger.Sugar(),
	}
}

// Debug logs a debug message
func (l *Logger) Debug(msg string, args ...interface{}) {
	l.sugar.Debugw(msg, args...)
}

// Info logs an info message
func (l *Logger) Info(msg string, args ...interface{}) {
	l.sugar.Infow(msg, args...)
}

// Warn logs a warning message
func (l *Logger) Warn(msg string, args ...interface{}) {
	l.sugar.Warnw(msg, args...)
}

// Error logs an error message
func (l *Logger) Error(msg string, args ...interface{}) {
	l.sugar.Errorw(msg, args...)
}

// Fatal logs a fatal message and exits
func (l *Logger) Fatal(msg string, args ...interface{}) {
	l.sugar.Fatalw(msg, args...)
}

// With returns a logger with the given fields
func (l *Logger) With(args ...interface{}) *Logger {
	return &Logger{
		logger: l.logger,
		sugar:  l.sugar.With(args...),
	}
}

// SetOutput sets the logger output
func (l *Logger) SetOutput(output io.Writer) {
	// Create a new core with the given output
	core := zapcore.NewCore(
		l.logger.Core().Encoder(),
		zapcore.AddSync(output),
		l.logger.Core().Level(),
	)

	// Create a new logger with the new core
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	l.logger = logger
	l.sugar = logger.Sugar()
}

// Sync flushes any buffered log entries
func (l *Logger) Sync() error {
	return l.logger.Sync()
}
