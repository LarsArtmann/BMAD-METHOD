# Enterprise-Grade Template Architecture

## Prompt Name: Enterprise-Grade Template Architecture

## Context
You need to design and implement enterprise-grade template architecture for generating sophisticated, production-ready applications with security, compliance, observability, and scalability features.

## Enterprise Requirements

### Security Features
- **mTLS (Mutual TLS)**: Client certificate authentication
- **RBAC (Role-Based Access Control)**: Permission-based access control
- **Audit Logging**: Comprehensive security event logging
- **Security Context**: Request-scoped security information
- **Input Validation**: Comprehensive request validation
- **Rate Limiting**: Protection against abuse

### Compliance Features
- **Audit Trails**: Immutable audit logs
- **Data Privacy**: GDPR/CCPA compliance patterns
- **Retention Policies**: Automated data lifecycle management
- **Compliance Reporting**: Automated compliance reports
- **Access Controls**: Fine-grained permission systems

### Observability Features
- **OpenTelemetry**: Distributed tracing and metrics
- **Structured Logging**: JSON-formatted, searchable logs
- **Health Checks**: Comprehensive health monitoring
- **Performance Metrics**: Application and business metrics
- **Alerting**: Proactive issue detection

## Architecture Patterns

### 1. Multi-Tier Template System
```
Enterprise Template Architecture:
├── Basic Tier
│   ├── Core health endpoints
│   ├── Basic logging
│   └── Docker support
├── Intermediate Tier
│   ├── Dependency health checks
│   ├── Server timing metrics
│   └── Basic observability
├── Advanced Tier
│   ├── Full OpenTelemetry integration
│   ├── CloudEvents support
│   └── Kubernetes manifests
└── Enterprise Tier
    ├── mTLS authentication
    ├── RBAC authorization
    ├── Audit logging
    ├── Compliance features
    └── Multi-environment support
```

### 2. Security Architecture
```go
// Enterprise security stack
type SecurityConfig struct {
    // mTLS Configuration
    MTLSEnabled    bool   `yaml:"mtls_enabled"`
    ClientCertPath string `yaml:"client_cert_path"`
    ClientKeyPath  string `yaml:"client_key_path"`
    CACertPath     string `yaml:"ca_cert_path"`
    
    // RBAC Configuration
    RBACEnabled    bool   `yaml:"rbac_enabled"`
    PolicyPath     string `yaml:"policy_path"`
    
    // Audit Configuration
    AuditEnabled   bool   `yaml:"audit_enabled"`
    AuditLogPath   string `yaml:"audit_log_path"`
}

// Security middleware stack
func SecurityMiddleware(config SecurityConfig) []Middleware {
    var middlewares []Middleware
    
    if config.MTLSEnabled {
        middlewares = append(middlewares, MTLSMiddleware(config))
    }
    
    if config.RBACEnabled {
        middlewares = append(middlewares, RBACMiddleware(config))
    }
    
    if config.AuditEnabled {
        middlewares = append(middlewares, AuditMiddleware(config))
    }
    
    return middlewares
}
```

### 3. Observability Architecture
```go
// Enterprise observability stack
type ObservabilityConfig struct {
    // Tracing
    TracingEnabled bool   `yaml:"tracing_enabled"`
    JaegerEndpoint string `yaml:"jaeger_endpoint"`
    
    // Metrics
    MetricsEnabled    bool   `yaml:"metrics_enabled"`
    PrometheusPort    int    `yaml:"prometheus_port"`
    
    // Logging
    LogLevel          string `yaml:"log_level"`
    LogFormat         string `yaml:"log_format"` // json, text
    
    // Health Checks
    HealthCheckInterval time.Duration `yaml:"health_check_interval"`
}

// Observability initialization
func InitObservability(config ObservabilityConfig) error {
    // Initialize tracing
    if config.TracingEnabled {
        if err := initTracing(config.JaegerEndpoint); err != nil {
            return err
        }
    }
    
    // Initialize metrics
    if config.MetricsEnabled {
        if err := initMetrics(config.PrometheusPort); err != nil {
            return err
        }
    }
    
    // Initialize structured logging
    initLogging(config.LogLevel, config.LogFormat)
    
    return nil
}
```

