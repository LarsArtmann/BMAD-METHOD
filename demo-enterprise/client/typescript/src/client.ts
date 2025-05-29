// Generated TypeScript client for demo-enterprise
// Generated at: 2024-01-01T00:00:00Z

import { HealthReport, ServerTime } from './types';

export interface HealthClientConfig {
  baseURL: string;
  timeout?: number;
  headers?: Record<string, string>;
}

export class HealthClient {
  private baseURL: string;
  private timeout: number;
  private headers: Record<string, string>;

  constructor(config: HealthClientConfig) {
    this.baseURL = config.baseURL.replace(/\/$/, '');
    this.timeout = config.timeout || 5000;
    this.headers = config.headers || {};
  }

  /**
   * Check the health status of the service
   */
  async checkHealth(): Promise<HealthReport> {
    return this.request<HealthReport>('/health');
  }

  /**
   * Get server time information
   */
  async getServerTime(): Promise<ServerTime> {
    return this.request<ServerTime>('/health/time');
  }

  /**
   * Check readiness status
   */
  async checkReadiness(): Promise<HealthReport> {
    return this.request<HealthReport>('/health/ready');
  }

  /**
   * Check liveness status
   */
  async checkLiveness(): Promise<HealthReport> {
    return this.request<HealthReport>('/health/live');
  }

  /**
   * Check startup status
   */
  async checkStartup(): Promise<HealthReport> {
    return this.request<HealthReport>('/health/startup');
  }

  private async request<T>(path: string): Promise<T> {
    const controller = new AbortController();
    const timeoutId = setTimeout(() => controller.abort(), this.timeout);

    try {
      const response = await fetch(`${this.baseURL}${path}`, {
        method: 'GET',
        headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json',
          ...this.headers,
        },
        signal: controller.signal,
      });

      clearTimeout(timeoutId);

      if (!response.ok) {
        throw new Error(`HTTP ${response.status}: ${response.statusText}`);
      }

      return await response.json();
    } catch (error) {
      clearTimeout(timeoutId);
      if (error instanceof Error && error.name === 'AbortError') {
        throw new Error(`Request timeout after ${this.timeout}ms`);
      }
      throw error;
    }
  }
}

// Default export for convenience
export default HealthClient;
