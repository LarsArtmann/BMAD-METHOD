// Generated TypeScript types for test-health-service
// Generated at: 2024-01-01T00:00:00Z

export interface HealthReport {
  status: string;
  timestamp: string;
  version: string;
  uptime: number;
  uptime_human: string;
}

export interface ServerTime {
  timestamp: string;
  timezone: string;
  unix: number;
  unix_milli: number;
  iso8601: string;
  formatted: string;
}

export type HealthStatus = 'healthy' | 'degraded' | 'unhealthy';
