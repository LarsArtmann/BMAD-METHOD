import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";
import "../health/health.tsp";
import "../health/server-time.tsp";
import "../health/health-api.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;
using TypeSpec.OpenAPI;
using HealthAPI;

/**
 * Intermediate tier health API - production-ready health checking
 * Deployment time: ~15 minutes
 * 
 * Features:
 * - All Basic tier features
 * - Dependency health checks (/health/dependencies)
 * - HTTP/2 support
 * - Structured logging
 * - Basic OpenTelemetry integration
 * - Rate limiting
 * - API versioning
 * - Caching layer (Redis/Memcached)
 */
namespace IntermediateTier;

/**
 * Dependency status
 */
alias DependencyStatus = "healthy" | "degraded" | "unhealthy" | "unreachable";

/**
 * Individual dependency health check
 */
model DependencyCheck {
  /** Name of the dependency */
  name: string;
  
  /** Type of dependency (database, cache, api, filesystem) */
  type: "database" | "cache" | "api" | "filesystem" | "other";
  
  /** Current status */
  status: DependencyStatus;
  
  /** Response time in milliseconds */
  responseTime?: string;
  
  /** Last successful check timestamp */
  lastCheck: utcDateTime;
  
  /** Error message if unhealthy */
  message?: string;
  
  /** Additional dependency-specific metrics */
  metrics?: Record<unknown>;
}

/**
 * Dependency health report
 */
model DependencyHealthReport {
  /** Overall dependency status */
  status: HealthStatus;
  
  /** Timestamp of the report */
  timestamp: utcDateTime;
  
  /** Map of dependency checks by name */
  dependencies: Record<DependencyCheck>;
  
  /** Total number of healthy dependencies */
  healthyCount: int32;
  
  /** Total number of degraded dependencies */
  degradedCount: int32;
  
  /** Total number of unhealthy dependencies */
  unhealthyCount: int32;
  
  /** OpenTelemetry trace ID */
  traceId?: string;
}

/**
 * Enhanced health report for intermediate tier
 */
model IntermediateHealthReport extends HealthAPI.HealthReport {
  /** Basic OpenTelemetry metrics */
  @doc("Basic OpenTelemetry metrics for observability")
  otelMetrics?: {
    /** Number of spans created */
    spanCount?: int32;
    
    /** Number of metrics exported */
    metricsExported?: int32;
    
    /** Trace sampling rate */
    samplingRate?: float64;
  };
  
  /** Cache statistics */
  @doc("Cache layer statistics if enabled")
  cacheStats?: {
    /** Cache hit rate percentage */
    hitRate?: float64;
    
    /** Number of cache entries */
    entryCount?: int32;
    
    /** Cache memory usage in bytes */
    memoryUsage?: int64;
  };
}

/**
 * Rate limit headers
 */
model RateLimitHeaders {
  @header("X-RateLimit-Limit")
  limit: int32;
  
  @header("X-RateLimit-Remaining")
  remaining: int32;
  
  @header("X-RateLimit-Reset")
  reset: int64;
  
  @header("Retry-After")
  retryAfter?: int32;
}

/**
 * Intermediate tier health API endpoints
 */
@route("/health")
@doc("Health endpoints with dependency checking and rate limiting")
interface IntermediateHealthAPI extends BasicTier.BasicHealthAPI {
  /**
   * Check health of external dependencies
   */
  @get
  @route("/dependencies")
  @doc("Returns health status of all external dependencies")
  
  checkDependencies(): {
    @statusCode statusCode: 200;
    ...RateLimitHeaders;
    @body body: DependencyHealthReport;
  } | {
    @statusCode statusCode: 503;
    ...RateLimitHeaders;
    @body body: HealthError;
  } | {
    @statusCode statusCode: 429;
    ...RateLimitHeaders;
    @body body: HealthError;
  };
  
  /**
   * Enhanced health check with dependency awareness
   */
  @get
  @doc("Enhanced health check that includes dependency status")
  
  checkHealthEnhanced(): {
    @statusCode statusCode: 200;
    ...RateLimitHeaders;
    @body body: IntermediateHealthReport;
  } | {
    @statusCode statusCode: 503;
    ...RateLimitHeaders;
    @body body: HealthError;
  } | {
    @statusCode statusCode: 429;
    ...RateLimitHeaders;
    @body body: HealthError;
  };
}

/**
 * API versioning support
 */
@route("/v{version}/health")
@doc("Versioned health API endpoints")
interface VersionedHealthAPI {
  @get
  checkHealth(@path version: "1" | "2"): {
    @statusCode statusCode: 200;
    @body body: IntermediateHealthReport;
  };
}