package generator

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/LarsArtmann/BMAD-METHOD/pkg/config"
)

func TestGenerator_Generate(t *testing.T) {
	tests := []struct {
		name    string
		config  *config.ProjectConfig
		wantErr bool
	}{
		{
			name: "basic tier generation",
			config: &config.ProjectConfig{
				Name:        "test-service",
				Description: "Test health service",
				GoModule:    "github.com/example/test-service",
				Tier:        config.TierBasic,
				Version:     "1.0.0",
				OutputDir:   "test-output-basic",
				Features: config.FeatureFlags{
					Kubernetes: true,
				},
			},
			wantErr: false,
		},
		{
			name: "intermediate tier generation",
			config: &config.ProjectConfig{
				Name:        "test-intermediate",
				Description: "Test intermediate service",
				GoModule:    "github.com/example/test-intermediate",
				Tier:        config.TierIntermediate,
				Version:     "1.0.0",
				OutputDir:   "test-output-intermediate",
				Features: config.FeatureFlags{
					Kubernetes:    true,
					OpenTelemetry: true,
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up before test
			os.RemoveAll(tt.config.OutputDir)
			defer os.RemoveAll(tt.config.OutputDir)

			generator, err := New(tt.config)
			if err != nil {
				t.Fatalf("Failed to create generator: %v", err)
			}

			err = generator.Generate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Generator.Generate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify essential files were created
				essentialFiles := []string{
					"README.md",
					"go.mod",
					"cmd/server/main.go",
					"internal/handlers/health.go",
					"internal/server/server.go",
					"deployments/kubernetes/deployment.yaml",
				}

				for _, file := range essentialFiles {
					fullPath := filepath.Join(tt.config.OutputDir, file)
					if _, err := os.Stat(fullPath); os.IsNotExist(err) {
						t.Errorf("Essential file not created: %s", file)
					}
				}
			}
		})
	}
}

func TestGenerator_ValidateHealthEndpoints(t *testing.T) {
	// Create a test project
	config := &config.ProjectConfig{
		Name:        "endpoint-test",
		Description: "Test endpoint validation",
		GoModule:    "github.com/example/endpoint-test",
		Tier:        config.TierBasic,
		Version:     "1.0.0",
		OutputDir:   "test-endpoint-validation",
		Features: config.FeatureFlags{
			Kubernetes: true,
		},
	}

	// Clean up
	defer os.RemoveAll(config.OutputDir)

	generator, err := New(config)
	if err != nil {
		t.Fatalf("Failed to create generator: %v", err)
	}

	err = generator.Generate()
	if err != nil {
		t.Fatalf("Failed to generate project: %v", err)
	}

	// Read the generated health handler
	handlerPath := filepath.Join(config.OutputDir, "internal/handlers/health.go")
	content, err := os.ReadFile(handlerPath)
	if err != nil {
		t.Fatalf("Failed to read health handler: %v", err)
	}

	handlerContent := string(content)

	// Verify all required endpoints are present
	requiredMethods := []string{
		"CheckHealth",
		"ServerTime", 
		"ReadinessCheck",
		"LivenessCheck",
		"StartupCheck", // This is the new endpoint we added
	}

	for _, method := range requiredMethods {
		if !contains(handlerContent, method) {
			t.Errorf("Required method %s not found in generated handler", method)
		}
	}

	// Verify server routing includes startup endpoint
	serverPath := filepath.Join(config.OutputDir, "internal/server/server.go")
	serverContent, err := os.ReadFile(serverPath)
	if err != nil {
		t.Fatalf("Failed to read server file: %v", err)
	}

	serverStr := string(serverContent)
	if !contains(serverStr, "/startup") {
		t.Error("Startup endpoint route not found in server configuration")
	}
}

func TestGenerator_TypeScriptClient(t *testing.T) {
	config := &config.ProjectConfig{
		Name:        "ts-client-test",
		Description: "Test TypeScript client generation",
		GoModule:    "github.com/example/ts-client-test",
		Tier:        config.TierBasic,
		Version:     "1.0.0",
		OutputDir:   "test-ts-client",
		Features: config.FeatureFlags{
			Kubernetes: true,
		},
	}

	// Clean up
	defer os.RemoveAll(config.OutputDir)

	generator, err := New(config)
	if err != nil {
		t.Fatalf("Failed to create generator: %v", err)
	}

	err = generator.Generate()
	if err != nil {
		t.Fatalf("Failed to generate project: %v", err)
	}

	// Verify TypeScript client includes startup method
	clientPath := filepath.Join(config.OutputDir, "client/typescript/src/client.ts")
	content, err := os.ReadFile(clientPath)
	if err != nil {
		t.Fatalf("Failed to read TypeScript client: %v", err)
	}

	clientContent := string(content)
	if !contains(clientContent, "checkStartup") {
		t.Error("checkStartup method not found in TypeScript client")
	}

	if !contains(clientContent, "/health/startup") {
		t.Error("Startup endpoint path not found in TypeScript client")
	}
}

func TestGenerator_Documentation(t *testing.T) {
	config := &config.ProjectConfig{
		Name:        "docs-test",
		Description: "Test documentation generation",
		GoModule:    "github.com/example/docs-test",
		Tier:        config.TierBasic,
		Version:     "1.0.0",
		OutputDir:   "test-docs",
		Features: config.FeatureFlags{
			Kubernetes: true,
		},
	}

	// Clean up
	defer os.RemoveAll(config.OutputDir)

	generator, err := New(config)
	if err != nil {
		t.Fatalf("Failed to create generator: %v", err)
	}

	err = generator.Generate()
	if err != nil {
		t.Fatalf("Failed to generate project: %v", err)
	}

	// Check README includes startup endpoint
	readmePath := filepath.Join(config.OutputDir, "README.md")
	readmeContent, err := os.ReadFile(readmePath)
	if err != nil {
		t.Fatalf("Failed to read README: %v", err)
	}

	if !contains(string(readmeContent), "/health/startup") {
		t.Error("Startup endpoint not documented in README")
	}

	// Check API documentation includes startup endpoint
	apiDocsPath := filepath.Join(config.OutputDir, "docs/API.md")
	apiContent, err := os.ReadFile(apiDocsPath)
	if err != nil {
		t.Fatalf("Failed to read API docs: %v", err)
	}

	if !contains(string(apiContent), "GET /health/startup") {
		t.Error("Startup endpoint not documented in API docs")
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && 
		   (s == substr || 
		    s[:len(substr)] == substr || 
		    s[len(s)-len(substr):] == substr ||
		    containsSubstring(s, substr))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
