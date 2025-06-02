# Multi-stage build for BMAD Method CLI
FROM node:20-alpine AS typespec-builder

# Install TypeSpec
RUN npm install -g @typespec/compiler @typespec/openapi3 @typespec/json-schema

FROM golang:1.21-alpine AS go-builder

# Install build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w -X main.version=${VERSION:-dev} -X main.buildDate=$(date -u +%Y-%m-%dT%H:%M:%SZ)" \
    -o bmad-method \
    ./cmd/generator

# Final stage
FROM alpine:3.19

# Install runtime dependencies
RUN apk add --no-cache \
    ca-certificates \
    git \
    curl \
    bash

# Create non-root user
RUN adduser -D -s /bin/sh bmad

# Copy TypeSpec from first stage
COPY --from=typespec-builder /usr/local/lib/node_modules /usr/local/lib/node_modules
COPY --from=typespec-builder /usr/local/bin/node /usr/local/bin/node
COPY --from=typespec-builder /usr/local/bin/npm /usr/local/bin/npm
RUN ln -s /usr/local/lib/node_modules/@typespec/compiler/cmd/tsp.js /usr/local/bin/tsp && \
    chmod +x /usr/local/bin/tsp

# Copy binary from builder stage
COPY --from=go-builder /app/bmad-method /usr/local/bin/bmad-method

# Copy templates and schemas
COPY --from=go-builder /app/template-health /app/template-health
COPY --from=go-builder /app/templates /app/templates
COPY --from=go-builder /app/pkg /app/pkg

# Set working directory
WORKDIR /workspace

# Switch to non-root user
USER bmad

# Default command
ENTRYPOINT ["bmad-method"]
CMD ["--help"]

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD bmad-method --version || exit 1

# Labels
LABEL maintainer="Lars Artmann <lars@example.com>"
LABEL org.opencontainers.image.title="BMAD Method"
LABEL org.opencontainers.image.description="TypeSpec-driven health endpoint template generator"
LABEL org.opencontainers.image.vendor="BMAD Method"
LABEL org.opencontainers.image.licenses="MIT"
LABEL org.opencontainers.image.source="https://github.com/LarsArtmann/BMAD-METHOD"
LABEL org.opencontainers.image.documentation="https://github.com/LarsArtmann/BMAD-METHOD#readme"