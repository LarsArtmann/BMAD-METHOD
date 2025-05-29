package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	"{{.Config.GoModule}}/internal/config"
	"{{.Config.GoModule}}/internal/handlers"
	"{{.Config.GoModule}}/internal/security"
	"{{.Config.GoModule}}/internal/compliance"
	"{{.Config.GoModule}}/internal/middleware"
)

// Server represents the HTTP server
type Server struct {
	config      *config.Config
	server      *http.Server
	handler     *handlers.HealthHandler
	auditLogger *compliance.AuditLogger
	rbacPolicy  *security.RBACPolicy
}

// EnterpriseServer represents the enterprise HTTP server with security features
type EnterpriseServer struct {
	*Server
	tlsConfig *tls.Config
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
	health.HandleFunc("/dependencies", healthHandler.DependenciesCheck).Methods("GET")
	health.HandleFunc("/metrics", healthHandler.MetricsCheck).Methods("GET")

	// Create HTTP server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Server{
		config:      cfg,
		server:      srv,
		handler:     healthHandler,
		auditLogger: nil,
		rbacPolicy:  nil,
	}, nil
}

// NewEnterprise creates a new enterprise server instance with security features
func NewEnterprise(cfg *config.Config, auditLogger *compliance.AuditLogger, rbacPolicy *security.RBACPolicy) (*EnterpriseServer, error) {
	// Create health handler
	healthHandler := handlers.NewHealthHandler(cfg)

	// Create router
	router := mux.NewRouter()

	// Add enterprise middleware
	router.Use(compliance.AuditMiddleware(auditLogger))
	router.Use(security.RBACMiddleware(rbacPolicy))
	router.Use(middleware.ServerTimingMiddleware)

	// Health endpoints
	health := router.PathPrefix("/health").Subrouter()
	health.HandleFunc("", healthHandler.CheckHealth).Methods("GET")
	health.HandleFunc("/", healthHandler.CheckHealth).Methods("GET")
	health.HandleFunc("/time", healthHandler.ServerTime).Methods("GET")
	health.HandleFunc("/ready", healthHandler.ReadinessCheck).Methods("GET")
	health.HandleFunc("/live", healthHandler.LivenessCheck).Methods("GET")
	health.HandleFunc("/startup", healthHandler.StartupCheck).Methods("GET")
	health.HandleFunc("/dependencies", healthHandler.DependenciesCheck).Methods("GET")
	health.HandleFunc("/metrics", healthHandler.MetricsCheck).Methods("GET")

	// Setup mTLS configuration
	mtlsConfig, err := security.SetupMTLS(security.MTLSConfig{
		CertFile:   "/etc/ssl/certs/server.crt",
		KeyFile:    "/etc/ssl/private/server.key",
		CAFile:     "/etc/ssl/certs/ca.crt",
		ClientAuth: tls.RequireAndVerifyClientCert,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to setup mTLS: %w", err)
	}

	// Add mTLS middleware
	router.Use(security.MTLSMiddleware)

	// Create HTTP server with TLS
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
		TLSConfig:    mtlsConfig,
	}

	baseServer := &Server{
		config:      cfg,
		server:      srv,
		handler:     healthHandler,
		auditLogger: auditLogger,
		rbacPolicy:  rbacPolicy,
	}

	return &EnterpriseServer{
		Server:    baseServer,
		tlsConfig: mtlsConfig,
	}, nil
}

// Start starts the HTTP server
func (s *Server) Start() error {
	return s.server.ListenAndServe()
}

// Start starts the enterprise HTTP server with TLS
func (es *EnterpriseServer) Start() error {
	return es.server.ListenAndServeTLS("", "")
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

// Shutdown gracefully shuts down the enterprise server
func (es *EnterpriseServer) Shutdown(ctx context.Context) error {
	return es.server.Shutdown(ctx)
}
