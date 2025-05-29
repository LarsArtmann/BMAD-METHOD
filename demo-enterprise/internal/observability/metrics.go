package observability

// Simplified metrics for intermediate tier
type MetricsProvider struct{}

func NewMetricsProvider(serviceName string) *MetricsProvider {
	return &MetricsProvider{}
}
