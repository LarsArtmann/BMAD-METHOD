package compliance

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"{{.ModuleName}}/internal/security"
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
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

// AuditLogger handles audit logging
type AuditLogger struct {
	logger   *log.Logger
	file     *os.File
	enabled  bool
	minLevel AuditLevel
}

// AuditLevel represents the audit logging level
type AuditLevel int

const (
	AuditLevelInfo AuditLevel = iota
	AuditLevelWarn
	AuditLevelError
	AuditLevelCritical
)

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
		minLevel: AuditLevelInfo,
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
		EventID:    generateEventID(),
		UserID:     userID,
		Action:     "http_request",
		Resource:   r.URL.Path,
		Method:     r.Method,
		Path:       r.URL.Path,
		StatusCode: statusCode,
		Duration:   duration.Milliseconds(),
		RequestID:  getRequestID(r),
		ClientIP:   getClientIP(r),
		UserAgent:  r.UserAgent(),
		Success:    statusCode < 400,
		Metadata: map[string]interface{}{
			"query_params": r.URL.RawQuery,
			"content_type": r.Header.Get("Content-Type"),
		},
	}

	if err != nil {
		event.ErrorMsg = err.Error()
	}

	if logErr := a.LogEvent(event); logErr != nil {
		log.Printf("Failed to log audit event: %v", logErr)
	}
}

// LogSecurityEvent logs a security-related audit event
func (a *AuditLogger) LogSecurityEvent(ctx context.Context, action, resource string, success bool, metadata map[string]interface{}) {
	if !a.enabled {
		return
	}

	userID := security.GetClientIdentity(ctx)
	if userID == "" {
		userID = "system"
	}

	event := AuditEvent{
		Timestamp: time.Now().UTC(),
		EventID:   generateEventID(),
		UserID:    userID,
		Action:    action,
		Resource:  resource,
		Success:   success,
		Metadata:  metadata,
	}

	if logErr := a.LogEvent(event); logErr != nil {
		log.Printf("Failed to log security audit event: %v", logErr)
	}
}

// AuditMiddleware creates middleware for HTTP request auditing
func AuditMiddleware(auditLogger *AuditLogger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			
			// Wrap response writer to capture status code
			wrapped := &responseWriter{ResponseWriter: w, statusCode: 200}
			
			// Process request
			next.ServeHTTP(wrapped, r)
			
			// Log the request
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

// ComplianceReport generates a compliance report
type ComplianceReport struct {
	GeneratedAt     time.Time              `json:"generated_at"`
	Period          string                 `json:"period"`
	TotalRequests   int                    `json:"total_requests"`
	SuccessfulReqs  int                    `json:"successful_requests"`
	FailedRequests  int                    `json:"failed_requests"`
	SecurityEvents  int                    `json:"security_events"`
	UserActivity    map[string]int         `json:"user_activity"`
	ResourceAccess  map[string]int         `json:"resource_access"`
	ErrorSummary    map[string]int         `json:"error_summary"`
	Metadata        map[string]interface{} `json:"metadata"`
}

// GenerateComplianceReport generates a compliance report from audit logs
func (a *AuditLogger) GenerateComplianceReport(startTime, endTime time.Time) (*ComplianceReport, error) {
	// This is a simplified implementation
	// In a real system, you would parse the audit log file and generate statistics
	
	report := &ComplianceReport{
		GeneratedAt:    time.Now().UTC(),
		Period:         fmt.Sprintf("%s to %s", startTime.Format(time.RFC3339), endTime.Format(time.RFC3339)),
		UserActivity:   make(map[string]int),
		ResourceAccess: make(map[string]int),
		ErrorSummary:   make(map[string]int),
		Metadata: map[string]interface{}{
			"audit_enabled": a.enabled,
			"log_file":      a.file.Name(),
		},
	}

	return report, nil
}

// Close closes the audit logger
func (a *AuditLogger) Close() error {
	if a.file != nil {
		return a.file.Close()
	}
	return nil
}

// Helper functions

func generateEventID() string {
	return fmt.Sprintf("audit_%d", time.Now().UnixNano())
}

func getRequestID(r *http.Request) string {
	if id := r.Header.Get("X-Request-ID"); id != "" {
		return id
	}
	return generateEventID()
}

func getClientIP(r *http.Request) string {
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return ip
	}
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	return r.RemoteAddr
}
