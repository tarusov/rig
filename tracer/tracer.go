// Package tracer contains opentracing rig for jaeger.
package tracer

import (
	"io"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

type (
	// Tracer struct.
	Tracer struct {
		opentracing.Tracer
		io.Closer
	}

	// tracerOptions is aux constructor struct.
	tracerOptions struct {
		agentHostPort string
		serviceName   string
		samplerType   string
		samplerParam  float64
		tags          opentracing.Tags
	}

	// TracerCloseFunc is close func.
	TracerCloseFunc func() error
)

// Defaults.
const (
	defaultServiceName  = "unknown-service"
	defaultSamplerType  = jaeger.SamplerTypeRemote
	defaultSamplerParam = 0.0
)

// New creates new jaeger tracer instance.
func New(opts ...tracerOption) (*Tracer, error) {

	var t = &tracerOptions{
		serviceName:  defaultServiceName,
		samplerType:  defaultSamplerType,
		samplerParam: defaultSamplerParam,
	}

	for _, opt := range opts {
		opt(t)
	}

	var jc = config.Configuration{
		ServiceName: t.serviceName,
		Disabled:    false,
		Sampler: &config.SamplerConfig{
			Type:    t.samplerType,
			Param:   t.samplerParam,
			Options: []jaeger.SamplerOption{},
		},
	}

	if len(t.tags) != 0 {
		var tags = make([]opentracing.Tag, 0)
		for k, v := range t.tags {
			tags = append(tags, opentracing.Tag{Key: k, Value: v})
		}
		jc.Tags = tags
	}

	if t.agentHostPort != "" {
		jc.Reporter = &config.ReporterConfig{
			LocalAgentHostPort: t.agentHostPort,
		}
	}

	ot, tc, err := jc.NewTracer()
	if err != nil {
		return nil, err
	}

	return &Tracer{
		Tracer: ot,
		Closer: tc,
	}, nil
}
