package logger

import (
	"context"
	"log/slog"
	"os"
)

type StdLogger struct {
	stdOut *slog.Logger
	stdErr *slog.Logger
}

func NewLogger(logLevel string) Logger {
	level := slog.LevelInfo
	switch logLevel {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "warning":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}
	options := &slog.HandlerOptions{
		Level: level,
	}
	return &StdLogger{
		stdOut: slog.New(slog.NewJSONHandler(os.Stdout, options)),
		stdErr: slog.New(slog.NewJSONHandler(os.Stdout, options)),
	}
}

func (l *StdLogger) With(args ...any) Logger {
	return &StdLogger{
		stdOut: l.stdOut.With(args...),
		stdErr: l.stdErr.With(args...),
	}
}

func (l *StdLogger) Info(msg string, args ...any) {
	l.stdOut.Info(msg, args...)
}

func (l *StdLogger) Warn(msg string, args ...any) {
	l.stdOut.Warn(msg, args...)
}

func (l *StdLogger) Error(msg string, args ...any) {
	l.stdErr.Error(msg, args...)
}

func (l *StdLogger) Debug(msg string, args ...any) {
	l.stdOut.Debug(msg, args...)
}

func (l *StdLogger) InfoContext(ctx context.Context, msg string, args ...any) {
	l.stdOut.InfoContext(ctx, msg, args...)
}

func (l *StdLogger) WarnContext(ctx context.Context, msg string, args ...any) {
	l.stdOut.WarnContext(ctx, msg, args...)
}

func (l *StdLogger) ErrorContext(ctx context.Context, msg string, args ...any) {
	l.stdErr.ErrorContext(ctx, msg, args...)
}

func (l *StdLogger) DebugContext(ctx context.Context, msg string, args ...any) {
	l.stdOut.DebugContext(ctx, msg, args...)
}
