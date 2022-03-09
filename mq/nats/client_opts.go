package nats

import "time"

// clientOption is constructor modification method.
type clientOption func(*clientOptions)

// WithName setup client name.
func WithName(name string) clientOption {
	return func(co *clientOptions) {
		co.name = name
	}
}

// WithDialTimeout setup connection timeout.
func WithDialTimeout(timeout time.Duration) clientOption {
	return func(co *clientOptions) {
		co.dialTimeout = timeout
	}
}

// WithDrainTimeout setup client drain timeout.
func WithDrainTimeout(timeout time.Duration) clientOption {
	return func(co *clientOptions) {
		co.drainTimeout = timeout
	}
}

// WithMaxReconnectCount set count of retries.
func WithMaxReconnectCount(n int) clientOption {
	return func(co *clientOptions) {
		co.maxReconnection = n
	}
}

// WithReconnectTimeout set interval for reconnection.
func WithReconnectTimeout(timeout time.Duration) clientOption {
	return func(co *clientOptions) {
		co.reconnectWait = timeout
	}
}

// WithRootCAs setup connection certificate.
func WithRootCAs(rootCAs ...string) clientOption {
	return func(co *clientOptions) {
		co.rootCAs = rootCAs[:]
	}
}

// WithCredentials setup nats credentials.
func WithCredentials(token, user, password string) clientOption {
	return func(co *clientOptions) {
		co.token = token
		co.user = user
		co.password = password
	}
}

// WithPingInterval setup ping interval.
func WithPingInterval(interval time.Duration) clientOption {
	return func(co *clientOptions) {
		co.pingInterval = interval
	}
}
