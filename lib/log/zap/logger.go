package main

import (
	"context"
	"fmt"
	"runtime"
	"time"
)

type Logger interface {
	Debug(ctx context.Context, args ...interface{})
	Debugf(ctx context.Context, format string, args ...interface{})
	Info(ctx context.Context, args ...interface{})
	Infof(ctx context.Context, format string, args ...interface{})
	Warn(ctx context.Context, args ...interface{})
	Warnf(ctx context.Context, format string, args ...interface{})
	Error(ctx context.Context, args ...interface{})
	Errorf(ctx context.Context, format string, args ...interface{})
	Fatal(ctx context.Context, args ...interface{})
	Fatalf(ctx context.Context, format string, args ...interface{})
}

var defaultLogger = NewDefaultLogger()

const (
	green  = "\033[32m"
	yellow = "\033[33m"
	red    = "\033[31m"
	black  = "\033[30m"
	reset  = "\033[0m"
)

func colorForLevel(level string) string {
	switch level {
	case "INFO":
		return green
	case "WARN":
		return yellow
	case "ERROR":
		return red
	case "FATAL":
		return black
	default:
		return ""
	}
}

func formatMessage(level string, format string, args ...interface{}) string {
	color := colorForLevel(level)
	timestamp := time.Now().Format("01-02-2006 15:04:05.000")
	_, file, line, _ := runtime.Caller(1) // Increase the caller depth to account for the additional function call
	numGoroutines := runtime.NumGoroutine()
	fileInfo := fmt.Sprintf("%s:%d", file, line)
	message := fmt.Sprintf(format, args...)
	return fmt.Sprintf("%s [%s%s%s] [%d] [%s] - %s", timestamp, color, level, reset, numGoroutines, fileInfo, message)
}

func Debug(ctx context.Context, args ...interface{}) {
	defaultLogger.Debug(ctx, formatMessage("DEBUG", "%v", args...))
}

func Debugf(ctx context.Context, format string, args ...interface{}) {
	defaultLogger.Debug(ctx, formatMessage("DEBUG", format, args...))
}

func Info(ctx context.Context, args ...interface{}) {
	defaultLogger.Info(ctx, formatMessage("INFO", "%v", args...))
}

func Infof(ctx context.Context, format string, args ...interface{}) {
	defaultLogger.Infof(ctx, formatMessage("INFO", format, args...))
}

func Warn(ctx context.Context, args ...interface{}) {
	defaultLogger.Warn(ctx, formatMessage("WARN", "%v", args...))
}

func Warnf(ctx context.Context, format string, args ...interface{}) {
	defaultLogger.Warnf(ctx, formatMessage("WARN", format, args...))
}

func Error(ctx context.Context, args ...interface{}) {
	defaultLogger.Error(ctx, formatMessage("ERROR", "%v", args...))
}

func Errorf(ctx context.Context, format string, args ...interface{}) {
	defaultLogger.Errorf(ctx, formatMessage("ERROR", format, args...))
}

func Fatal(ctx context.Context, args ...interface{}) {
	defaultLogger.Fatal(ctx, formatMessage("FATAL", "%v", args...))
}

func Fatalf(ctx context.Context, format string, args ...interface{}) {
	defaultLogger.Fatalf(ctx, formatMessage("FATAL", format, args...))
}

func main() {
	ctx := context.Background()
	for i := 1; i <= 1000; i++ {
		Infof(ctx, "This is an info message %d", i)
	}

	// Debug(ctx, "This is a debug message")
	// Debugf(ctx, "This is a debug message with args: %s", "argument")
	// Info(ctx, "This is an info message")
	// Infof(ctx, "This is an info message with args: %s", "argument")
	// Warn(ctx, "This is a warn message")
	// Warnf(ctx, "This is a warn message with args: %s", "argument")
	// Error(ctx, "This is an error message")
	// Errorf(ctx, "This is an error message with args: %s", "argument")
	// Fatal and Fatalf will terminate the program, comment out for testing other logs
	// Fatal(ctx, "This is a fatal message")
	// Fatalf(ctx, "This is a fatal message with args: %s", "argument")
}
