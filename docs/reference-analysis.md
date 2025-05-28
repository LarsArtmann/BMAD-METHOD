# Reference Implementation Analysis: LarsArtmann/CV health.go

## Overview

This document analyzes the reference health endpoint implementation from `LarsArtmann/CV@master/internal/health/health.go` to extract patterns for TypeSpec schema design and template generation.

## Change 1: Core Type Definitions and Constants

### Health Status Types
```go
// HealthStatus represents the overall health status
type HealthStatus string

const (
    StatusHealthy   HealthStatus = "healthy"
    StatusDegraded  HealthStatus = "degraded"
    StatusUnhealthy HealthStatus = "unhealthy"
)

// CheckStatus represents the status of an individual check
type CheckStatus string

const (
    CheckHealthy   CheckStatus = "healthy"
    CheckDegraded  CheckStatus = "degraded"
    CheckUnhealthy CheckStatus = "unhealthy"
    CheckDisabled  CheckStatus = "disabled"
    CheckError     CheckStatus = "error"
)
```

### Key Insights for TypeSpec Mapping:
1. **String-based enums**: Use TypeSpec union types for status values
2. **Hierarchical status**: Overall health vs individual check status
3. **Extensible status**: Support for disabled and error states
4. **Clear naming**: Consistent naming patterns for status constants

### TypeSpec Schema Implications:
- Create union types for status enumerations
- Separate models for overall health vs individual checks
- Include all status states including error and disabled
- Use descriptive field names matching Go conventions

## Change 2: Core Data Structures

### HealthCheck Structure
```go
type HealthCheck struct {
    Name        string      `json:"name"`
    Status      CheckStatus `json:"status"`
    Message     string      `json:"message,omitempty"`
    Duration    string      `json:"duration,omitempty"`
    LastChecked time.Time   `json:"last_checked"`
    Details     interface{} `json:"details,omitempty"`
}
```

### HealthReport Structure
```go
type HealthReport struct {
    Status           HealthStatus                   `json:"status"`
    Timestamp        time.Time                      `json:"timestamp"`
    Version          string                         `json:"version"`
    Uptime           time.Duration                  `json:"uptime"`
    UptimeHuman      string                         `json:"uptime_human"`
    Checks           map[string]*HealthCheck        `json:"checks"`
    SystemInfo       *SystemInfo                    `json:"system_info,omitempty"`
    Performance      *PerformanceMetrics            `json:"performance,omitempty"`
    ServerTime       *ServerTimeInfo                `json:"server_time,omitempty"`
    Degradation      *DegradationInfo               `json:"degradation,omitempty"`
    CircuitBreakers  map[string]interface{}         `json:"circuit_breakers,omitempty"`
}
```

### ServerTimeInfo Structure (Key for ServerTime API)
```go
type ServerTimeInfo struct {
    Timezone     string    `json:"timezone"`
    UTCTime      time.Time `json:"utc_time"`
    LocalTime    time.Time `json:"local_time"`
    UTCOffset    string    `json:"utc_offset"`
    Epoch        int64     `json:"epoch"`
    ISO8601      string    `json:"iso8601"`
    RFC3339Nano  string    `json:"rfc3339_nano"`
}
```

### Key Insights for TypeSpec:
1. **Comprehensive timestamps**: Multiple format support (epoch, ISO8601, RFC3339)
2. **Flexible details**: Interface{} for extensible check details
3. **Performance metrics**: Dedicated structures for system monitoring
4. **Circuit breakers**: Integration with resilience patterns
5. **Human-readable formats**: Uptime and duration formatting

## Change 3: Health Check Implementation Patterns

### HealthChecker Core Pattern
```go
type HealthChecker struct {
    startTime time.Time
    version   string
    checks    map[string]func(context.Context) *HealthCheck
}

func NewHealthChecker(version string) *HealthChecker {
    hc := &HealthChecker{
        startTime: time.Now(),
        version:   version,
        checks:    make(map[string]func(context.Context) *HealthCheck),
    }

    // Register default checks
    hc.RegisterCheck("filesystem", hc.checkFilesystem)
    hc.RegisterCheck("memory", hc.checkMemory)
    hc.RegisterCheck("disk", hc.checkDisk)
    hc.RegisterCheck("database", hc.checkDatabase)
    hc.RegisterCheck("degradation", hc.checkDegradation)
    hc.RegisterCheck("circuit_breakers", hc.checkCircuitBreakers)

    return hc
}
```

