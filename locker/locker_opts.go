package locker

import "time"

// LockerOption is locker constructor optional modificator.
type LockerOption func(*Locker)

// WithRetryCount set custom retry count.
func WithRetryCount(n int) LockerOption {
	return func(l *Locker) {
		l.retryCount = n
	}
}

// WithRetryTimeout set custom retry timeout.
func WithRetryTimeout(t time.Duration) LockerOption {
	return func(l *Locker) {
		l.retryTimeout = t
	}
}
