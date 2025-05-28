import "@typespec/http";
import "@typespec/rest";
import "../health/health.tsp";
import "../health/server-time.tsp";
import "../health/health-api.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;
using HealthAPI;

/**
 * Basic tier health API - minimal health checking functionality
 * Deployment time: ~5 minutes
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
}
