package handlers

import (
	"encoding/json"
	"net/http"
	"runtime"
	"time"
)

// MetricsResponse represents the response for metrics endpoint
type MetricsResponse struct {
	Service     ServiceMetrics     `json:"service"`
	Runtime     RuntimeMetrics     `json:"runtime"`
	Health      HealthMetrics      `json:"health"`
	Timestamp   time.Time          `json:"timestamp"`
}

// ServiceMetrics contains service-level metrics
type ServiceMetrics struct {
	Name            string        `json:"name"`
	Version         string        `json:"version"`
	Uptime          time.Duration `json:"uptime"`
	RequestCount    int64         `json:"request_count"`
	ErrorCount      int64         `json:"error_count"`
	AverageResponse string        `json:"average_response_time"`
}

// RuntimeMetrics contains Go runtime metrics
type RuntimeMetrics struct {
	GoVersion      string `json:"go_version"`
	Goroutines     int    `json:"goroutines"`
	MemoryAlloc    uint64 `json:"memory_alloc_bytes"`
	MemoryTotal    uint64 `json:"memory_total_bytes"`
	MemorySys      uint64 `json:"memory_sys_bytes"`
	GCRuns         uint32 `json:"gc_runs"`
	NextGC         uint64 `json:"next_gc_bytes"`
}

// HealthMetrics contains health check metrics
type HealthMetrics struct {
	LastHealthCheck      time.Time `json:"last_health_check"`
	HealthCheckCount     int64     `json:"health_check_count"`
	DependencyCheckCount int64     `json:"dependency_check_count"`
	FailedChecks         int64     `json:"failed_checks"`
}

var (
	serviceStartTime = time.Now()
	requestCount     int64
	errorCount       int64
	healthCheckCount int64
	depCheckCount    int64
	failedChecks     int64
	lastHealthCheck  time.Time
)

// MetricsCheck handles metrics endpoint requests
func (h *HealthHandler) MetricsCheck(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	
	// Increment request count
	requestCount++
	
	// Get runtime metrics
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	
	response := MetricsResponse{
		Service: ServiceMetrics{
			Name:            "{{.Config.Name}}",
			Version:         "{{.Version}}",
			Uptime:          time.Since(serviceStartTime),
			RequestCount:    requestCount,
			ErrorCount:      errorCount,
			AverageResponse: "< 100ms", // TODO: Calculate actual average
		},
		Runtime: RuntimeMetrics{
			GoVersion:   runtime.Version(),
			Goroutines:  runtime.NumGoroutine(),
			MemoryAlloc: m.Alloc,
			MemoryTotal: m.TotalAlloc,
			MemorySys:   m.Sys,
			GCRuns:      m.NumGC,
			NextGC:      m.NextGC,
		},
		Health: HealthMetrics{
			LastHealthCheck:      lastHealthCheck,
			HealthCheckCount:     healthCheckCount,
			DependencyCheckCount: depCheckCount,
			FailedChecks:         failedChecks,
		},
		Timestamp: time.Now(),
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Response-Time", time.Since(start).String())
	
	json.NewEncoder(w).Encode(response)
}

// IncrementHealthCheckCount increments the health check counter
func IncrementHealthCheckCount() {
	healthCheckCount++
	lastHealthCheck = time.Now()
}

// IncrementDependencyCheckCount increments the dependency check counter
func IncrementDependencyCheckCount() {
	depCheckCount++
}

// IncrementErrorCount increments the error counter
func IncrementErrorCount() {
	errorCount++
	failedChecks++
}
