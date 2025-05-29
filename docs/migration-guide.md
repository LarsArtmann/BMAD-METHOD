# Migration Guide - Upgrading Between Template Tiers

## 🎯 Overview

This guide helps you migrate between different template tiers as your service requirements evolve. Each migration path is designed to be safe, incremental, and minimize downtime.

## 🔄 Migration Paths

### Basic → Intermediate
**Estimated Time:** 30-45 minutes  
**Downtime:** < 5 minutes  
**Complexity:** Low

### Intermediate → Advanced
**Estimated Time:** 1-2 hours  
**Downtime:** < 10 minutes  
**Complexity:** Medium

### Advanced → Enterprise
**Estimated Time:** 2-4 hours  
**Downtime:** < 15 minutes  
**Complexity:** High

### Skip Tiers (e.g., Basic → Advanced)
**Estimated Time:** 2-3 hours  
**Downtime:** < 15 minutes  
**Complexity:** Medium-High

## 📋 Pre-Migration Checklist

### Before Starting Any Migration

- [ ] **Backup current service**
  ```bash
  # Create backup of current service
  cp -r my-service my-service-backup-$(date +%Y%m%d)
  ```

- [ ] **Document current configuration**
  ```bash
  # Save current environment variables
  env | grep -E "(PORT|VERSION|HEALTH_)" > current-config.env
  
  # Save current Kubernetes configuration
  kubectl get deployment my-service -o yaml > current-k8s-config.yaml
  ```

- [ ] **Test current service**
  ```bash
  # Verify all endpoints work
  curl http://localhost:8080/health
  curl http://localhost:8080/health/ready
  curl http://localhost:8080/health/live
  curl http://localhost:8080/health/startup
  ```

- [ ] **Check dependencies**
  ```bash
  # Verify external dependencies
  go mod verify
  go mod tidy
  ```

- [ ] **Plan rollback strategy**
  - Keep backup of working version
  - Document rollback steps
  - Prepare rollback scripts

## 🚀 Migration: Basic → Intermediate

### What Changes
- ✅ Adds dependency health checks
- ✅ Adds basic OpenTelemetry integration
- ✅ Enhances configuration management
- ✅ Adds structured logging

### Step-by-Step Migration

#### 1. Generate New Intermediate Service

```bash
# Generate intermediate tier in new directory
./bin/template-health-endpoint generate \
  --name my-service-intermediate \
  --tier intermediate \
  --module github.com/yourorg/my-service

# Compare with current service
diff -r my-service my-service-intermediate
```

#### 2. Update Configuration

**Add dependency configuration:**
```yaml
# config/config.yaml
dependencies:
  database:
    enabled: true
    url: "postgres://user:pass@localhost:5432/db"
    timeout: 5s
  cache:
    enabled: true
    url: "redis://localhost:6379"
    timeout: 2s
  external_api:
    enabled: true
    url: "https://api.external.com"
    timeout: 10s
```

**Update environment variables:**
```bash
# Add to your environment
export OTEL_SERVICE_NAME=my-service
export OTEL_EXPORTER_OTLP_ENDPOINT=http://jaeger:14268/api/traces
export LOG_LEVEL=info
export LOG_FORMAT=json
```

#### 3. Update Dependencies

```bash
cd my-service-intermediate

# Add new dependencies
go mod tidy

# Verify dependencies
go mod verify
```

#### 4. Test New Features

```bash
# Build and run
go build -o bin/my-service cmd/server/main.go
./bin/my-service &

# Test new dependency endpoint
curl http://localhost:8080/health/dependencies

# Expected response:
# {
#   "database": {"status": "healthy", "response_time": "2ms"},
#   "cache": {"status": "healthy", "response_time": "1ms"},
#   "external_api": {"status": "healthy", "response_time": "50ms"}
# }
```

#### 5. Update Kubernetes Manifests

