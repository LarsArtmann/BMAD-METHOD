import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";
import "../health/health.tsp";
import "../health/server-time.tsp";
import "../health/health-api.tsp";
import "../health/cloudevents.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;
using TypeSpec.OpenAPI;
using HealthAPI;
using HealthEvents;

/**
 * Advanced tier health API - full observability
 * Deployment time: ~30 minutes
 * 
 * Features:
 * - All Intermediate tier features
 * - Full OpenTelemetry observability
 * - HTTP/3 support
 * - Server Timing API
 * - CloudEvents emission
 * - Custom metrics endpoint (/health/metrics)
 * - Config endpoint (/health/config)
 * - WebSocket support
 * - gRPC support
 * - Circuit breaker patterns
 * - Retry mechanisms
 * - Message queue integration
 * - Health check webhooks
 * - Custom health check plugins
 */
namespace AdvancedTier;

/**
 * WebSocket event types for real-time health monitoring
 */
enum WebSocketEventType {
  healthUpdate: "health.update",
  dependencyChange: "dependency.change",
  metricUpdate: "metric.update",
  alertTriggered: "alert.triggered",
  configChange: "config.change",
}

/**
 * WebSocket health event
 */
model WebSocketHealthEvent {
  /** Event type */
  type: WebSocketEventType;
  
  /** Event timestamp */
  timestamp: utcDateTime;
  
  /** Event payload */
  data: unknown;
  
  /** Event ID for deduplication */
  eventId: string;
  
  /** Correlation ID for tracing */
  correlationId?: string;
}

/**
 * Circuit breaker state
 */
enum CircuitBreakerState {
  closed: "closed",
  open: "open",
  halfOpen: "half-open",
}

/**
 * Circuit breaker status
 */
model CircuitBreakerStatus {
  /** Current state */
  state: CircuitBreakerState;
  
  /** Number of consecutive failures */
  failureCount: int32;
  
  /** Success rate percentage */
  successRate: float64;
  
  /** Last state change timestamp */
  lastStateChange: utcDateTime;
  
  /** Next retry allowed timestamp (if open) */
  nextRetryAt?: utcDateTime;
}

/**
 * Advanced system metrics
 */
model AdvancedMetrics {
  /** System information */
  systemInfo: HealthAPI.SystemInfo;
  
  /** Performance metrics */
  performance: HealthAPI.PerformanceMetrics;
  
  /** Circuit breaker statuses by dependency */
  circuitBreakers?: Record<CircuitBreakerStatus>;
  
  /** Message queue statistics */
  messageQueueStats?: {
    /** Messages in queue */
    queueDepth: int32;
    
    /** Messages processed per second */
    throughput: float64;
    
    /** Average processing time */
    avgProcessingTime: string;
    
    /** Dead letter queue size */
    dlqSize: int32;
  };
  
  /** WebSocket connections */
  websocketStats?: {
    /** Active connections */
    activeConnections: int32;
    
    /** Messages sent per second */
    messageRate: float64;
    
    /** Average latency */
    avgLatency: string;
  };
  
  /** Custom metrics from plugins */
  customMetrics?: Record<unknown>;
  
  /** Timestamp */
  timestamp: utcDateTime;
  
  /** Trace ID */
  traceId?: string;
}

/**
 * Health check configuration
 */
model HealthCheckConfig {
  /** Check interval in seconds */
  interval: int32;
  
  /** Timeout in seconds */
  timeout: int32;
  
  /** Number of retries */
  retries: int32;
  
  /** Circuit breaker threshold */
  failureThreshold: int32;
  
  /** Recovery time in seconds */
  recoveryTime: int32;
  
  /** Enabled checks */
  enabledChecks: string[];
  
  /** Webhook URLs for notifications */
  webhooks?: string[];
  
  /** Custom plugin configurations */
  plugins?: Record<unknown>;
}

/**
 * Health configuration response
 */
