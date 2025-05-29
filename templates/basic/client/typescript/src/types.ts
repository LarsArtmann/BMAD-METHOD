// Generated TypeScript types for {{.Config.Name}}
// Generated at: {{.Timestamp}}

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
