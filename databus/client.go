package databus

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/tarusov/rig/logger"
)

type (
	// Client struct.
	Client struct {
		conn  *nats.EncodedConn
		wg    sync.WaitGroup
		pubMw []PublishMiddlewareFunc
	}

	// clientOptions is auxilary constructor struct.
	clientOptions struct {
		dialTimeout     time.Duration
		drainTimeout    time.Duration
		encoder         nats.Encoder
		encoderName     string
		maxReconnection int
		name            string
		password        string
		pingInterval    time.Duration
		reconnectWait   time.Duration
		rootCAs         []string
		token           string
		user            string
	}

	// PublishMiddlewareFunc is middleware func type.
	PublishMiddlewareFunc func(ctx context.Context, subject string, v interface{}, fn PublishMiddlewareFunc) error
)

// Defaults.
const (
	defaultName            = "client"
	defaultDialTimeout     = 3 * time.Second
	defaultDrainTimeout    = 30 * time.Second
	defaultMaxReconnection = 3
	defaultReconnectWait   = 1 * time.Second
	defaultPingInterval    = 3 * time.Second
)

// NewClient create new databus client instance.
func NewClient(servers []string, opts ...clientOption) (*Client, error) {

	if len(servers) == 0 {
		return nil, errors.New("NATS servers list is empty")
	}

	var (
		c  = &Client{}
		co = &clientOptions{
			dialTimeout:     defaultDialTimeout,
			drainTimeout:    defaultDrainTimeout,
			maxReconnection: defaultMaxReconnection,
			name:            defaultName,
			reconnectWait:   defaultReconnectWait,
			pingInterval:    defaultPingInterval,
		}
	)

	for _, opt := range opts {
		opt(co)
	}

	c.wg.Add(1)
	conn, err := nats.Connect("", mkOptions(&c.wg, servers, co)...)
	if err != nil {
		return nil, fmt.Errorf("failed to connect NATS: %w", err)
	}

	// Setup encoder if it set, or choose json.
	if co.encoder != nil && co.encoderName != "" {
		nats.RegisterEncoder(co.encoderName, co.encoder)
	} else {
		co.encoderName = nats.JSON_ENCODER
	}

	c.conn, err = nats.NewEncodedConn(conn, co.encoderName)
	if err != nil {
		return nil, fmt.Errorf("failed to create encoded connection: %w", err)
	}

	return c, nil
}

func mkOptions(clientWg *sync.WaitGroup, servers []string, co *clientOptions) []nats.Option {

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
			clientWg.Done()
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

	return natsOpts
}

// Publish
func (c *Client) Publish(ctx context.Context, subject string, v interface{}) error {

	logger.FromContext(ctx).WithField("subject", subject).Debug("sending message")

	return nil
}
