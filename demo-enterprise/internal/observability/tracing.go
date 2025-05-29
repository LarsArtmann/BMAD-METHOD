package observability

// Simplified tracing for intermediate tier
type TracingProvider struct{}

func NewTracingProvider(serviceName string) *TracingProvider {
	return &TracingProvider{}
}
