import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";
import "./health.tsp";
import "./server-time.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;
using HealthAPI;

/**
 * Error response model for health API endpoints
 */
model HealthError {
  /** Error code identifier */
  code: string;

  /** Human-readable error message */
  message: string;

  /** Additional error details */
  details?: unknown;

  /** Timestamp when the error occurred */
  timestamp: utcDateTime;

  /** Optional trace ID for debugging */
  traceId?: string;
}

/**
 * Dependency health check result
 */
model DependencyStatus {
  /** Dependency name identifier */
  name: string;

  /** Current status of the dependency */
  status: CheckStatus;

  /** Response time in milliseconds */
  responseTime?: float64;

  /** Error message if dependency is unhealthy */
  error?: string;

  /** Last successful check timestamp */
  lastCheck: utcDateTime;

  /** Dependency endpoint or connection string (sanitized) */
  endpoint?: string;

  /** Dependency type (database, cache, api, etc.) */
  type?: string;

  /** Additional dependency-specific metadata */
  metadata?: unknown;
}

/**
 * Health API endpoints for service monitoring
 */
@route("/health")
interface HealthEndpoints {
  /**
   * Basic health check endpoint
   * Returns overall service health status with individual check results
   */
  @get
  checkHealth(): {
    @statusCode statusCode: 200;
    @body body: HealthReport;
  } | {
    @statusCode statusCode: 503;
    @body body: HealthError;
  };

  /**
   * Server time endpoint with comprehensive timestamp formats
   * Useful for time synchronization and debugging
   */
  @get
  @route("/time")

  serverTime(): {
    @statusCode statusCode: 200;
    @body body: ServerTime;
  } | {
    @statusCode statusCode: 500;
    @body body: HealthError;
  };

  /**
   * Kubernetes readiness probe endpoint
   * Indicates if the service is ready to receive traffic
   */
  @get
  @route("/ready")

  readinessCheck(): {
    @statusCode statusCode: 200;
    @body body: HealthReport;
  } | {
    @statusCode statusCode: 503;
    @body body: HealthError;
  };

  /**
   * Kubernetes liveness probe endpoint
   * Indicates if the service is alive and should not be restarted
   */
  @get
  @route("/live")

  livenessCheck(): {
    @statusCode statusCode: 200;
    @body body: HealthReport;
  } | {
    @statusCode statusCode: 503;
    @body body: HealthError;
  };

  /**
   * Individual dependency health checks
   * Returns status of external dependencies (databases, APIs, etc.)
   */
  @get
  @route("/dependencies")

  dependencyChecks(): {
    @statusCode statusCode: 200;
    @body body: Record<DependencyStatus>;
  } | {
    @statusCode statusCode: 503;
    @body body: HealthError;
  };

  /**
   * Startup probe endpoint for Kubernetes
   * Indicates if the service has finished starting up
   */
  @get
  @route("/startup")

  startupCheck(): {
    @statusCode statusCode: 200;
    @body body: HealthReport;
  } | {
    @statusCode statusCode: 503;
    @body body: HealthError;
  };
}

/**
 * Extended health API for advanced monitoring (Intermediate+ tiers)
 */
@route("/health")
interface ExtendedHealthEndpoints {
  /**
   * Detailed system metrics and performance data
   */
  @get
  @route("/metrics")

  getMetrics(): {
    @statusCode statusCode: 200;
    @body body: {
      systemInfo: SystemInfo;
      performance: PerformanceMetrics;
      timestamp: utcDateTime;
      traceId?: string;
    };
  } | {
    @statusCode statusCode: 500;
    @body body: HealthError;
  };

  /**
   * Health check configuration and status
   */
  @get
  @route("/config")

  getHealthConfig(): {
    @statusCode statusCode: 200;
    @body body: {
      checks: Record<{
        enabled: boolean;
        interval: duration;
        timeout: duration;
        threshold: {
          warning?: float64;
          critical?: float64;
        };
      }>;
      version: string;
      timestamp: utcDateTime;
    };
  } | {
    @statusCode statusCode: 500;
    @body body: HealthError;
  };
}
