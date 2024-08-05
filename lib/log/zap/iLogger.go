package logger

import (
	"context"

	"go.uber.org/zap"
)

type Logger interface {
	// Fields set fields to always be logged
	Fields(fields map[string]interface{}) Logger

	Debug(ctx context.Context, args ...interface{})
	// Logf Debug
	Debugf(ctx context.Context, format string, args ...interface{})
	// Log Info
	Info(ctx context.Context, args ...interface{})
	// Logf Info
	Infof(ctx context.Context, format string, args ...interface{})
	// Log Warn
	Warn(ctx context.Context, args ...interface{})
	// Logf Warn
	Warnf(ctx context.Context, format string, args ...interface{})
	// Log Error
	Error(ctx context.Context, args ...interface{})
	// Logf Error
	Errorf(ctx context.Context, format string, args ...interface{})
	// Log Fatal
	Fatal(ctx context.Context, args ...interface{})
	// Logf Fatal
	Fatalf(ctx context.Context, format string, args ...interface{})

	// String returns the name of logger
	String() string
}

type LoggerImpl struct {
	sugar *zap.SugaredLogger
}

func NewLoggerImpl() Logger {
	logger, _ := zap.NewProduction()
	sugar := logger.Sugar()
	return &LoggerImpl{
		sugar: sugar,
	}
}

// Debug implements Logger.
func (l *LoggerImpl) Debug(ctx context.Context, args ...interface{}) {
	panic("unimplemented")
}

// Debugf implements Logger.
func (l *LoggerImpl) Debugf(ctx context.Context, format string, args ...interface{}) {
	panic("unimplemented")
}

// Error implements Logger.
func (l *LoggerImpl) Error(ctx context.Context, args ...interface{}) {
	panic("unimplemented")
}

// Errorf implements Logger.
func (l *LoggerImpl) Errorf(ctx context.Context, format string, args ...interface{}) {
	panic("unimplemented")
}

// Fatal implements Logger.
func (l *LoggerImpl) Fatal(ctx context.Context, args ...interface{}) {
	panic("unimplemented")
}

// Fatalf implements Logger.
func (l *LoggerImpl) Fatalf(ctx context.Context, format string, args ...interface{}) {
	panic("unimplemented")
}

// Fields implements Logger.
func (l *LoggerImpl) Fields(fields map[string]interface{}) Logger {
	panic("unimplemented")
}

// Info implements Logger.
func (l *LoggerImpl) Info(ctx context.Context, args ...interface{}) {
	panic("unimplemented")
}

// Infof implements Logger.
func (l *LoggerImpl) Infof(ctx context.Context, format string, args ...interface{}) {
	panic("unimplemented")
}

// String implements Logger.
func (l *LoggerImpl) String() string {
	panic("unimplemented")
}

// Warn implements Logger.
func (l *LoggerImpl) Warn(ctx context.Context, args ...interface{}) {
	panic("unimplemented")
}

// Warnf implements Logger.
func (l *LoggerImpl) Warnf(ctx context.Context, format string, args ...interface{}) {
	panic("unimplemented")
}
