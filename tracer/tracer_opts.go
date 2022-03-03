package tracer

import "github.com/opentracing/opentracing-go"

// TracerOption is jaeger tracer constructor option.
type TracerOption func(*jaegerTracer)

// WithAgentAddress setup agent host port.
func WithAgentAddress(hostPort string) TracerOption {
	return func(jt *jaegerTracer) {
		jt.agentHostPort = hostPort
	}
}

// WithServiceName setup service name.
func WithServiceName(name string) TracerOption {
	return func(jt *jaegerTracer) {
		jt.serviceName = name
	}
}

// WithSamplerType setup sampler type.
func WithSamplerType(samplerType string) TracerOption {
	return func(jt *jaegerTracer) {
		jt.samplerType = samplerType
	}
}

// WithSamplerParam setup sampler param.
func WithSamplerParam(samplerParam float64) TracerOption {
	return func(jt *jaegerTracer) {
		jt.samplerParam = samplerParam
	}
}

// WithTags setup tags.
func WithTags(tags opentracing.Tags) TracerOption {
	return func(jt *jaegerTracer) {
		jt.tags = tags
	}
}