model HealthConfigResponse {
  /** Current configuration */
  config: HealthCheckConfig;
  
  /** Configuration version */
  version: string;
  
  /** Last modified timestamp */
  lastModified: utcDateTime;
  
  /** Applied overrides */
  overrides?: Record<unknown>;
}

/**
 * Advanced tier health API endpoints
 */
@route("/health")
@doc("Advanced health endpoints with full observability")
interface AdvancedHealthAPI extends IntermediateTier.IntermediateHealthAPI {
  /**
   * Get detailed system metrics
   */
  @get
  @route("/metrics")
  @doc("Returns detailed system and application metrics")
  
  getMetrics(): {
    @statusCode statusCode: 200;
    @header("Server-Timing") serverTiming?: string;
    ...IntermediateTier.RateLimitHeaders;
    @body body: AdvancedMetrics;
  } | {
    @statusCode statusCode: 503;
    ...IntermediateTier.RateLimitHeaders;
    @body body: HealthError;
  };
  
  /**
   * Get health check configuration
   */
  @get
  @route("/config")
  @doc("Returns current health check configuration")
  
  getConfig(): {
    @statusCode statusCode: 200;
    @header("Server-Timing") serverTiming?: string;
    ...IntermediateTier.RateLimitHeaders;
    @body body: HealthConfigResponse;
  } | {
    @statusCode statusCode: 403;
    @body body: HealthError;
  };
  
  /**
   * Update health check configuration
   */
  @put
  @route("/config")
  @doc("Updates health check configuration")
  
  updateConfig(@body config: HealthCheckConfig): {
    @statusCode statusCode: 200;
    @header("Server-Timing") serverTiming?: string;
    ...IntermediateTier.RateLimitHeaders;
    @body body: HealthConfigResponse;
  } | {
    @statusCode statusCode: 400;
    @body body: HealthError;
  } | {
    @statusCode statusCode: 403;
    @body body: HealthError;
  };
  
  /**
   * WebSocket endpoint for real-time health monitoring
   */
  @route("/ws")
  @doc("WebSocket connection for real-time health updates")
  websocket(): void;
  
  /**
   * Trigger health check webhook test
   */
  @post
  @route("/webhook/test")
  @doc("Tests webhook configuration by sending test event")
  
  testWebhook(@body webhookUrl: { url: string }): {
    @statusCode statusCode: 200;
    @body body: { success: boolean; message: string };
  } | {
    @statusCode statusCode: 400;
    @body body: HealthError;
  };
}

/**
 * Health status event data for CloudEvents
 */
model HealthStatusEventData {
  /** Health status that triggered the event */
  status: HealthStatus;
  
  /** Component that changed status */
  component?: string;
  
  /** Previous status */
  previousStatus?: HealthStatus;
  
  /** Health report snapshot */
  healthReport?: HealthAPI.HealthReport;
  
  /** Timestamp */
  timestamp: utcDateTime;
  
  /** Trace ID */
  traceId?: string;
}

/**
 * Health event for CloudEvents emission
 */
model HealthCloudEvent {
  /** CloudEvents specification version (always "1.0") */
  specversion: "1.0";
  
  /** Event type for health-specific events */
  type: "com.health.status.changed" | "com.health.check.failed" | "com.health.check.recovered";
  
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
  
  /** Health-specific event data */
  data: HealthStatusEventData;
}

/**
 * gRPC service definition for health checks
 * Note: This is a TypeSpec representation - actual gRPC implementation
 * will be generated from protobuf definitions
 */
@doc("gRPC health service interface")
interface HealthGrpcService {
  @doc("Check health via gRPC")
  Check(request: { service?: string }): {
    status: HealthStatus;
    timestamp: utcDateTime;
  };
  
  @doc("Watch health changes via gRPC streaming")
  Watch(request: { service?: string }): {
    status: HealthStatus;
    timestamp: utcDateTime;
  };
}