### Check Execution Pattern
```go
func (hc *HealthChecker) Check(ctx context.Context) *HealthReport {
    uptime := time.Since(hc.startTime)
    report := &HealthReport{
        Timestamp:    time.Now(),
        Version:      hc.version,
        Uptime:       uptime,
        UptimeHuman:  formatUptimeDuration(uptime),
        Checks:       make(map[string]*HealthCheck),
        SystemInfo:   hc.getSystemInfo(),
        Performance:  hc.getPerformanceMetrics(),
    }

    // Run all checks with timing
    overallStatus := StatusHealthy
    for name, checkFunc := range hc.checks {
        result := utils.MeasureOperationWithContext(ctx, func(ctx context.Context) (*HealthCheck, error) {
            return checkFunc(ctx), nil
        })

        check := result.Result
        check.Duration = utils.FormatDuration(result.Duration)
        check.LastChecked = time.Now()

        report.Checks[name] = check

        // Determine overall status
        switch check.Status {
        case CheckUnhealthy, CheckError:
            overallStatus = StatusUnhealthy
        case CheckDegraded:
            if overallStatus == StatusHealthy {
                overallStatus = StatusDegraded
            }
        }
    }

    report.Status = overallStatus
    return report
}
```

### Individual Check Examples
```go
// Filesystem check
func (hc *HealthChecker) checkFilesystem(ctx context.Context) *HealthCheck {
    check := &HealthCheck{Name: "filesystem"}

    tempFile := "/tmp/health_check_test"
    if err := os.WriteFile(tempFile, []byte("test"), 0644); err != nil {
        check.Status = CheckUnhealthy
        check.Message = fmt.Sprintf("Cannot write to filesystem: %v", err)
        return check
    }

    os.Remove(tempFile)
    check.Status = CheckHealthy
    check.Message = "Filesystem is accessible"
    return check
}

// Memory check with thresholds
func (hc *HealthChecker) checkMemory(ctx context.Context) *HealthCheck {
    check := &HealthCheck{Name: "memory"}

    var m runtime.MemStats
    runtime.ReadMemStats(&m)

    allocMB := float64(m.Alloc) / 1024 / 1024

    check.Details = map[string]interface{}{
        "alloc_mb": allocMB,
        "sys_mb":   float64(m.Sys) / 1024 / 1024,
        "num_gc":   m.NumGC,
    }

    if allocMB > 1000 { // 1GB threshold
        check.Status = CheckUnhealthy
        check.Message = fmt.Sprintf("Critical memory usage: %.2f MB allocated", allocMB)
    } else if allocMB > 500 { // 500MB threshold
        check.Status = CheckDegraded
        check.Message = fmt.Sprintf("High memory usage: %.2f MB allocated", allocMB)
    } else {
        check.Status = CheckHealthy
        check.Message = fmt.Sprintf("Memory usage normal: %.2f MB allocated", allocMB)
    }

    return check
}
```

### Key Patterns for Template Generation:
1. **Pluggable architecture**: Function-based check registration
2. **Context-aware**: All checks accept context for cancellation
3. **Timing integration**: Automatic duration measurement
4. **Threshold-based logic**: Configurable health thresholds
5. **Detailed reporting**: Rich details in check responses
6. **Status aggregation**: Smart overall status determination

## Change 4: ServerTime API Implementation

### ServerTime HTTP Handler
```go
func (hc *HealthChecker) ServerTimeHandler(w http.ResponseWriter, r *http.Request) {
    middleware.SetJSONContentType(w)

    now := time.Now()
    location := now.Location()

    serverTime := &ServerTimeInfo{
        Timezone:    location.String(),
        UTCTime:     now.UTC(),
        LocalTime:   now,
        UTCOffset:   now.Format("-07:00"),
        Epoch:       now.Unix(),
        ISO8601:     now.Format(time.RFC3339),
        RFC3339Nano: now.Format(time.RFC3339Nano),
    }

    // Add additional formatted timestamps
    response := map[string]interface{}{
        "server_time": serverTime,
        "formatted": map[string]string{
            "human":     now.Format("Monday, January 2, 2006 at 3:04:05 PM MST"),
            "date":      now.Format("2006-01-02"),
            "time":      now.Format("15:04:05"),
            "datetime":  now.Format("2006-01-02 15:04:05"),
            "rfc822":    now.Format(time.RFC822),
            "rfc850":    now.Format(time.RFC850),
            "rfc1123":   now.Format(time.RFC1123),
        },
        "unix_timestamps": map[string]int64{
            "seconds":      now.Unix(),
            "milliseconds": now.UnixMilli(),
            "microseconds": now.UnixMicro(),
            "nanoseconds":  now.UnixNano(),
        },
    }

    if err := json.NewEncoder(w).Encode(response); err != nil {
        logging.Logger.Error("Failed to encode server time response", "error", err)
        http.Error(w, "Internal server error", http.StatusInternalServerError)
    }
}
```

