package handlers

import (
	"encoding/json"
	"net/http"
	"time"
)

// DependenciesHandler handles dependency health checks
type DependenciesHandler struct{}

// NewDependenciesHandler creates a new dependencies handler
func NewDependenciesHandler() *DependenciesHandler {
	return &DependenciesHandler{}
}

// DependencyStatus represents the status of a dependency
type DependencyStatus struct {
	Name      string        `json:"name"`
	Status    string        `json:"status"`
	Latency   time.Duration `json:"latency"`
	Error     string        `json:"error,omitempty"`
	Timestamp time.Time     `json:"timestamp"`
}

// DependenciesResponse represents the dependencies health response
type DependenciesResponse struct {
	Status       string              `json:"status"`
	Dependencies []DependencyStatus  `json:"dependencies"`
	Timestamp    time.Time           `json:"timestamp"`
}

// CheckDependencies handles GET /health/dependencies requests
func (h *DependenciesHandler) CheckDependencies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Simulate dependency checks
	dependencies := []DependencyStatus{
		{
			Name:      "database",
			Status:    "healthy",
			Latency:   5 * time.Millisecond,
			Timestamp: time.Now(),
		},
		{
			Name:      "cache",
			Status:    "healthy",
			Latency:   2 * time.Millisecond,
			Timestamp: time.Now(),
		},
	}

	// Determine overall status
	overallStatus := "healthy"
	for _, dep := range dependencies {
		if dep.Status != "healthy" {
			overallStatus = "degraded"
			break
		}
	}

	response := DependenciesResponse{
		Status:       overallStatus,
		Dependencies: dependencies,
		Timestamp:    time.Now(),
	}

	if overallStatus == "healthy" {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
