package logger

import (
	"context"
)

var staticLogger Logger

func Setup(logger Logger) {
	staticLogger = logger
}

func With(args ...any) Logger {
	return staticLogger.With(args...)
}

func Info(msg string, args ...any) {
	staticLogger.Info(msg, args...)
}

func Warn(msg string, args ...any) {
	staticLogger.Warn(msg, args...)
}

func Error(msg string, args ...any) {
	staticLogger.Error(msg, args...)
}

func Debug(msg string, args ...any) {
	staticLogger.Debug(msg, args...)
}

func InfoContext(ctx context.Context, msg string, args ...any) {
	staticLogger.InfoContext(ctx, msg, args...)
}

func WarnContext(ctx context.Context, msg string, args ...any) {
	staticLogger.WarnContext(ctx, msg, args...)
}

func ErrorContext(ctx context.Context, msg string, args ...any) {
	staticLogger.ErrorContext(ctx, msg, args...)
}

func DebugContext(ctx context.Context, msg string, args ...any) {
	staticLogger.DebugContext(ctx, msg, args...)
}
