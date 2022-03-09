package mq

import "context"

type (
	// Client define message queue client methods.
	Client interface {
		Publish(ctx context.Context, queue string, msg []byte) error
		Close() error
	}
)
