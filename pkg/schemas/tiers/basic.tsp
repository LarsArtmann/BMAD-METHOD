import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";
import "../../../template-health/schemas/health.tsp";
import "../../../template-health/schemas/server-time.tsp";
import "../../../template-health/schemas/health-api.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;
using TypeSpec.OpenAPI;
using HealthAPI;

/**
 * Basic tier health API - minimal health checking functionality
 * Deployment time: ~5 minutes
 * 
 * Features:
 * - Core health endpoints (/health, /time, /ready, /live, /startup)
 * - HTTP/1.1 support
 * - HTTP native streaming
 * - Basic security headers
 * - TLS/HTTPS support
 * - API documentation (OpenAPI)
 * - Docker containerization ready
 * - Kubernetes manifests with health probes
 */
namespace BasicTier;

/**
 * Simplified health report for basic tier
 */
model BasicHealthReport {
  /** Overall health status */
  status: HealthStatus;

  /** Report timestamp */
  timestamp: utcDateTime;

  /** Service version */
  version: string;

  /** Service uptime */
  uptime: duration;

  /** Human-readable uptime */
  uptimeHuman: string;

  /** Optional trace ID for basic observability */
  traceId?: string;
}

/**
 * Simplified server time for basic tier
 */
model BasicServerTime {
  /** Current timestamp */
  timestamp: utcDateTime;

  /** Server timezone */
  timezone: string;

  /** Unix timestamp in seconds */
  unix: int64;

  /** Unix timestamp in milliseconds */
  unixMilli: int64;

  /** ISO8601 formatted timestamp */
  iso8601: string;

  /** Human-readable formatted time */
  formatted: string;

  /** Optional trace ID */
  traceId?: string;
}

/**
 * Basic tier health API endpoints
 */
@route("/health")
interface BasicHealthAPI {
  /**
   * Basic health check endpoint
   */
  @get

  checkHealth(): {
    @statusCode statusCode: 200;
    @body body: BasicHealthReport;
  } | {
    @statusCode statusCode: 503;
    @body body: HealthError;
  };

  /**
   * Server time endpoint
   */
  @get
  @route("/time")

  serverTime(): {
    @statusCode statusCode: 200;
    @body body: BasicServerTime;
  } | {
    @statusCode statusCode: 500;
    @body body: HealthError;
  };

  /**
   * Kubernetes liveness probe
   */
  @get
  @route("/live")

  livenessCheck(): {
    @statusCode statusCode: 200;
    @body body: BasicHealthReport;
  } | {
    @statusCode statusCode: 503;
    @body body: HealthError;
  };

  /**
   * Kubernetes readiness probe
   */
  @get
  @route("/ready")

  readinessCheck(): {
    @statusCode statusCode: 200;
    @body body: BasicHealthReport;
  } | {
    @statusCode statusCode: 503;
    @body body: HealthError;
  };

  /**
   * Kubernetes startup probe
   */
  @get
  @route("/startup")

  startupCheck(): {
    @statusCode statusCode: 200;
    @body body: BasicHealthReport;
  } | {
    @statusCode statusCode: 503;
    @body body: HealthError;
  };
}