### Time Formatting Utilities
```go
// formatUptimeDuration formats a duration into human-readable format
func formatUptimeDuration(d time.Duration) string {
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

// Additional time utilities
func getTimezoneInfo() map[string]interface{} {
    now := time.Now()
    _, offset := now.Zone()

    return map[string]interface{}{
        "name":           now.Location().String(),
        "abbreviation":   now.Format("MST"),
        "offset_seconds": offset,
        "offset_hours":   float64(offset) / 3600,
        "is_dst":         now.IsDST(),
    }
}
```

### Performance Timing Integration
```go
// Server timing middleware for performance metrics
func (hc *HealthChecker) withServerTiming(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()

        // Capture various timing metrics
        timings := make(map[string]time.Duration)

        // Database timing (if applicable)
        if dbStart := time.Now(); hc.hasDatabase() {
            // Simulate database ping
            time.Sleep(time.Millisecond * 2)
            timings["db"] = time.Since(dbStart)
        }

        // Cache timing (if applicable)
        if cacheStart := time.Now(); hc.hasCache() {
            // Simulate cache lookup
            time.Sleep(time.Millisecond * 1)
            timings["cache"] = time.Since(cacheStart)
        }

        // Execute the handler
        next(w, r)

        // Add Server-Timing header
        totalDuration := time.Since(start)
        timings["total"] = totalDuration

        var serverTimingValues []string
        for name, duration := range timings {
            ms := float64(duration.Nanoseconds()) / 1e6
            serverTimingValues = append(serverTimingValues,
                fmt.Sprintf(`%s;dur=%.1f;desc="%s"`, name, ms, name))
        }

        if len(serverTimingValues) > 0 {
            w.Header().Set("Server-Timing", strings.Join(serverTimingValues, ", "))
        }
    }
}
```

### Key ServerTime API Insights:
1. **Comprehensive timestamp formats**: RFC3339, ISO8601, Unix variants, human-readable
2. **Timezone awareness**: Full timezone information including DST detection
3. **Multiple precision levels**: Seconds, milliseconds, microseconds, nanoseconds
4. **Server Timing integration**: Performance metrics in HTTP headers
5. **Extensible formatting**: Support for various date/time format standards
6. **Error handling**: Graceful degradation with proper HTTP status codes

## Change 5: TypeSpec Mapping Strategy and Recommendations

### Go to TypeSpec Type Mappings

| Go Type | TypeSpec Type | Notes |
|---------|---------------|-------|
| `string` | `string` | Direct mapping |
| `time.Time` | `utcDateTime` | Use TypeSpec's built-in datetime type |
| `time.Duration` | `duration` | Use TypeSpec's duration type |
| `int64` | `int64` | For Unix timestamps |
| `float64` | `float64` | For memory metrics, percentages |
| `interface{}` | `unknown` or specific union | Use union types for known variants |
| `map[string]*HealthCheck` | `Record<HealthCheck>` | TypeSpec Record type |
| `[]string` | `string[]` | Array type |

### Core TypeSpec Schema Design

```typescript
// health.tsp - Core health models
import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";

using TypeSpec.Http;
using TypeSpec.Rest;

@service({
  title: "Health API",
  version: "1.0.0",
})
namespace HealthAPI;

// Status enumerations
alias HealthStatus = "healthy" | "degraded" | "unhealthy";
alias CheckStatus = "healthy" | "degraded" | "unhealthy" | "disabled" | "error";

// Core models
model HealthCheck {
  /** Check name identifier */
  name: string;

  /** Current status of the check */
  status: CheckStatus;

  /** Optional status message */
  message?: string;

  /** Duration of the check execution */
  duration?: string;

  /** Timestamp of last check */
  lastChecked: utcDateTime;

  /** Additional check-specific details */
  details?: unknown;
}

model HealthReport {
  /** Overall health status */
  status: HealthStatus;

  /** Report timestamp */
  timestamp: utcDateTime;

  /** Service version */
  version: string;

  /** Service uptime duration */
  uptime: duration;

  /** Human-readable uptime */
  uptimeHuman: string;

  /** Individual health checks */
  checks: Record<HealthCheck>;

  /** System information (optional) */
  systemInfo?: SystemInfo;

  /** Performance metrics (optional) */
  performance?: PerformanceMetrics;

  /** Server time information (optional) */
  serverTime?: ServerTimeInfo;
}

model ServerTimeInfo {
  /** Server timezone */
  timezone: string;

  /** UTC timestamp */
  utcTime: utcDateTime;

  /** Local timestamp */
  localTime: utcDateTime;

  /** UTC offset string */
  utcOffset: string;

  /** Unix timestamp (seconds) */
  epoch: int64;

  /** ISO8601 formatted timestamp */
  iso8601: string;

  /** RFC3339 nano precision timestamp */
  rfc3339Nano: string;
}

model ServerTime {
  /** Core server time information */
  serverTime: ServerTimeInfo;

  /** Various formatted timestamps */
  formatted: FormattedTimestamps;

  /** Unix timestamps in different precisions */
  unixTimestamps: UnixTimestamps;
}

model FormattedTimestamps {
  /** Human-readable format */
  human: string;

  /** Date only (YYYY-MM-DD) */
  date: string;

  /** Time only (HH:MM:SS) */
  time: string;

  /** DateTime (YYYY-MM-DD HH:MM:SS) */
  datetime: string;

  /** RFC822 format */
  rfc822: string;

  /** RFC850 format */
  rfc850: string;

  /** RFC1123 format */
  rfc1123: string;
}

model UnixTimestamps {
  /** Unix seconds */
  seconds: int64;

  /** Unix milliseconds */
  milliseconds: int64;

  /** Unix microseconds */
  microseconds: int64;

  /** Unix nanoseconds */
  nanoseconds: int64;
}
```

