// Package locker contains redis based mutex locker.
package locker

import (
	"context"
	"time"

	"github.com/bsm/redislock"
	"github.com/go-redis/redis/v8"
)

type (
	// Locker struct.
	Locker struct {
		client       redis.UniversalClient
		retryCount   int
		retryTimeout time.Duration
	}

	// UnlockFunc is method for unlock func.
	UnlockFunc func() error

	// Target key mutex lock.
	lock struct {
		*redislock.Lock
		ctx context.Context
	}
)

// Defaults.
const (
	defaultRetryCount   = 3
	defaultRetryTimeout = 3 * time.Second
)

// New creates new locker instance.
func New(client redis.UniversalClient, opts ...LockerOption) (*Locker, error) {

	var l = &Locker{
		client:       client,
		retryCount:   defaultRetryCount,
		retryTimeout: defaultRetryTimeout,
	}

	for _, opt := range opts {
		opt(l)
	}

	return l, nil
}

// Lock method create new redis mutex lock. Return unlock func or error.
func (l *Locker) Lock(ctx context.Context, key string, ttl time.Duration) (UnlockFunc, error) {

	obtained, err := redislock.Obtain(
		ctx,
		l.client,
		key,
		ttl,
		&redislock.Options{
			RetryStrategy: redislock.LimitRetry(redislock.LinearBackoff(l.retryTimeout), l.retryCount),
		},
	)
	if err != nil {
		return nil, err
	}

	r := lock{
		Lock: obtained,
		ctx:  ctx,
	}

	return func() error {
		return r.Unlock()
	}, nil
}

// Unlock method try to release current mutex.
func (l *lock) Unlock() error {
	return l.Lock.Release(l.ctx)
}
