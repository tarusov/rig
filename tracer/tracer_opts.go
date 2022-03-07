package tracer

import "github.com/opentracing/opentracing-go"

// tracerOption is jaeger tracer constructor option.
type tracerOption func(*tracerOptions)

// WithAgentAddress setup agent host port.
func WithAgentAddress(hostPort string) tracerOption {
	return func(jt *tracerOptions) {
		jt.agentHostPort = hostPort
	}
}

// WithServiceName setup service name.
func WithServiceName(name string) tracerOption {
	return func(jt *tracerOptions) {
		jt.serviceName = name
	}
}

// WithSamplerType setup sampler type.
func WithSamplerType(samplerType string) tracerOption {
	return func(jt *tracerOptions) {
		jt.samplerType = samplerType
	}
}

// WithSamplerParam setup sampler param.
func WithSamplerParam(samplerParam float64) tracerOption {
	return func(jt *tracerOptions) {
		jt.samplerParam = samplerParam
	}
}

// WithTags setup tags.
func WithTags(tags opentracing.Tags) tracerOption {
	return func(jt *tracerOptions) {
		jt.tags = tags
	}
}