```bash
# Apply new ConfigMap with dependency config
kubectl apply -f deployments/kubernetes/configmap.yaml

# Update deployment with new image
kubectl set image deployment/my-service my-service=my-service:intermediate

# Monitor rollout
kubectl rollout status deployment/my-service
```

#### 6. Verify Migration

```bash
# Check all endpoints
curl http://your-service/health
curl http://your-service/health/dependencies

# Check Kubernetes health probes
kubectl describe pod -l app=my-service

# Check logs for OpenTelemetry traces
kubectl logs -l app=my-service | grep -i trace
```

### Rollback Plan (Basic ← Intermediate)

```bash
# If issues occur, rollback to basic tier
kubectl set image deployment/my-service my-service=my-service:basic
kubectl rollout status deployment/my-service

# Or restore from backup
rm -rf my-service
mv my-service-backup-* my-service
```

## 🔥 Migration: Intermediate → Advanced

### What Changes
- ✅ Full OpenTelemetry observability
- ✅ Server Timing API
- ✅ CloudEvents emission
- ✅ Advanced metrics endpoints
- ✅ Circuit breaker patterns

### Step-by-Step Migration

#### 1. Prepare Infrastructure

**Set up OpenTelemetry Collector:**
```yaml
# otel-collector.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: otel-collector-config
data:
  config.yaml: |
    receivers:
      otlp:
        protocols:
          grpc:
            endpoint: 0.0.0.0:4317
    exporters:
      jaeger:
        endpoint: jaeger:14250
      prometheus:
        endpoint: "0.0.0.0:8889"
    service:
      pipelines:
        traces:
          receivers: [otlp]
          exporters: [jaeger]
        metrics:
          receivers: [otlp]
          exporters: [prometheus]
```

**Set up CloudEvents infrastructure:**
```bash
# Install CloudEvents broker (e.g., Knative Eventing)
kubectl apply -f https://github.com/knative/eventing/releases/download/knative-v1.8.0/eventing-crds.yaml
kubectl apply -f https://github.com/knative/eventing/releases/download/knative-v1.8.0/eventing-core.yaml
```

#### 2. Generate Advanced Service

```bash
./bin/template-health-endpoint generate \
  --name my-service-advanced \
  --tier advanced \
  --module github.com/yourorg/my-service \
  --features opentelemetry,cloudevents,server-timing
```

#### 3. Update Configuration

**Add observability configuration:**
```yaml
# config/config.yaml
observability:
  opentelemetry:
    enabled: true
    service_name: my-service
    endpoint: http://otel-collector:4317
  server_timing:
    enabled: true
  cloudevents:
    enabled: true
    broker_url: http://broker-ingress.knative-eventing.svc.cluster.local
```

#### 4. Test Advanced Features

```bash
# Test metrics endpoint
curl http://localhost:8080/health/metrics

# Expected response with detailed metrics:
# {
#   "systemInfo": {
#     "cpu_usage": 15.5,
#     "memory_usage": 256000000,
#     "goroutines": 42
#   },
#   "performance": {
#     "requests_per_second": 150.5,
#     "average_response_time": "5ms"
#   },
#   "traceId": "abc123def456"
# }

# Check Server-Timing headers
curl -I http://localhost:8080/health
# Look for: Server-Timing: total;dur=2.5, db;dur=1.2, cache;dur=0.8
```

#### 5. Deploy with Blue-Green Strategy

```bash
# Deploy advanced version alongside current
kubectl apply -f deployments/kubernetes/ --dry-run=client -o yaml | \
  sed 's/my-service/my-service-advanced/g' | \
  kubectl apply -f -

# Test advanced version
kubectl port-forward svc/my-service-advanced 8081:8080 &
curl http://localhost:8081/health/metrics

# Switch traffic gradually
kubectl patch service my-service -p '{"spec":{"selector":{"version":"advanced"}}}'
```

#### 6. Monitor Migration

```bash
# Monitor traces in Jaeger
open http://jaeger-ui:16686

# Monitor metrics in Prometheus
open http://prometheus:9090

# Check CloudEvents
kubectl logs -l app=cloudevents-broker
```

