package logger

import (
	"context"
)

// loggerContextKey is custom context-logger
type loggerContextKey struct{}

// FromContext extract logger from context or return new logger instance.
func FromContext(ctx context.Context) *Logger {
	if ctx != nil {
		if v, ok := ctx.Value(loggerContextKey{}).(*Logger); ok {
			return v
		}
	}
	return New()
}

// ContextWithLogger insert logger into context.
func ContextWithLogger(ctx context.Context, l *Logger) context.Context {
	return context.WithValue(ctx, loggerContextKey{}, l)
}