## Template Implementation Patterns

### 1. Progressive Feature Enhancement
```yaml
# Template feature matrix
features:
  basic:
    core: true
    docker: true
    basic_logging: true
  
  intermediate:
    core: true
    docker: true
    basic_logging: true
    dependencies: true
    server_timing: true
    basic_metrics: true
  
  advanced:
    core: true
    docker: true
    basic_logging: true
    dependencies: true
    server_timing: true
    basic_metrics: true
    opentelemetry: true
    cloudevents: true
    kubernetes: true
    typescript: true
  
  enterprise:
    core: true
    docker: true
    basic_logging: true
    dependencies: true
    server_timing: true
    basic_metrics: true
    opentelemetry: true
    cloudevents: true
    kubernetes: true
    typescript: true
    mtls: true
    rbac: true
    audit: true
    compliance: true
    multi_env: true
```

### 2. Configuration Management
```go
// Enterprise configuration structure
type EnterpriseConfig struct {
    // Application settings
    App AppConfig `yaml:"app"`
    
    // Server settings
    Server ServerConfig `yaml:"server"`
    
    // Security settings
    Security SecurityConfig `yaml:"security"`
    
    // Observability settings
    Observability ObservabilityConfig `yaml:"observability"`
    
    // Compliance settings
    Compliance ComplianceConfig `yaml:"compliance"`
    
    // Environment-specific overrides
    Environment string `yaml:"environment"`
}

// Environment-specific configuration loading
func LoadConfig(env string) (*EnterpriseConfig, error) {
    // Load base configuration
    config, err := loadBaseConfig()
    if err != nil {
        return nil, err
    }
    
    // Apply environment-specific overrides
    if err := applyEnvironmentOverrides(config, env); err != nil {
        return nil, err
    }
    
    // Validate configuration
    if err := config.Validate(); err != nil {
        return nil, err
    }
    
    return config, nil
}
```

### 3. Kubernetes Integration
```yaml
# Enterprise Kubernetes manifests
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{.Config.Name}}
  labels:
    app: {{.Config.Name}}
    tier: enterprise
spec:
  replicas: 3
  selector:
    matchLabels:
      app: {{.Config.Name}}
  template:
    metadata:
      labels:
        app: {{.Config.Name}}
      annotations:
        prometheus.io/scrape: "true"
        prometheus.io/port: "8080"
        prometheus.io/path: "/metrics"
    spec:
      containers:
      - name: {{.Config.Name}}
        image: {{.Config.Image}}
        ports:
        - containerPort: 8080
          name: http
        - containerPort: 8443
          name: https
        env:
        - name: ENVIRONMENT
          value: "production"
        - name: LOG_LEVEL
          value: "info"
        livenessProbe:
          httpGet:
            path: /health/live
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health/ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        volumeMounts:
        - name: tls-certs
          mountPath: /etc/tls
          readOnly: true
        - name: config
          mountPath: /etc/config
          readOnly: true
      volumes:
      - name: tls-certs
        secret:
          secretName: {{.Config.Name}}-tls
      - name: config
        configMap:
          name: {{.Config.Name}}-config
```

## Implementation Guidelines

