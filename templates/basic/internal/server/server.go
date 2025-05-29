package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"{{.Config.GoModule}}/internal/config"
	"{{.Config.GoModule}}/internal/handlers"
)

// Server represents the HTTP server
type Server struct {
	config  *config.Config
	server  *http.Server
	handler *handlers.HealthHandler
}

// New creates a new server instance
func New(cfg *config.Config) (*Server, error) {
	// Create health handler
	healthHandler := handlers.NewHealthHandler(cfg)

	// Create router
	router := mux.NewRouter()

	// Health endpoints
	health := router.PathPrefix("/health").Subrouter()
	health.HandleFunc("", healthHandler.CheckHealth).Methods("GET")
	health.HandleFunc("/", healthHandler.CheckHealth).Methods("GET")
	health.HandleFunc("/time", healthHandler.ServerTime).Methods("GET")
	health.HandleFunc("/ready", healthHandler.ReadinessCheck).Methods("GET")
	health.HandleFunc("/live", healthHandler.LivenessCheck).Methods("GET")
	health.HandleFunc("/startup", healthHandler.StartupCheck).Methods("GET")

	// Create HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Server{
		config:  cfg,
		server:  srv,
		handler: healthHandler,
	}, nil
}

// Start starts the HTTP server
func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
