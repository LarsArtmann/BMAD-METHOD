# BMAD Method Template Health Endpoint Generator - Setup Guide

## 🚀 Quick Start (5 minutes)

### Prerequisites

- **Go 1.21+** - [Download Go](https://golang.org/dl/)
- **Node.js 16+** (optional, for TypeScript client) - [Download Node.js](https://nodejs.org/)
- **Docker** (optional, for containerization) - [Download Docker](https://docker.com/)
- **kubectl** (optional, for Kubernetes deployment) - [Install kubectl](https://kubernetes.io/docs/tasks/tools/)

### Installation

1. **Clone the repository:**
   ```bash
   git clone https://github.com/LarsArtmann/BMAD-METHOD.git
   cd BMAD-METHOD
   ```

2. **Build the CLI tool:**
   ```bash
   go build -o bin/template-health-endpoint cmd/generator/main.go
   ```

3. **Verify installation:**
   ```bash
   ./bin/template-health-endpoint --help
   ```

### Generate Your First Health Service

1. **Generate a basic health service:**
   ```bash
   ./bin/template-health-endpoint generate \
     --name my-health-service \
     --tier basic \
     --module github.com/yourorg/my-health-service
   ```

2. **Build and run the service:**
   ```bash
   cd my-health-service
   go mod tidy
   go run cmd/server/main.go
   ```

3. **Test the health endpoints:**
   ```bash
   # Basic health check
   curl http://localhost:8080/health
   
   # Server time information
   curl http://localhost:8080/health/time
   
   # Kubernetes probes
   curl http://localhost:8080/health/ready
   curl http://localhost:8080/health/live
   curl http://localhost:8080/health/startup
   ```

## 📋 Template Tiers

### Basic Tier (~5 minutes deployment)
- ✅ Core health endpoints (`/health`, `/health/time`, `/health/ready`, `/health/live`, `/health/startup`)
- ✅ Go server implementation with graceful shutdown
- ✅ TypeScript client SDK
- ✅ Docker containerization
- ✅ Kubernetes manifests with health probes
- ✅ Comprehensive documentation

**Perfect for:** Simple services, microservices, development environments

### Intermediate Tier (~15 minutes deployment)
- 🔄 Everything from Basic tier
- 🔄 Dependency health checks (database, cache, external APIs)
- 🔄 Basic OpenTelemetry integration
- 🔄 Enhanced error handling and logging
- 🔄 Configuration management

**Perfect for:** Production services with external dependencies

### Advanced Tier (~30 minutes deployment)
- 🔄 Everything from Intermediate tier
- 🔄 Full OpenTelemetry observability (traces, metrics, logs)
- 🔄 Server Timing API for performance metrics
- 🔄 CloudEvents emission for event-driven monitoring
- 🔄 Advanced health check strategies

**Perfect for:** Mission-critical services requiring full observability

### Enterprise Tier (~45 minutes deployment)
- 🔄 Everything from Advanced tier
- 🔄 Compliance features (audit logging, data governance)
- 🔄 Advanced security (mTLS, RBAC integration)
- 🔄 ServiceMonitor for Prometheus integration
- 🔄 Multi-environment configuration

**Perfect for:** Enterprise applications with strict compliance requirements

## 🛠️ CLI Commands

### Generate Command

```bash
./bin/template-health-endpoint generate [flags]
```

**Flags:**
- `--name` (required) - Project name
- `--tier` (required) - Template tier (basic, intermediate, advanced, enterprise)
- `--module` (required) - Go module path
- `--output` - Output directory (default: project name)
- `--dry-run` - Preview generation without creating files
- `--features` - Comma-separated feature flags

**Examples:**
```bash
# Basic service
./bin/template-health-endpoint generate \
  --name user-service \
  --tier basic \
  --module github.com/myorg/user-service

# Advanced service with specific features
./bin/template-health-endpoint generate \
  --name payment-service \
  --tier advanced \
  --module github.com/myorg/payment-service \
  --features opentelemetry,cloudevents

# Dry run to preview
./bin/template-health-endpoint generate \
  --name test-service \
  --tier basic \
  --module github.com/test/service \
  --dry-run
```

### Validate Command

```bash
./bin/template-health-endpoint validate [flags]
```

**Flags:**
- `--schemas` - Path to TypeSpec schemas (default: pkg/schemas/)
- `--emit` - Output formats (openapi3, json-schema)
- `--output` - Output directory for generated schemas
- `--verbose` - Detailed validation output

**Examples:**
```bash
# Validate all schemas
./bin/template-health-endpoint validate --verbose

# Generate OpenAPI and JSON Schema
./bin/template-health-endpoint validate \
  --emit openapi3,json-schema \
  --output generated-schemas
```

## 🏗️ Project Structure

Generated projects follow a consistent, production-ready structure:

```
my-health-service/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── handlers/
│   │   └── health.go            # Health endpoint handlers
│   ├── models/
│   │   └── health.go            # Data models
│   ├── server/
│   │   └── server.go            # HTTP server setup
│   └── config/
│       └── config.go            # Configuration management
├── client/
│   └── typescript/              # TypeScript client SDK
│       ├── src/
│       │   ├── client.ts        # API client
│       │   └── types.ts         # Type definitions
│       ├── package.json
│       └── tsconfig.json
├── deployments/
│   └── kubernetes/              # Kubernetes manifests
│       ├── deployment.yaml      # Deployment with health probes
│       ├── service.yaml         # Service definition
│       └── configmap.yaml       # Configuration
├── docs/
│   └── API.md                   # API documentation
├── scripts/
│   ├── build.sh                 # Build script
│   └── test.sh                  # Test script
├── Dockerfile                   # Multi-stage Docker build
├── docker-compose.yml           # Local development
├── Makefile                     # Build automation
├── go.mod                       # Go module definition
└── README.md                    # Project documentation
```

## 🔧 Configuration

### Environment Variables

Generated services support configuration via environment variables:

```bash
# Server configuration
PORT=8080                        # HTTP server port
VERSION=1.0.0                    # Service version

# Health check configuration
HEALTH_CHECK_TIMEOUT=5s          # Health check timeout
HEALTH_CHECK_INTERVAL=30s        # Health check interval

# Observability (Advanced+ tiers)
OTEL_SERVICE_NAME=my-service     # OpenTelemetry service name
OTEL_EXPORTER_OTLP_ENDPOINT=...  # OTLP endpoint
```

### Configuration Files

For complex configurations, use YAML files:

```yaml
# config.yaml
server:
  port: 8080
  timeout: 30s

health:
  checks:
    database:
      enabled: true
      timeout: 5s
    cache:
      enabled: true
      timeout: 2s

observability:
  opentelemetry:
    enabled: true
    service_name: my-service
```

## 🐳 Docker Deployment

### Build and Run

```bash
# Build Docker image
docker build -t my-health-service:latest .

# Run container
docker run -p 8080:8080 my-health-service:latest

# Test health endpoints
curl http://localhost:8080/health
```

### Docker Compose

```bash
# Start with dependencies
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down
```

## ☸️ Kubernetes Deployment

### Deploy to Kubernetes

```bash
# Apply manifests
kubectl apply -f deployments/kubernetes/

# Check deployment status
kubectl get pods -l app=my-health-service

# Check health probes
kubectl describe pod <pod-name>

# Test service
kubectl port-forward svc/my-health-service 8080:8080
curl http://localhost:8080/health
```

### Health Probe Configuration

Generated Kubernetes manifests include properly configured health probes:

```yaml
# Liveness probe - restart if unhealthy
livenessProbe:
  httpGet:
    path: /health/live
    port: 8080
  initialDelaySeconds: 30
  periodSeconds: 10

# Readiness probe - remove from load balancer if not ready
readinessProbe:
  httpGet:
    path: /health/ready
    port: 8080
  initialDelaySeconds: 5
  periodSeconds: 5

# Startup probe - allow time for startup
startupProbe:
  httpGet:
    path: /health/startup
    port: 8080
  initialDelaySeconds: 10
  periodSeconds: 10
  failureThreshold: 30
```

## 🧪 Testing

### Run Tests

```bash
# Run all tests
./scripts/test.sh

# Run unit tests only
go test -v ./pkg/generator/...

# Run integration tests only
go test -v ./tests/...
```

### Test Generated Services

```bash
# Build and test generated service
cd my-health-service
go mod tidy
go test ./...
go build -o bin/my-health-service cmd/server/main.go

# Start service and test endpoints
./bin/my-health-service &
curl http://localhost:8080/health
curl http://localhost:8080/health/startup
```

## 🔍 Troubleshooting

### Common Issues

**1. Build Failures**
```bash
# Ensure Go version is 1.21+
go version

# Clean module cache
go clean -modcache
go mod tidy
```

**2. Port Already in Use**
```bash
# Find process using port 8080
lsof -i :8080

# Kill process
kill -9 <PID>

# Or use different port
PORT=8081 ./bin/my-health-service
```

**3. Health Endpoints Not Responding**
```bash
# Check if service is running
ps aux | grep my-health-service

# Check logs
./bin/my-health-service 2>&1 | tee service.log

# Test with verbose curl
curl -v http://localhost:8080/health
```

**4. Kubernetes Deployment Issues**
```bash
# Check pod status
kubectl get pods -l app=my-health-service

# Check pod logs
kubectl logs -l app=my-health-service

# Check events
kubectl get events --sort-by=.metadata.creationTimestamp
```

### Getting Help

- 📖 **Documentation**: Check `docs/` directory
- 🐛 **Issues**: [GitHub Issues](https://github.com/LarsArtmann/BMAD-METHOD/issues)
- 💬 **Discussions**: [GitHub Discussions](https://github.com/LarsArtmann/BMAD-METHOD/discussions)
- 📧 **Contact**: git@lars.software

## 🎯 Next Steps

1. **Generate your first service** using the basic tier
2. **Deploy to Kubernetes** and test health probes
3. **Explore advanced tiers** for production requirements
4. **Customize templates** for your organization's needs
5. **Contribute back** improvements and new features

---

**Generated by BMAD Method Template Health Endpoint Generator**  
*Following the Business, Management, Architecture, Development methodology*
