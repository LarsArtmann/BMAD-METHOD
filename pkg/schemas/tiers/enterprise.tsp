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
 * Enterprise tier health API - compliance & security
 * Deployment time: ~45 minutes
 * 
 * Features:
 * - All Advanced tier features
 * - mTLS support
 * - RBAC integration
 * - Audit logging
 * - Data governance
 * - Compliance reporting
 * - Multi-environment configuration
 * - Blue/Green deployment support
 * - Gateway API support
 * - ServiceMonitor (Prometheus) integration
 * - Ingress configuration
 */
namespace EnterpriseTier;

/**
 * Audit event types
 */
enum AuditEventType {
  healthCheckAccessed: "health.check.accessed",
  configurationChanged: "config.changed",
  dependencyCheckFailed: "dependency.check.failed",
  thresholdExceeded: "threshold.exceeded",
  authenticationFailed: "auth.failed",
  authorizationDenied: "auth.denied",
  complianceViolation: "compliance.violation",
}

/**
 * Audit log entry
 */
model AuditLogEntry {
  /** Unique event ID */
  eventId: string;
  
  /** Event type */
  eventType: AuditEventType;
  
  /** Timestamp of the event */
  timestamp: utcDateTime;
  
  /** User or service that triggered the event */
  principal: {
    /** User ID or service account */
    id: string;
    
    /** Principal type */
    type: "user" | "service" | "system";
    
    /** Principal name */
    name?: string;
    
    /** IP address */
    ipAddress?: string;
    
    /** User agent */
    userAgent?: string;
  };
  
  /** Resource accessed */
  resource: {
    /** Resource type */
    type: string;
    
    /** Resource identifier */
    id: string;
    
    /** Action performed */
    action: string;
  };
  
  /** Result of the action */
  result: {
    /** Success or failure */
    status: "success" | "failure";
    
    /** HTTP status code */
    statusCode?: int32;
    
    /** Error message if failed */
    error?: string;
  };
  
  /** Additional context */
  context?: Record<unknown>;
  
  /** Data classification */
  dataClassification?: "public" | "internal" | "confidential" | "restricted";
  
  /** Compliance tags */
  complianceTags?: string[];
}

/**
 * RBAC permissions for health endpoints
 */
enum HealthPermission {
  healthRead: "health:read",
  healthWrite: "health:write",
  metricsRead: "metrics:read",
  configRead: "config:read",
  configWrite: "config:write",
  auditRead: "audit:read",
  complianceRead: "compliance:read",
}

/**
 * RBAC context
 */
model RBACContext {
  /** Authenticated principal */
  principal: {
    /** User or service ID */
    id: string;
    
    /** Assigned roles */
    roles: string[];
    
    /** Direct permissions */
    permissions: HealthPermission[];
  };
  
  /** JWT claims */
  claims?: Record<unknown>;
  
  /** mTLS certificate info */
  mtls?: {
    /** Certificate subject */
    subject: string;
    
    /** Certificate issuer */
    issuer: string;
    
    /** Certificate serial */
    serial: string;
    
    /** Valid from */
    validFrom: utcDateTime;
    
    /** Valid to */
    validTo: utcDateTime;
  };
}

/**
 * Compliance report
 */
model ComplianceReport {
  /** Report ID */
  reportId: string;
  
  /** Report timestamp */
  timestamp: utcDateTime;
  
  /** Compliance framework */
  framework: "SOC2" | "HIPAA" | "GDPR" | "ISO27001" | "PCI-DSS" | "Custom";
  
  /** Overall compliance status */
  status: "compliant" | "non-compliant" | "partial";
  
  /** Compliance checks */
  checks: {
    /** Check name */
    name: string;
    
    /** Check description */
    description: string;
    
    /** Check result */
    result: "pass" | "fail" | "warning" | "not-applicable";
    
    /** Evidence */
    evidence?: string[];
    
    /** Remediation steps if failed */
    remediation?: string;
  }[];
  
  /** Data governance metrics */
  dataGovernance: {
    /** Data classification coverage */
    classificationCoverage: float64;
    
    /** Sensitive data endpoints */
    sensitiveEndpoints: int32;
    
    /** Encryption status */
    encryptionStatus: {
      atRest: boolean;
      inTransit: boolean;
      keyRotation: boolean;
    };
  };
  
  /** Security metrics */
  securityMetrics: {
    /** mTLS enabled endpoints */
    mtlsEnabled: int32;
    
    /** RBAC coverage */
    rbacCoverage: float64;
    
    /** Audit log retention days */
    auditRetentionDays: int32;
    
    /** Last security scan */
    lastSecurityScan?: utcDateTime;
  };
  
  /** Next audit date */
  nextAudit?: utcDateTime;
}

/**
 * Multi-environment configuration
 */
