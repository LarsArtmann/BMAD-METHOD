package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"{{.Config.GoModule}}/internal/config"
	"{{.Config.GoModule}}/internal/server"
	"{{.Config.GoModule}}/internal/security"
	"{{.Config.GoModule}}/internal/compliance"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize audit logger
	auditLogger, err := compliance.NewAuditLogger("/var/log/{{.Config.Name}}/audit.log", true)
	if err != nil {
		log.Fatalf("Failed to initialize audit logger: %v", err)
	}
	defer auditLogger.Close()

	// Initialize RBAC policy
	rbacPolicy := security.DefaultRBACPolicy()

	// Create server with enterprise features
	srv, err := server.NewEnterprise(cfg, auditLogger, rbacPolicy)
	if err != nil {
		log.Fatalf("Failed to create enterprise server: %v", err)
	}

	// Start server
	go func() {
		fmt.Printf("🚀 Starting {{.Config.Name}} server on :%d\n", cfg.Port)
		if err := srv.Start(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("🛑 Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	fmt.Println("✅ Server exited")
}
