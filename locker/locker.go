// Package locker contains redis based mutex locker.
package locker

import (
	"context"
	"errors"
	"time"

	"github.com/bsm/redislock"
	"github.com/go-redis/redis/v8"
	"github.com/tarusov/rig/logger"
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

	// Target key mutex lock struct.
	lock struct {
		*redislock.Lock
		key string
		ctx context.Context
	}
)

// Defaults.
const (
	defaultRetryCount   = 3
	defaultRetryTimeout = 3 * time.Second
)

// New creates new locker instance.
func New(client redis.UniversalClient, opts ...lockerOption) (*Locker, error) {

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

// Aux error types.
var (
	ErrLockNotObtained = errors.New("failed to obtain lock")        // Unable to obtain lock.
	ErrNotLocked       = errors.New("failed to unlock - not exist") // No lock exist.
)

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
		if err == redislock.ErrNotObtained {
			return nil, ErrLockNotObtained
		}
		return nil, err
	}

	logger.FromContext(ctx).WithField("key", key).Debug("lock obtained")

	r := lock{
		Lock: obtained,
		key:  key,
		ctx:  ctx,
	}

	return func() error {
		return r.Unlock()
	}, nil
}

// Unlock method try to release current mutex.
func (l *lock) Unlock() (err error) {
	defer func() {
		logger.FromContext(l.ctx).WithField("key", l.key).WithErr(err).Debug("unlocked")
	}()

	err = l.Lock.Release(l.ctx)
	if err == redislock.ErrLockNotHeld {
		return ErrNotLocked
	}

	return err
}
