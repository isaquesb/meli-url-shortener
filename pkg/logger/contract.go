package logger

import "context"

type Logger interface {
	With(args ...any) Logger
	Info(string, ...any)
	Warn(string, ...any)
	Error(string, ...any)
	Debug(string, ...any)
	InfoContext(context.Context, string, ...any)
	WarnContext(context.Context, string, ...any)
	ErrorContext(context.Context, string, ...any)
	DebugContext(context.Context, string, ...any)
}
