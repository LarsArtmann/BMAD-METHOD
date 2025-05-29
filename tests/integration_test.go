package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/LarsArtmann/BMAD-METHOD/pkg/config"
	"github.com/LarsArtmann/BMAD-METHOD/pkg/generator"
)

// HealthReport represents the health response structure
type HealthReport struct {
	Status      string    `json:"status"`
	Timestamp   time.Time `json:"timestamp"`
	Version     string    `json:"version"`
	Uptime      int64     `json:"uptime"`
	UptimeHuman string    `json:"uptime_human"`
}

// ServerTime represents the server time response structure
type ServerTime struct {
	Timestamp string `json:"timestamp"`
	Timezone  string `json:"timezone"`
	Unix      int64  `json:"unix"`
	UnixMilli int64  `json:"unix_milli"`
	ISO8601   string `json:"iso8601"`
	Formatted string `json:"formatted"`
}

func TestEndToEndGeneration(t *testing.T) {
	projectName := "integration-test-service"
	outputDir := "test-integration-output"

	// Clean up
	defer os.RemoveAll(outputDir)

	// Create project configuration
	config := &config.ProjectConfig{
		Name:        projectName,
		Description: "Integration test health service",
		GoModule:    "github.com/example/" + projectName,
		Tier:        config.TierBasic,
		Version:     "1.0.0",
		OutputDir:   outputDir,
		Features: config.FeatureFlags{
			Kubernetes: true,
		},
	}

	// Generate project
	gen, err := generator.New(config)
	if err != nil {
		t.Fatalf("Failed to create generator: %v", err)
	}

	err = gen.Generate()
	if err != nil {
		t.Fatalf("Failed to generate project: %v", err)
	}

	// Test that project compiles
	err = testProjectCompilation(outputDir)
	if err != nil {
		t.Fatalf("Generated project failed to compile: %v", err)
	}

	// Test that all endpoints work
	err = testHealthEndpoints(outputDir, projectName)
	if err != nil {
		t.Fatalf("Health endpoints test failed: %v", err)
	}
}

func testProjectCompilation(projectDir string) error {
	// Change to project directory
	originalDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}
	defer os.Chdir(originalDir)

	err = os.Chdir(projectDir)
	if err != nil {
		return fmt.Errorf("failed to change to project directory: %w", err)
	}

	// Run go mod tidy
	cmd := exec.Command("go", "mod", "tidy")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go mod tidy failed: %w", err)
	}

	// Build the project
	cmd = exec.Command("go", "build", "-o", "bin/test-service", "cmd/server/main.go")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("go build failed: %w", err)
	}

	return nil
}

func testHealthEndpoints(projectDir, projectName string) error {
	// Change to project directory
	originalDir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %w", err)
	}
	defer os.Chdir(originalDir)

	err = os.Chdir(projectDir)
	if err != nil {
		return fmt.Errorf("failed to change to project directory: %w", err)
	}

	// Start the server
	cmd := exec.Command("./bin/test-service")
	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}
	defer cmd.Process.Kill()

	// Wait for server to start
	time.Sleep(2 * time.Second)

	// Test all health endpoints
	endpoints := []struct {
		path     string
		expected string
	}{
		{"/health", "healthy"},
		{"/health/ready", "healthy"},
		{"/health/live", "healthy"},
		{"/health/startup", "healthy"}, // This is our new endpoint
	}

	for _, endpoint := range endpoints {
		err := testHealthEndpoint("http://localhost:8080"+endpoint.path, endpoint.expected)
		if err != nil {
			return fmt.Errorf("endpoint %s failed: %w", endpoint.path, err)
		}
	}

	// Test server time endpoint
	err = testServerTimeEndpoint("http://localhost:8080/health/time")
	if err != nil {
		return fmt.Errorf("server time endpoint failed: %w", err)
	}

	return nil
}

func testHealthEndpoint(url, expectedStatus string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	var health HealthReport
	err = json.NewDecoder(resp.Body).Decode(&health)
	if err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	if health.Status != expectedStatus {
		return fmt.Errorf("expected status %s, got %s", expectedStatus, health.Status)
	}

	// Verify required fields are present
	if health.Version == "" {
		return fmt.Errorf("version field is empty")
	}

	if health.UptimeHuman == "" {
		return fmt.Errorf("uptime_human field is empty")
	}

	return nil
}

func testServerTimeEndpoint(url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	var serverTime ServerTime
	err = json.NewDecoder(resp.Body).Decode(&serverTime)
	if err != nil {
		return fmt.Errorf("failed to decode response: %w", err)
	}

	// Verify required fields are present
	if serverTime.Timestamp == "" {
		return fmt.Errorf("timestamp field is empty")
	}

	if serverTime.Timezone == "" {
		return fmt.Errorf("timezone field is empty")
	}

	if serverTime.Unix == 0 {
		return fmt.Errorf("unix timestamp is zero")
	}

	if serverTime.ISO8601 == "" {
		return fmt.Errorf("iso8601 field is empty")
	}

	return nil
}

func TestCLICommands(t *testing.T) {
	// Test CLI help command
	cmd := exec.Command("./bin/template-health-endpoint", "--help")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("CLI help command failed: %v", err)
	}

	helpOutput := string(output)
	if !contains(helpOutput, "generate") {
		t.Error("CLI help should mention 'generate' command")
	}

	if !contains(helpOutput, "validate") {
		t.Error("CLI help should mention 'validate' command")
	}
}

func TestValidateCommand(t *testing.T) {
	// Test TypeSpec validation
	cmd := exec.Command("./bin/template-health-endpoint", "validate", "--verbose")
	err := cmd.Run()
	if err != nil {
		t.Fatalf("TypeSpec validation failed: %v", err)
	}
}

func TestGenerateCommandDryRun(t *testing.T) {
	// Test dry run generation
	cmd := exec.Command("./bin/template-health-endpoint", "generate", 
		"--name", "dry-run-test", 
		"--tier", "basic", 
		"--module", "github.com/example/dry-run-test",
		"--dry-run")
	
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("Dry run generation failed: %v", err)
	}

	dryRunOutput := string(output)
	if !contains(dryRunOutput, "dry run") || !contains(dryRunOutput, "would generate") {
		t.Error("Dry run should indicate it's a simulation")
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
