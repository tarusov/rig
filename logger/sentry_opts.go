package logger

import (
	"time"

	"github.com/rs/zerolog"
)

// SentryOption is sentry notifier constructor option.
type SentryOption func(*SentryNotifier)

// WithSentryNotifyLevel setup minimum severnity for sentry messages.
// If level description is invalid,
// debug level will be set.
func WithSentryNotifyLevel(level Level) SentryOption {
	return func(sn *SentryNotifier) {
		zl, err := zerolog.ParseLevel(string(level))
		if err != nil {
			zl = zerolog.DebugLevel
		}
		sn.level = zl
	}
}

// WithSentryEnvironment setup notify env flag.
func WithSentryEnvironment(env string) SentryOption {
	return func(sn *SentryNotifier) {
		sn.env = env
	}
}

// WithSentryEnvironment setup notify release flag.
func WithSentryRelease(release string) SentryOption {
	return func(sn *SentryNotifier) {
		sn.release = release
	}
}

// WithSentryTimeout flush timeout.
func WithSentryTimeout(timeout time.Duration) SentryOption {
	return func(sn *SentryNotifier) {
		sn.timeout = timeout
	}
}
