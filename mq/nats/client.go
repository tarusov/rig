package nats

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/tarusov/rig/logger"
	"github.com/tarusov/rig/mq"
)

type (
	// Client struct.
	Client struct {
		conn *nats.Conn
		wg   sync.WaitGroup
	}

	// clientOptions is auxilary constructor struct.
	clientOptions struct {
		dialTimeout     time.Duration
		drainTimeout    time.Duration
		maxReconnection int
		name            string
		password        string
		pingInterval    time.Duration
		reconnectWait   time.Duration
		rootCAs         []string
		token           string
		user            string
	}
)

// Defaults.
const (
	defaultDialTimeout     = 3 * time.Second
	defaultDrainTimeout    = 30 * time.Second
	defaultMaxReconnection = 3
	defaultName            = "nats_client"
	defaultPingInterval    = 3 * time.Second
	defaultReconnectWait   = 1 * time.Second
)

// New create new NATS client instance.
func New(servers []string, opts ...clientOption) (mq.Client, error) {

	if len(servers) == 0 {
		return nil, errors.New("NATS servers list is empty")
	}

	var (
		co = &clientOptions{
			dialTimeout:     defaultDialTimeout,
			drainTimeout:    defaultDrainTimeout,
			maxReconnection: defaultMaxReconnection,
			name:            defaultName,
			reconnectWait:   defaultReconnectWait,
			pingInterval:    defaultPingInterval,
		}
		c   = &Client{}
		err error
	)

	for _, opt := range opts {
		opt(co)
	}

	c.wg.Add(1)

	var natsOpts = []nats.Option{
		nats.DrainTimeout(co.drainTimeout),
		nats.MaxReconnects(co.maxReconnection),
		nats.Name(co.name),
		nats.PingInterval(co.pingInterval),
		nats.ReconnectWait(co.reconnectWait),
		nats.Timeout(co.dialTimeout),

		func(opts *nats.Options) error {
			opts.Servers = servers
			return nil
		},

		nats.DisconnectErrHandler(func(_ *nats.Conn, err error) {
			logger.New().Errorf("NATS Disconnected because of err: %v", err)
		}),

		nats.ReconnectHandler(func(_ *nats.Conn) {
			logger.New().Info("NATS Reconnecting...")
		}),

		nats.ClosedHandler(func(_ *nats.Conn) {
			logger.New().Debug("NATS Connection closed...")
			c.wg.Done()
		}),
	}

	if len(co.rootCAs) != 0 {
		natsOpts = append(natsOpts, nats.RootCAs(co.rootCAs...))
	}

	if co.token != "" && co.user != "" && co.password != "" {
		natsOpts = append(natsOpts, func(opts *nats.Options) error {
			opts.Token = co.token
			opts.User = co.user
			opts.Password = co.password
			return nil
		})
	}

	c.conn, err = nats.Connect("", natsOpts...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect NATS: %w", err)
	}

	return c, nil
}

// Publish method implements mq.Client Publish method.
func (c *Client) Publish(_ context.Context, queue string, msg []byte) error {
	return c.conn.Publish(queue, msg)
}

// Close method implements mq.Client Close method.
func (c *Client) Close() error {
	c.conn.Close()
	return nil
}