### 1. Security Implementation
```go
// mTLS implementation
func setupMTLS(config SecurityConfig) (*tls.Config, error) {
    // Load client certificate
    cert, err := tls.LoadX509KeyPair(config.ClientCertPath, config.ClientKeyPath)
    if err != nil {
        return nil, err
    }
    
    // Load CA certificate
    caCert, err := os.ReadFile(config.CACertPath)
    if err != nil {
        return nil, err
    }
    
    caCertPool := x509.NewCertPool()
    caCertPool.AppendCertsFromPEM(caCert)
    
    return &tls.Config{
        Certificates: []tls.Certificate{cert},
        ClientCAs:    caCertPool,
        ClientAuth:   tls.RequireAndVerifyClientCert,
    }, nil
}

// RBAC implementation
func (p *RBACPolicy) HasPermission(userID string, permission Permission) bool {
    user, exists := p.Users[userID]
    if !exists {
        return false
    }
    
    for _, role := range user.Roles {
        for _, perm := range role.Permissions {
            if perm == permission {
                return true
            }
        }
    }
    
    return false
}
```

### 2. Audit Logging
```go
// Comprehensive audit logging
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

func (a *AuditLogger) LogHTTPRequest(r *http.Request, statusCode int, duration time.Duration, err error) {
    event := AuditEvent{
        Timestamp:  time.Now().UTC(),
        EventID:    generateEventID(),
        UserID:     getUserID(r.Context()),
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
    
    a.LogEvent(event)
}
```

### 3. Compliance Features
```go
// Data retention and compliance
type ComplianceManager struct {
    retentionPolicies map[string]time.Duration
    encryptionKey     []byte
    auditLogger       *AuditLogger
}

func (c *ComplianceManager) ProcessDataRetention() error {
    for dataType, retention := range c.retentionPolicies {
        cutoff := time.Now().Add(-retention)
        
        if err := c.deleteExpiredData(dataType, cutoff); err != nil {
            return err
        }
        
        c.auditLogger.LogEvent(AuditEvent{
            Action:   "data_retention",
            Resource: dataType,
            Success:  true,
        })
    }
    
    return nil
}
```

## Testing Enterprise Features

### 1. Security Testing
```go
func TestMTLSAuthentication(t *testing.T) {
    // Test valid client certificate
    client := createMTLSClient(validCert)
    resp, err := client.Get("https://localhost:8443/health")
    assert.NoError(t, err)
    assert.Equal(t, 200, resp.StatusCode)
    
    // Test invalid client certificate
    client = createMTLSClient(invalidCert)
    _, err = client.Get("https://localhost:8443/health")
    assert.Error(t, err)
}

func TestRBACAuthorization(t *testing.T) {
    policy := createTestRBACPolicy()
    
    // Test admin user has all permissions
    assert.True(t, policy.HasPermission("admin", PermissionHealthRead))
    assert.True(t, policy.HasPermission("admin", PermissionHealthWrite))
    
    // Test service user has limited permissions
    assert.True(t, policy.HasPermission("service", PermissionHealthRead))
    assert.False(t, policy.HasPermission("service", PermissionHealthWrite))
}
```

### 2. Compliance Testing
```go
func TestAuditLogging(t *testing.T) {
    logger := createTestAuditLogger()
    
    // Simulate HTTP request
    req := createTestRequest()
    logger.LogHTTPRequest(req, 200, time.Millisecond*100, nil)
    
    // Verify audit event was logged
    events := logger.GetEvents()
    assert.Len(t, events, 1)
    assert.Equal(t, "http_request", events[0].Action)
    assert.Equal(t, 200, events[0].StatusCode)
}
```

## Success Criteria

### Security
- ✅ mTLS authentication working
- ✅ RBAC authorization enforced
- ✅ Audit logging comprehensive
- ✅ Security context propagated

### Compliance
- ✅ Audit trails immutable
- ✅ Data retention policies enforced
- ✅ Access controls granular
- ✅ Compliance reports automated

### Observability
- ✅ Distributed tracing enabled
- ✅ Metrics collection comprehensive
- ✅ Structured logging implemented
- ✅ Health checks reliable

### Performance
- ✅ Sub-100ms response times
- ✅ High availability (99.9%+)
- ✅ Horizontal scalability
- ✅ Resource efficiency

This enterprise-grade architecture ensures production-ready, secure, compliant, and observable applications.
