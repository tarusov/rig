package locker

import (
	"time"
)

// lockerOption is locker constructor optional modificator.
type lockerOption func(*Locker)

// WithRetryCount set custom retry count.
func WithRetryCount(n int) lockerOption {
	return func(l *Locker) {
		l.retryCount = n
	}
}

// WithRetryTimeout set custom retry timeout.
func WithRetryTimeout(t time.Duration) lockerOption {
	return func(l *Locker) {
		l.retryTimeout = t
	}
}
