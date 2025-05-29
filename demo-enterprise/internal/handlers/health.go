package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/demo/enterprise/internal/config"
	"github.com/demo/enterprise/internal/models"
)

// HealthHandler handles health-related HTTP requests
type HealthHandler struct {
	config    *config.Config
	startTime time.Time
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(cfg *config.Config) *HealthHandler {
	return &HealthHandler{
		config:    cfg,
		startTime: time.Now(),
	}
}

// CheckHealth handles GET /health requests
func (h *HealthHandler) CheckHealth(w http.ResponseWriter, r *http.Request) {
	h.setJSONContentType(w)

	uptime := time.Since(h.startTime)
	status := models.HealthReport{
		Status:      "healthy",
		Timestamp:   time.Now(),
		Version:     h.config.Version,
		Uptime:      uptime,
		UptimeHuman: h.formatUptime(uptime),
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(status)
}

// ServerTime handles GET /health/time requests
func (h *HealthHandler) ServerTime(w http.ResponseWriter, r *http.Request) {
	h.setJSONContentType(w)

	now := time.Now()
	location := now.Location()

	serverTime := models.ServerTime{
		Timestamp:   now,
		Timezone:    location.String(),
		Unix:        now.Unix(),
		UnixMilli:   now.UnixMilli(),
		ISO8601:     now.Format(time.RFC3339),
		Formatted:   now.Format("Monday, January 2, 2006 at 3:04:05 PM MST"),
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(serverTime)
}

// ReadinessCheck handles GET /health/ready requests
func (h *HealthHandler) ReadinessCheck(w http.ResponseWriter, r *http.Request) {
	h.setJSONContentType(w)

	// For basic tier, readiness is same as health
	uptime := time.Since(h.startTime)
	status := models.HealthReport{
		Status:      "healthy",
		Timestamp:   time.Now(),
		Version:     h.config.Version,
		Uptime:      uptime,
		UptimeHuman: h.formatUptime(uptime),
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(status)
}

// LivenessCheck handles GET /health/live requests
func (h *HealthHandler) LivenessCheck(w http.ResponseWriter, r *http.Request) {
	h.setJSONContentType(w)

	// For basic tier, liveness is same as health
	uptime := time.Since(h.startTime)
	status := models.HealthReport{
		Status:      "healthy",
		Timestamp:   time.Now(),
		Version:     h.config.Version,
		Uptime:      uptime,
		UptimeHuman: h.formatUptime(uptime),
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(status)
}

// StartupCheck handles GET /health/startup requests
func (h *HealthHandler) StartupCheck(w http.ResponseWriter, r *http.Request) {
	h.setJSONContentType(w)

	// For basic tier, startup is same as health
	uptime := time.Since(h.startTime)
	status := models.HealthReport{
		Status:      "healthy",
		Timestamp:   time.Now(),
		Version:     h.config.Version,
		Uptime:      uptime,
		UptimeHuman: h.formatUptime(uptime),
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(status)
}

// setJSONContentType sets the JSON content type header
func (h *HealthHandler) setJSONContentType(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

// formatUptime formats a duration into human-readable format
func (h *HealthHandler) formatUptime(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%.1f seconds", d.Seconds())
	}
	if d < time.Hour {
		return fmt.Sprintf("%.1f minutes", d.Minutes())
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%.1f hours", d.Hours())
	}
	days := int(d.Hours() / 24)
	hours := int(d.Hours()) % 24
	return fmt.Sprintf("%d days, %d hours", days, hours)
}
