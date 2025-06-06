# {{.Config.Name}} API Documentation

This document describes the health endpoints provided by {{.Config.Name}}.

## Base URL

```
http://localhost:8080
```

## Endpoints

### GET /health

Returns the overall health status of the service.

**Response:**
```json
{
  "status": "healthy",
  "timestamp": "2024-01-01T12:00:00Z",
  "version": "1.0.0",
  "uptime": 3600000000000,
  "uptime_human": "1.0 hours"
}
```

### GET /health/time

Returns server time information in multiple formats.

**Response:**
```json
{
  "timestamp": "2024-01-01T12:00:00Z",
  "timezone": "UTC",
  "unix": 1704110400,
  "unix_milli": 1704110400000,
  "iso8601": "2024-01-01T12:00:00Z",
  "formatted": "Monday, January 1, 2024 at 12:00:00 PM UTC"
}
```

### GET /health/ready

Kubernetes readiness probe endpoint.

**Response:** Same as /health

### GET /health/live

Kubernetes liveness probe endpoint.

**Response:** Same as /health

### GET /health/startup

Kubernetes startup probe endpoint.

**Response:** Same as /health

## Status Codes

- `200 OK` - Service is healthy
- `503 Service Unavailable` - Service is unhealthy

## Generated by

Template Health Endpoint Generator v{{.Version}}
Generated at: {{.Timestamp}}
