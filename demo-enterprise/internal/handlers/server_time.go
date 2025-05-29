package handlers

import (
	"encoding/json"
	"net/http"
	"time"
)

// ServerTimeHandler handles server time requests
type ServerTimeHandler struct{}

// NewServerTimeHandler creates a new server time handler
func NewServerTimeHandler() *ServerTimeHandler {
	return &ServerTimeHandler{}
}

// GetServerTime handles GET /health/time requests
func (h *ServerTimeHandler) GetServerTime(w http.ResponseWriter, r *http.Request) {
	now := time.Now()

	response := map[string]interface{}{
		"timestamp":    now,
		"unix":         now.Unix(),
		"unix_milli":   now.UnixMilli(),
		"rfc3339":      now.Format(time.RFC3339),
		"timezone":     now.Location().String(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
