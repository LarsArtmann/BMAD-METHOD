import "@typespec/http";
import "@typespec/rest";
import "@typespec/openapi3";

// Import all health endpoint schemas
import "./pkg/schemas/health/health.tsp";
import "./pkg/schemas/health/server-time.tsp";
import "./pkg/schemas/health/health-api.tsp";
import "./pkg/schemas/health/cloudevents.tsp";

// Import tier-specific schemas
import "./pkg/schemas/tiers/basic.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;

/**
 * Template Health Endpoint API
 *
 * A comprehensive TypeSpec-driven health endpoint template system
 * supporting four tiers of complexity: Basic, Intermediate, Advanced, Enterprise
 */
@service({
  title: "Template Health Endpoint API",
  version: "1.0.0",
  description: "TypeSpec-driven health endpoint templates with progressive complexity tiers"
})
@server("https://api.example.com", "Production server")
@server("http://localhost:8080", "Development server")
namespace TemplateHealthEndpoint;

// Re-export main interfaces for easy access
alias MainHealthAPI = HealthEndpoints;
alias ExtendedAPI = ExtendedHealthEndpoints;
alias BasicAPI = BasicTier.BasicHealthAPI;

// Re-export core models
alias HealthStatus = HealthAPI.HealthStatus;
alias HealthReport = HealthAPI.HealthReport;
alias HealthEvent = HealthEvents.HealthEvent;
