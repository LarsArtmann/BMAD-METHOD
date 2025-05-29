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
import "./pkg/schemas/tiers/intermediate.tsp";
import "./pkg/schemas/tiers/advanced.tsp";
import "./pkg/schemas/tiers/enterprise.tsp";

using TypeSpec.Http;
using TypeSpec.Rest;

/**
 * Template Health Endpoint API
 *
 * A comprehensive TypeSpec-driven health endpoint template system
 * supporting four tiers of complexity: Basic, Intermediate, Advanced, Enterprise
 */
@server("https://api.example.com", "Production server")
@server("http://localhost:8080", "Development server")
namespace TemplateHealthEndpoint;

// Re-export main interfaces for easy access
alias MainHealthAPI = HealthEndpoints;
alias ExtendedAPI = ExtendedHealthEndpoints;
alias BasicAPI = BasicTier.BasicHealthAPI;
alias IntermediateAPI = IntermediateTier.IntermediateHealthAPI;
alias AdvancedAPI = AdvancedTier.AdvancedHealthAPI;
alias EnterpriseAPI = EnterpriseTier.EnterpriseHealthAPI;

// Re-export core models
alias HealthStatus = HealthAPI.HealthStatus;
alias HealthReport = HealthAPI.HealthReport;
alias HealthEvent = HealthEvents.HealthEvent;
