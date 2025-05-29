package middleware

import (
	"fmt"
	"net/http"
	"time"
)

// ServerTimingMiddleware adds Server-Timing headers to responses
func ServerTimingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		
		// Create a response writer that captures the status code
		wrapped := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}
		
		// Call the next handler
		next.ServeHTTP(wrapped, r)
		
		// Calculate total duration
		duration := time.Since(start)
		
		// Add Server-Timing header
		serverTiming := fmt.Sprintf("total;dur=%.1f", float64(duration.Nanoseconds())/1e6)
		
		// Add additional timing metrics if available
		if processingTime := r.Context().Value("processing_time"); processingTime != nil {
			if pt, ok := processingTime.(time.Duration); ok {
				serverTiming += fmt.Sprintf(", processing;dur=%.1f", float64(pt.Nanoseconds())/1e6)
			}
		}
		
		if dbTime := r.Context().Value("db_time"); dbTime != nil {
			if dt, ok := dbTime.(time.Duration); ok {
				serverTiming += fmt.Sprintf(", db;dur=%.1f", float64(dt.Nanoseconds())/1e6)
			}
		}
		
		w.Header().Set("Server-Timing", serverTiming)
	})
}

// responseWriter wraps http.ResponseWriter to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// TimingContext helps track timing information in request context
type TimingContext struct {
	timings map[string]time.Duration
}

// NewTimingContext creates a new timing context
func NewTimingContext() *TimingContext {
	return &TimingContext{
		timings: make(map[string]time.Duration),
	}
}

// AddTiming adds a timing measurement
func (tc *TimingContext) AddTiming(name string, duration time.Duration) {
	tc.timings[name] = duration
}

// GetTimings returns all timing measurements
func (tc *TimingContext) GetTimings() map[string]time.Duration {
	return tc.timings
}