### HTTP Interface Design

```typescript
// health-api.tsp - HTTP interface definitions
@route("/health")
interface HealthAPI {
  /** Basic health check endpoint */
  @get
  checkHealth(): HealthReport;

  /** Server time endpoint */
  @get
  @route("/time")
  serverTime(): ServerTime;

  /** Readiness probe endpoint */
  @get
  @route("/ready")
  readinessCheck(): HealthReport;

  /** Liveness probe endpoint */
  @get
  @route("/live")
  livenessCheck(): HealthReport;

  /** Individual dependency checks */
  @get
  @route("/dependencies")
  dependencyChecks(): Record<HealthCheck>;
}
```

### Template Tier Progression Strategy

#### Basic Tier TypeSpec
- Core HealthReport model (status, timestamp, version, uptime)
- ServerTime model with essential formats
- Basic health and time endpoints
- No dependency checks or advanced features

#### Intermediate Tier Additions
- HealthCheck model for dependency monitoring
- SystemInfo and PerformanceMetrics models
- Dependency checks endpoint
- Error handling and status aggregation

#### Advanced Tier Extensions
- OpenTelemetry integration models (trace IDs, span context)
- Server Timing metrics models
- CloudEvents schema integration
- Circuit breaker status models

#### Enterprise Tier Features
- Compliance and audit models
- Authentication and authorization schemas
- Advanced monitoring and alerting models
- Multi-environment configuration schemas

### Code Generation Strategy

1. **Go Code Generation**:
   - Generate structs with proper JSON tags
   - Create HTTP handlers with middleware integration
   - Include validation and error handling
   - Add OpenTelemetry instrumentation hooks

2. **TypeScript Generation**:
   - Generate interfaces for all models
   - Create client SDK with proper error handling
   - Include type guards and validation functions
   - Add OpenAPI client generation support

3. **Kubernetes Integration**:
   - Generate health probe configurations
   - Create ServiceMonitor for Prometheus
   - Include Ingress routing configurations
   - Add deployment and service manifests

### Best Practices and Recommendations

1. **Schema Design**:
   - Use union types for enumerations (status values)
   - Make optional fields truly optional with `?`
   - Include comprehensive documentation with `/** */`
   - Use semantic versioning for schema evolution

2. **API Design**:
   - Follow RESTful conventions for endpoint naming
   - Use appropriate HTTP status codes
   - Include proper error response models
   - Support content negotiation (JSON, XML if needed)

3. **Observability Integration**:
   - Include trace ID correlation in all responses
   - Add Server-Timing headers for performance metrics
   - Support CloudEvents for health status changes
   - Integrate with OpenTelemetry for comprehensive monitoring

4. **Template Generation**:
   - Use progressive complexity across tiers
   - Maintain backward compatibility between tiers
   - Include comprehensive testing for generated code
   - Provide clear migration paths between tiers

## Analysis Complete! ✅

### Summary of Key Findings:

1. **Reference Implementation Strengths**:
   - Comprehensive health checking with pluggable architecture
   - Rich ServerTime API with multiple timestamp formats
   - Excellent error handling and status aggregation
   - Performance timing integration with Server-Timing API
   - Extensible design with detailed reporting

2. **TypeSpec Schema Strategy**:
   - Direct mapping from Go types to TypeSpec models
   - Progressive complexity across four template tiers
   - Comprehensive timestamp and timezone support
   - Integration hooks for observability features

3. **Template Generation Approach**:
   - Build upon proven patterns from reference implementation
   - Maintain consistency across Go and TypeScript generation
   - Include Kubernetes integration from the start
   - Support enterprise features with compliance considerations

**Ready for Story 1.2: Design Core TypeSpec Health Schemas** 🚀