### Rollback Plan (Intermediate ← Advanced)

```bash
# Quick rollback - switch service selector
kubectl patch service my-service -p '{"spec":{"selector":{"version":"intermediate"}}}'

# Full rollback - remove advanced deployment
kubectl delete deployment my-service-advanced
```

## 🏢 Migration: Advanced → Enterprise

### What Changes
- ✅ mTLS and advanced security
- ✅ RBAC integration
- ✅ Audit logging
- ✅ Compliance features
- ✅ Multi-environment configuration

### Step-by-Step Migration

#### 1. Prepare Security Infrastructure

**Set up Certificate Authority:**
```bash
# Generate CA certificate
openssl genrsa -out ca-key.pem 4096
openssl req -new -x509 -days 365 -key ca-key.pem -out ca.pem

# Create Kubernetes secret
kubectl create secret tls ca-secret --cert=ca.pem --key=ca-key.pem
```

**Set up RBAC:**
```yaml
# rbac.yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: health-endpoint-reader
rules:
- apiGroups: [""]
  resources: ["pods", "services"]
  verbs: ["get", "list"]
---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: health-endpoint-binding
subjects:
- kind: ServiceAccount
  name: my-service
roleRef:
  kind: Role
  name: health-endpoint-reader
  apiGroup: rbac.authorization.k8s.io
```

#### 2. Generate Enterprise Service

```bash
./bin/template-health-endpoint generate \
  --name my-service-enterprise \
  --tier enterprise \
  --module github.com/yourorg/my-service \
  --features opentelemetry,cloudevents,security,compliance
```

#### 3. Configure Security

**Add security configuration:**
```yaml
# config/security.yaml
security:
  mtls:
    enabled: true
    ca_cert_path: /etc/certs/ca.pem
    cert_path: /etc/certs/server.pem
    key_path: /etc/certs/server-key.pem
  rbac:
    enabled: true
    provider: kubernetes
  audit:
    enabled: true
    log_path: /var/log/audit.log
    retention_days: 90
```

#### 4. Test Enterprise Features

```bash
# Test with mTLS
curl --cert client.pem --key client-key.pem --cacert ca.pem \
  https://localhost:8443/health

# Test audit logging
tail -f /var/log/audit.log

# Test compliance endpoint
curl https://localhost:8443/health/compliance
```

#### 5. Deploy with Security

```bash
# Create service account
kubectl apply -f deployments/kubernetes/serviceaccount.yaml

# Apply RBAC
kubectl apply -f rbac.yaml

# Deploy with security context
kubectl apply -f deployments/kubernetes/
```

### Rollback Plan (Advanced ← Enterprise)

```bash
# Remove security features
kubectl delete -f rbac.yaml
kubectl patch deployment my-service --type='json' \
  -p='[{"op": "remove", "path": "/spec/template/spec/securityContext"}]'
```

## 🔧 Common Migration Issues

### Issue: Port Conflicts
**Problem:** New tier uses different ports
**Solution:**
```bash
# Update service configuration
kubectl patch service my-service -p '{"spec":{"ports":[{"port":8080,"targetPort":8080}]}}'
```

### Issue: Configuration Mismatch
**Problem:** Environment variables not compatible
**Solution:**
```bash
# Create migration script
cat > migrate-config.sh << 'EOF'
#!/bin/bash
# Convert old config to new format
OLD_CONFIG=${1:-config.env}
NEW_CONFIG=${2:-config.yaml}

# Convert environment variables to YAML
echo "server:" > $NEW_CONFIG
echo "  port: ${PORT:-8080}" >> $NEW_CONFIG
echo "  version: ${VERSION:-1.0.0}" >> $NEW_CONFIG
EOF

chmod +x migrate-config.sh
./migrate-config.sh
```