model EnvironmentConfig {
  /** Environment name */
  environment: "development" | "staging" | "production" | "disaster-recovery";
  
  /** Configuration version */
  version: string;
  
  /** Feature flags */
  features: {
    /** mTLS required */
    mtlsRequired: boolean;
    
    /** Audit logging enabled */
    auditEnabled: boolean;
    
    /** Compliance mode */
    complianceMode: boolean;
    
    /** Blue/Green deployment active */
    blueGreenActive: boolean;
    
    /** Gateway API enabled */
    gatewayApiEnabled: boolean;
  };
  
  /** Environment-specific overrides */
  overrides: Record<unknown>;
  
  /** Last deployed */
  lastDeployed: utcDateTime;
  
  /** Deployment metadata */
  deployment?: {
    /** Deployment strategy */
    strategy: "rolling" | "blue-green" | "canary";
    
    /** Current version */
    currentVersion: string;
    
    /** Previous version */
    previousVersion?: string;
    
    /** Rollback available */
    rollbackAvailable: boolean;
  };
}

/**
 * Enterprise tier health API endpoints
 */
@route("/health")
@doc("Enterprise health endpoints with security and compliance")
interface EnterpriseHealthAPI extends AdvancedTier.AdvancedHealthAPI {
  /**
   * Get audit logs
   */
  @get
  @route("/audit")
  @doc("Returns audit logs (requires audit:read permission)")
  
  getAuditLogs(
    @query startTime?: utcDateTime,
    @query endTime?: utcDateTime,
    @query eventType?: AuditEventType,
    @query principalId?: string,
    @query @maxValue(1000) limit?: int32 = 100,
    @query offset?: int32 = 0
  ): {
    @statusCode statusCode: 200;
    @header("X-Total-Count") totalCount: int32;
    ...IntermediateTier.RateLimitHeaders;
    @body body: {
      entries: AuditLogEntry[];
      hasMore: boolean;
    };
  } | {
    @statusCode statusCode: 403;
    @body body: HealthError;
  };
  
  /**
   * Get compliance report
   */
  @get
  @route("/compliance")
  @doc("Returns compliance report (requires compliance:read permission)")
  
  getComplianceReport(@query framework?: ComplianceReport.framework): {
    @statusCode statusCode: 200;
    ...IntermediateTier.RateLimitHeaders;
    @body body: ComplianceReport;
  } | {
    @statusCode statusCode: 403;
    @body body: HealthError;
  };
  
  /**
   * Get environment configuration
   */
  @get
  @route("/environments/{environment}")
  @doc("Returns environment-specific configuration")
  
  getEnvironmentConfig(
    @path environment: EnvironmentConfig.environment
  ): {
    @statusCode statusCode: 200;
    @body body: EnvironmentConfig;
  } | {
    @statusCode statusCode: 403;
    @body body: HealthError;
  } | {
    @statusCode statusCode: 404;
    @body body: HealthError;
  };
  
  /**
   * Trigger blue/green deployment switch
   */
  @post
  @route("/deployment/switch")
  @doc("Switches blue/green deployment (requires health:write permission)")
  
  switchDeployment(@body target: { environment: string; version: string }): {
    @statusCode statusCode: 200;
    @body body: {
      success: boolean;
      previousVersion: string;
      currentVersion: string;
      rollbackToken: string;
    };
  } | {
    @statusCode statusCode: 403;
    @body body: HealthError;
  } | {
    @statusCode statusCode: 409;
    @body body: HealthError;
  };
  
  /**
   * Get RBAC context for current user
   */
  @get
  @route("/auth/context")
  @doc("Returns RBAC context for authenticated user")
  
  getRBACContext(): {
    @statusCode statusCode: 200;
    @body body: RBACContext;
  } | {
    @statusCode statusCode: 401;
    @body body: HealthError;
  };
}

/**
 * Prometheus ServiceMonitor configuration
 */
model ServiceMonitorConfig {
  /** Scrape interval */
  interval: string = "30s";
  
  /** Scrape timeout */
  timeout: string = "10s";
  
  /** Metrics path */
  path: string = "/metrics";
  
  /** Port name */
  port: string = "http";
  
  /** Additional labels */
  labels?: Record<string>;
  
  /** Relabeling configs */
  relabelings?: Record<unknown>[];
}

/**
 * Gateway API route configuration
 */
model GatewayRouteConfig {
  /** Route name */
  name: string;
  
  /** Hostnames */
  hostnames: string[];
  
  /** Path matches */
  paths: {
    type: "Exact" | "PathPrefix" | "RegularExpression";
    value: string;
  }[];
  
  /** Backend services */
  backends: {
    name: string;
    port: int32;
    weight?: int32;
  }[];
  
  /** TLS configuration */
  tls?: {
    certificateRefs: string[];
    mode: "Terminate" | "Passthrough";
  };
}