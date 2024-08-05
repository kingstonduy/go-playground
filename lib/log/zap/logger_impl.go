package main

import (
	"context"

	"github.com/gammazero/workerpool"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type DefaultLogger struct {
	log        *zap.SugaredLogger
	workerpool workerpool.WorkerPool
}

func NewDefaultLogger() Logger {
	cfg := zap.Config{
		Encoding:    "console",
		Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
		OutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:       "message",
			ConsoleSeparator: " ",
		},
	}

	logger, _ := cfg.Build()
	sugar := logger.Sugar()
	return &DefaultLogger{
		log:        sugar,
		workerpool: *workerpool.New(1000),
	}
}

// DefaultLogger methods implementation
func (d *DefaultLogger) Debug(ctx context.Context, args ...interface{}) {
	d.workerpool.SubmitWait(func() {
		d.log.Debug(args...)
	})
}

func (d *DefaultLogger) Debugf(ctx context.Context, format string, args ...interface{}) {
	d.workerpool.SubmitWait(func() {
		d.log.Debugf(format, args...)
	})
}

func (d *DefaultLogger) Info(ctx context.Context, args ...interface{}) {
	d.workerpool.SubmitWait(func() {
		d.log.Info(args...)
	})
}

func (d *DefaultLogger) Infof(ctx context.Context, format string, args ...interface{}) {
	d.workerpool.SubmitWait(func() {
		d.log.Infof(format, args...)
	})
}

func (d *DefaultLogger) Warn(ctx context.Context, args ...interface{}) {
	d.workerpool.SubmitWait(func() {
		d.log.Warn(args...)
	})
}

func (d *DefaultLogger) Warnf(ctx context.Context, format string, args ...interface{}) {
	d.workerpool.SubmitWait(func() {
		d.log.Warnf(format, args...)
	})
}

func (d *DefaultLogger) Error(ctx context.Context, args ...interface{}) {
	d.workerpool.SubmitWait(func() {
		d.log.Error(args...)
	})
}

func (d *DefaultLogger) Errorf(ctx context.Context, format string, args ...interface{}) {
	d.workerpool.SubmitWait(func() {
		d.log.Errorf(format, args...)
	})
}

func (d *DefaultLogger) Fatal(ctx context.Context, args ...interface{}) {
	d.workerpool.SubmitWait(func() {
		d.log.Fatal(args...)
	})
}

func (d *DefaultLogger) Fatalf(ctx context.Context, format string, args ...interface{}) {
	d.workerpool.SubmitWait(func() {
		d.log.Fatalf(format, args...)
	})
}