### Issue: Health Check Failures
**Problem:** New endpoints not responding
**Solution:**
```bash
# Check endpoint availability
kubectl exec -it deployment/my-service -- curl localhost:8080/health/dependencies

# Check logs
kubectl logs -l app=my-service --tail=100

# Verify configuration
kubectl get configmap my-service-config -o yaml
```

### Issue: Performance Degradation
**Problem:** New tier has higher resource usage
**Solution:**
```bash
# Update resource limits
kubectl patch deployment my-service -p '{
  "spec": {
    "template": {
      "spec": {
        "containers": [{
          "name": "my-service",
          "resources": {
            "requests": {"memory": "128Mi", "cpu": "100m"},
            "limits": {"memory": "256Mi", "cpu": "200m"}
          }
        }]
      }
    }
  }
}'
```

## 📊 Migration Validation Checklist

### Post-Migration Verification

- [ ] **All health endpoints respond correctly**
  ```bash
  curl http://your-service/health
  curl http://your-service/health/ready
  curl http://your-service/health/live
  curl http://your-service/health/startup
  # Tier-specific endpoints
  curl http://your-service/health/dependencies  # Intermediate+
  curl http://your-service/health/metrics       # Advanced+
  ```

- [ ] **Kubernetes health probes working**
  ```bash
  kubectl describe pod -l app=my-service | grep -A 10 "Liveness\|Readiness\|Startup"
  ```

- [ ] **Performance within acceptable range**
  ```bash
  # Load test
  ab -n 1000 -c 10 http://your-service/health
  ```

- [ ] **Observability data flowing**
  ```bash
  # Check traces (Advanced+)
  curl http://jaeger:16686/api/traces?service=my-service
  
  # Check metrics (Advanced+)
  curl http://prometheus:9090/api/v1/query?query=up{job="my-service"}
  ```

- [ ] **Security features working** (Enterprise)
  ```bash
  # Test mTLS
  curl --cert client.pem --key client-key.pem --cacert ca.pem https://your-service/health
  
  # Check audit logs
  kubectl logs -l app=my-service | grep audit
  ```

## 🎯 Best Practices

### 1. Gradual Migration
- Migrate one environment at a time (dev → staging → production)
- Use feature flags to enable new features gradually
- Monitor each step before proceeding

### 2. Testing Strategy
- Run comprehensive tests before migration
- Use canary deployments for production
- Have automated rollback triggers

### 3. Documentation
- Document all configuration changes
- Update runbooks and operational procedures
- Train team on new features

### 4. Monitoring
- Set up alerts for new metrics
- Monitor performance during migration
- Track error rates and response times

## 🆘 Emergency Procedures

### Immediate Rollback
```bash
# Emergency rollback script
cat > emergency-rollback.sh << 'EOF'
#!/bin/bash
set -e

echo "🚨 Emergency rollback initiated"

# Restore from backup
kubectl set image deployment/my-service my-service=my-service:backup
kubectl rollout status deployment/my-service

# Verify health
kubectl wait --for=condition=ready pod -l app=my-service --timeout=60s

echo "✅ Rollback complete"
EOF

chmod +x emergency-rollback.sh
```

### Health Check Script
```bash
# health-check.sh
#!/bin/bash
SERVICE_URL=${1:-http://localhost:8080}

echo "🔍 Checking service health..."

# Basic health check
if curl -f $SERVICE_URL/health > /dev/null 2>&1; then
  echo "✅ Basic health check passed"
else
  echo "❌ Basic health check failed"
  exit 1
fi

# Check all endpoints
for endpoint in ready live startup; do
  if curl -f $SERVICE_URL/health/$endpoint > /dev/null 2>&1; then
    echo "✅ /$endpoint endpoint working"
  else
    echo "❌ /$endpoint endpoint failed"
    exit 1
  fi
done

echo "🎉 All health checks passed!"
```

---

**Need help with migration?** Check the [troubleshooting section](setup-guide.md#troubleshooting) or [open an issue](https://github.com/LarsArtmann/BMAD-METHOD/issues).
