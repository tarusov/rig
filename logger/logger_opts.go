package logger

import (
	"io"

	"github.com/rs/zerolog"
)

// loggerOption is logger constructor optional modificator.
type loggerOption func(*loggerOptions)

// WithLoggingOutput setup output targets (Stdout, Sentry, files, etc).
// More than one call of this option will replace previous setup.
func WithLoggingOutput(out ...io.Writer) loggerOption {
	return func(l *loggerOptions) {
		if len(out) == 0 {
			return
		}
		l.outputs = out[:]
	}
}

// WithLoggingLevel setup minimum severnity for log messages.
// If level description is invalid - debug level will be set.
func WithLoggingLevel(level Level) loggerOption {
	return func(l *loggerOptions) {
		zl, err := zerolog.ParseLevel(string(level))
		if err != nil {
			zl = zerolog.DebugLevel
		}
		l.level = zl
	}
}

// WithLoggingFormat setup output formatting type.
// If format description are invalid - json format will be set.
func WithLoggingFormat(format Format) loggerOption {
	return func(l *loggerOptions) {
		if format != FormatConsole &&
			format != FormatJSON &&
			format != FormatText {
			format = FormatJSON
		}

		l.format = format
	}
}

// WithLoggingTimestampFormat setup output timestamp formatting (rfc3339).
func WithLoggingTimestampFormat(tsf string) loggerOption {
	return func(l *loggerOptions) {
		l.timestampFormat = tsf
	}
}

// WithLoggingTimestampField setup output timestamp field name (ts, time).
func WithLoggingTimestampField(tsf string) loggerOption {
	return func(l *loggerOptions) {
		l.timestampName = tsf
	}
}
