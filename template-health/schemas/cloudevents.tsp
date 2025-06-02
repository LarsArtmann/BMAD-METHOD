import "@typespec/http";
import "@typespec/rest";
import "./health.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;
using HealthAPI;

/**
 * CloudEvents v1.0 specification compliance for health events
 */
namespace HealthEvents;

/**
 * CloudEvents context attributes for health events
 */
model CloudEventContext {
  /** CloudEvents specification version (always "1.0") */
  specversion: "1.0";
  
  /** Event type identifier */
  type: string;
  
  /** Event source identifier (service name/instance) */
  source: string;
  
  /** Unique event identifier */
  id: string;
  
  /** Event timestamp */
  time: utcDateTime;
  
  /** Content type of the data payload */
  datacontenttype: "application/json";
  
  /** Optional subject for the event */
  subject?: string;
  
  /** Optional schema reference for the data */
  dataschema?: string;
}

/**
 * Health status change event data
 */
model HealthStatusChangeData {
  /** Current health status */
  status: HealthStatus;
  
  /** Previous health status */
  previousStatus: HealthStatus;
  
  /** Service version */
  version: string;
  
  /** Timestamp of the status change */
  timestamp: utcDateTime;
  
  /** Reason for the status change */
  reason?: string;
  
  /** OpenTelemetry trace ID for correlation */
  traceId?: string;
  
  /** Affected health checks */
  affectedChecks?: string[];
  
  /** Service instance identifier */
  instanceId?: string;
  
  /** Additional context about the change */
  metadata?: Record<unknown>;
}

/**
 * Dependency status change event data
 */
model DependencyStatusChangeData {
  /** Dependency name */
  dependencyName: string;
  
  /** Current dependency status */
  status: CheckStatus;
  
  /** Previous dependency status */
  previousStatus: CheckStatus;
  
  /** Timestamp of the status change */
  timestamp: utcDateTime;
  
  /** Error message if dependency failed */
  error?: string;
  
  /** Response time when the change occurred */
  responseTime?: float64;
  
  /** OpenTelemetry trace ID for correlation */
  traceId?: string;
  
  /** Dependency endpoint (sanitized) */
  endpoint?: string;
  
  /** Dependency type */
  type?: string;
}

/**
 * Service startup event data
 */
model ServiceStartupData {
  /** Service version */
  version: string;
  
  /** Startup timestamp */
  timestamp: utcDateTime;
  
  /** Time taken to start up */
  startupDuration: duration;
  
  /** Initial health status */
  initialStatus: HealthStatus;
  
  /** Service instance identifier */
  instanceId?: string;
  
  /** Environment information */
  environment?: string;
  
  /** Configuration hash for change detection */
  configHash?: string;
}

/**
 * Service shutdown event data
 */
model ServiceShutdownData {
  /** Service version */
  version: string;
  
  /** Shutdown timestamp */
  timestamp: utcDateTime;
  
  /** Final health status */
  finalStatus: HealthStatus;
  
  /** Reason for shutdown */
  reason?: string;
  
  /** Service instance identifier */
  instanceId?: string;
  
  /** Graceful shutdown indicator */
  graceful: boolean;
  
  /** Uptime duration */
  uptime: duration;
}

/**
 * Performance degradation event data
 */
model PerformanceDegradationData {
  /** Metric that triggered the degradation */
  metric: string;
  
  /** Current metric value */
  currentValue: float64;
  
  /** Threshold that was exceeded */
  threshold: float64;
  
  /** Severity level */
  severity: "warning" | "critical";
  
  /** Timestamp when degradation was detected */
  timestamp: utcDateTime;
  
  /** Duration of the degradation */
  duration?: duration;
  
  /** OpenTelemetry trace ID for correlation */
  traceId?: string;
  
  /** Additional performance metrics */
  additionalMetrics?: Record<float64>;
}

/**
 * Complete CloudEvents health event
 */
model HealthEvent {
  ...CloudEventContext;
  
  /** Event data payload */
  data: HealthStatusChangeData | DependencyStatusChangeData | ServiceStartupData | ServiceShutdownData | PerformanceDegradationData;
}

/**
 * Predefined event types for health monitoring
 */
namespace EventTypes {
  /** Service health status changed */
  alias HealthStatusChanged = "com.health.status.changed";
  
  /** Dependency status changed */
  alias DependencyStatusChanged = "com.health.dependency.changed";
  
  /** Service started up */
  alias ServiceStartup = "com.health.service.startup";
  
  /** Service shutting down */
  alias ServiceShutdown = "com.health.service.shutdown";
  
  /** Performance degradation detected */
  alias PerformanceDegradation = "com.health.performance.degradation";
  
  /** Health check failed */
  alias HealthCheckFailed = "com.health.check.failed";
  
  /** Health check recovered */
  alias HealthCheckRecovered = "com.health.check.recovered";
}

/**
 * Event emission configuration
 */
model EventConfig {
  /** Whether event emission is enabled */
  enabled: boolean;
  
  /** Event broker endpoint */
  brokerEndpoint?: string;
  
  /** Event topic/subject pattern */
  topicPattern: string;
  
  /** Retry configuration for failed events */
  retryConfig?: {
    maxRetries: int32;
    retryDelay: duration;
    backoffMultiplier: float64;
  };
  
  /** Event filtering rules */
  filters?: {
    includeEvents: string[];
    excludeEvents?: string[];
    minSeverity?: "info" | "warning" | "error";
  };
}
