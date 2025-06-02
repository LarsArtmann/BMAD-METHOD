import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";

using TypeSpec.Http;
using TypeSpec.Rest;

namespace HealthAPI;

/**
 * Overall health status of the service
 */
alias HealthStatus = "healthy" | "degraded" | "unhealthy";

/**
 * Status of an individual health check
 */
alias CheckStatus = "healthy" | "degraded" | "unhealthy" | "disabled" | "error";

/**
 * Individual health check result
 */
model HealthCheck {
  /** Unique identifier for the health check */
  name: string;
  
  /** Current status of this check */
  status: CheckStatus;
  
  /** Optional descriptive message about the check result */
  message?: string;
  
  /** Duration taken to execute this check (e.g., "15.2ms") */
  duration?: string;
  
  /** Timestamp when this check was last executed */
  lastChecked: utcDateTime;
  
  /** Additional check-specific details and metrics */
  details?: unknown;
}

/**
 * Comprehensive health report for the service
 */
model HealthReport {
  /** Overall aggregated health status */
  status: HealthStatus;
  
  /** Timestamp when this report was generated */
  timestamp: utcDateTime;
  
  /** Service version identifier */
  version: string;
  
  /** Service uptime duration */
  uptime: duration;
  
  /** Human-readable uptime description */
  uptimeHuman: string;
  
  /** Map of individual health check results */
  checks: Record<HealthCheck>;
  
  /** Optional system information */
  systemInfo?: SystemInfo;
  
  /** Optional performance metrics */
  performance?: PerformanceMetrics;
  
  /** Optional server time information */
  serverTime?: ServerTimeInfo;
  
  /** Optional OpenTelemetry trace ID for correlation */
  traceId?: string;
  
  /** Optional Server-Timing metrics */
  serverTiming?: ServerTimingMetrics;
}

/**
 * System information and resource metrics
 */
model SystemInfo {
  /** Operating system name */
  os: string;
  
  /** System architecture */
  arch: string;
  
  /** Number of CPU cores */
  cpuCores: int32;
  
  /** Total system memory in bytes */
  totalMemory: int64;
  
  /** Available memory in bytes */
  availableMemory: int64;
  
  /** Memory usage percentage */
  memoryUsagePercent: float64;
}

/**
 * Performance metrics for the service
 */
model PerformanceMetrics {
  /** Current memory allocation in bytes */
  allocatedMemory: int64;
  
  /** System memory usage in bytes */
  systemMemory: int64;
  
  /** Number of garbage collection cycles */
  gcCount: int32;
  
  /** Number of active goroutines */
  goroutines: int32;
  
  /** CPU usage percentage */
  cpuUsagePercent: float64;
}

/**
 * Core server time information with timezone details
 */
model ServerTimeInfo {
  /** Server timezone identifier (e.g., "UTC", "America/New_York") */
  timezone: string;
  
  /** Current time in UTC */
  utcTime: utcDateTime;
  
  /** Current time in server's local timezone */
  localTime: utcDateTime;
  
  /** UTC offset string (e.g., "+00:00", "-05:00") */
  utcOffset: string;
  
  /** Unix timestamp in seconds */
  epoch: int64;
  
  /** ISO8601 formatted timestamp */
  iso8601: string;
  
  /** RFC3339 nano precision timestamp */
  rfc3339Nano: string;
  
  /** Timezone abbreviation (e.g., "UTC", "EST") */
  timezoneAbbr: string;
  
  /** Whether daylight saving time is active */
  isDST: boolean;
}

/**
 * Server-Timing API metrics for performance analysis
 */
model ServerTimingMetrics {
  /** Database query time in milliseconds */
  dbQuery?: float64;
  
  /** Cache lookup time in milliseconds */
  cacheHit?: float64;
  
  /** Total processing time in milliseconds */
  total?: float64;
  
  /** External API call time in milliseconds */
  externalApi?: float64;
  
  /** Authentication time in milliseconds */
  auth?: float64;
  
  /** Validation time in milliseconds */
  validation?: float64;
}
