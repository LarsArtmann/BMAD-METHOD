package compliance

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/demo/enterprise/internal/security"
)

// AuditEvent represents an audit log event
type AuditEvent struct {
	Timestamp   time.Time `json:"timestamp"`
	EventID     string    `json:"event_id"`
	UserID      string    `json:"user_id"`
	Action      string    `json:"action"`
	Resource    string    `json:"resource"`
	Method      string    `json:"method"`
	Path        string    `json:"path"`
	StatusCode  int       `json:"status_code"`
	Duration    int64     `json:"duration_ms"`
	RequestID   string    `json:"request_id"`
	ClientIP    string    `json:"client_ip"`
	UserAgent   string    `json:"user_agent"`
	Success     bool      `json:"success"`
	ErrorMsg    string    `json:"error_message,omitempty"`
}

// AuditLogger handles audit logging
type AuditLogger struct {
	logger   *log.Logger
	file     *os.File
	enabled  bool
}

// NewAuditLogger creates a new audit logger
func NewAuditLogger(logFile string, enabled bool) (*AuditLogger, error) {
	if !enabled {
		return &AuditLogger{enabled: false}, nil
	}

	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open audit log file: %w", err)
	}

	logger := log.New(file, "", 0)

	return &AuditLogger{
		logger:   logger,
		file:     file,
		enabled:  true,
	}, nil
}

// LogEvent logs an audit event
func (a *AuditLogger) LogEvent(event AuditEvent) error {
	if !a.enabled {
		return nil
	}

	eventJSON, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal audit event: %w", err)
	}

	a.logger.Println(string(eventJSON))
	return nil
}

// LogHTTPRequest logs an HTTP request audit event
func (a *AuditLogger) LogHTTPRequest(r *http.Request, statusCode int, duration time.Duration, err error) {
	if !a.enabled {
		return
	}

	userID := security.GetClientIdentity(r.Context())
	if userID == "" {
		userID = "anonymous"
	}

	event := AuditEvent{
		Timestamp:  time.Now().UTC(),
		EventID:    fmt.Sprintf("audit_%d", time.Now().UnixNano()),
		UserID:     userID,
		Action:     "http_request",
		Resource:   r.URL.Path,
		Method:     r.Method,
		Path:       r.URL.Path,
		StatusCode: statusCode,
		Duration:   duration.Milliseconds(),
		RequestID:  r.Header.Get("X-Request-ID"),
		ClientIP:   r.RemoteAddr,
		UserAgent:  r.UserAgent(),
		Success:    statusCode < 400,
	}

	if err != nil {
		event.ErrorMsg = err.Error()
	}

	if logErr := a.LogEvent(event); logErr != nil {
		log.Printf("Failed to log audit event: %v", logErr)
	}
}

// AuditMiddleware creates middleware for HTTP request auditing
func AuditMiddleware(auditLogger *AuditLogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			wrapped := &responseWriter{ResponseWriter: w, statusCode: 200}

			next.ServeHTTP(wrapped, r)

			duration := time.Since(start)
			auditLogger.LogHTTPRequest(r, wrapped.statusCode, duration, nil)
		})
	}
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

// Close closes the audit logger
func (a *AuditLogger) Close() error {
	if a.file != nil {
		return a.file.Close()
	}
	return nil
}
