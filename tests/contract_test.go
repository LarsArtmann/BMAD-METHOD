package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/LarsArtmann/BMAD-METHOD/pkg/config"
	"github.com/LarsArtmann/BMAD-METHOD/pkg/generator"
)

// ContractTestSuite runs comprehensive contract tests for generated projects
type ContractTestSuite struct {
	suite.Suite
	tempDir     string
	testProjects map[string]string // tier -> project path
	servers     map[string]*TestServer
}

// TestServer represents a running test server
type TestServer struct {
	cmd     *exec.Cmd
	port    int
	baseURL string
	tier    string
}

// SetupSuite prepares the test environment
func (suite *ContractTestSuite) SetupSuite() {
	var err error
	suite.tempDir, err = ioutil.TempDir("", "bmad-contract-tests")
	require.NoError(suite.T(), err)

	suite.testProjects = make(map[string]string)
	suite.servers = make(map[string]*TestServer)

	// Generate test projects for each tier
	tiers := []string{"basic", "intermediate", "advanced", "enterprise"}
	for i, tier := range tiers {
		suite.generateTestProject(tier, 8080+i)
	}
}

// TearDownSuite cleans up the test environment
func (suite *ContractTestSuite) TearDownSuite() {
	// Stop all servers
	for _, server := range suite.servers {
		if server.cmd != nil && server.cmd.Process != nil {
			server.cmd.Process.Kill()
			server.cmd.Wait()
		}
	}

	// Clean up temp directory
	if suite.tempDir != "" {
		os.RemoveAll(suite.tempDir)
	}
}

// generateTestProject generates and starts a test project for the given tier
func (suite *ContractTestSuite) generateTestProject(tier string, port int) {
	projectName := fmt.Sprintf("test-%s", tier)
	projectPath := filepath.Join(suite.tempDir, projectName)

	// Create project configuration
	cfg := &config.ProjectConfig{
		Name:      projectName,
		Tier:      config.TemplateTier(tier),
		OutputDir: projectPath,
		GoModule:  fmt.Sprintf("test/%s", projectName),
		Features: config.Features{
			TypeScript:    true,
			Docker:        true,
			Kubernetes:    tier != "basic",
			OpenTelemetry: tier == "advanced" || tier == "enterprise",
			CloudEvents:   tier == "enterprise",
			ServerTiming:  tier != "basic",
		},
	}

	if tier == "enterprise" {
		cfg.Features.Security = true
		cfg.Features.Compliance = true
	}

	// Generate project
	gen, err := generator.New(cfg)
	require.NoError(suite.T(), err)

	err = gen.Generate()
	require.NoError(suite.T(), err)

	suite.testProjects[tier] = projectPath

	// Build and start the project
	suite.buildAndStartProject(tier, projectPath, port)
}

// buildAndStartProject builds and starts a generated project
func (suite *ContractTestSuite) buildAndStartProject(tier, projectPath string, port int) {
	// Change to project directory
	originalDir, _ := os.Getwd()
	err := os.Chdir(projectPath)
	require.NoError(suite.T(), err)
	defer os.Chdir(originalDir)

	// Initialize go module
	cmd := exec.Command("go", "mod", "tidy")
	err = cmd.Run()
	require.NoError(suite.T(), err)

	// Build the project
	binaryPath := filepath.Join(projectPath, "server")
	cmd = exec.Command("go", "build", "-o", binaryPath, "./cmd/server")
	err = cmd.Run()
	require.NoError(suite.T(), err)

	// Start the server
	cmd = exec.Command(binaryPath)
	cmd.Env = append(os.Environ(), fmt.Sprintf("PORT=%d", port))
	err = cmd.Start()
	require.NoError(suite.T(), err)

	server := &TestServer{
		cmd:     cmd,
		port:    port,
		baseURL: fmt.Sprintf("http://localhost:%d", port),
		tier:    tier,
	}

	suite.servers[tier] = server

	// Wait for server to start
	suite.waitForServer(server.baseURL, 30*time.Second)
}

// waitForServer waits for a server to become available
func (suite *ContractTestSuite) waitForServer(baseURL string, timeout time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			suite.T().Fatalf("Server at %s did not start within timeout", baseURL)
		case <-ticker.C:
			resp, err := http.Get(baseURL + "/health")
			if err == nil && resp.StatusCode == http.StatusOK {
				resp.Body.Close()
				return
			}
			if resp != nil {
				resp.Body.Close()
			}
		}
	}
}

// TestBasicHealthEndpoints tests basic health endpoints for all tiers
func (suite *ContractTestSuite) TestBasicHealthEndpoints() {
	endpoints := []string{"/health", "/health/ready", "/health/live", "/health/startup"}

	for tier, server := range suite.servers {
		suite.T().Run(fmt.Sprintf("Tier_%s", tier), func(t *testing.T) {
			for _, endpoint := range endpoints {
				t.Run(strings.TrimPrefix(endpoint, "/"), func(t *testing.T) {
					resp, err := http.Get(server.baseURL + endpoint)
					require.NoError(t, err)
					defer resp.Body.Close()

					assert.Equal(t, http.StatusOK, resp.StatusCode)
					assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

					// Parse response body
					var healthReport map[string]interface{}
					err = json.NewDecoder(resp.Body).Decode(&healthReport)
					require.NoError(t, err)

					// Verify required fields
					assert.Contains(t, healthReport, "status")
					assert.Contains(t, healthReport, "timestamp")
					assert.Contains(t, healthReport, "version")
					assert.Equal(t, "healthy", healthReport["status"])
				})
			}
		})
	}
}

// TestServerTimeEndpoint tests server time endpoint
func (suite *ContractTestSuite) TestServerTimeEndpoint() {
	for tier, server := range suite.servers {
		suite.T().Run(fmt.Sprintf("Tier_%s", tier), func(t *testing.T) {
			resp, err := http.Get(server.baseURL + "/health/time")
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusOK, resp.StatusCode)

			var serverTime map[string]interface{}
			err = json.NewDecoder(resp.Body).Decode(&serverTime)
			require.NoError(t, err)

			// Verify time fields
			assert.Contains(t, serverTime, "timestamp")
			assert.Contains(t, serverTime, "timezone")
			assert.Contains(t, serverTime, "unix")
		})
	}
}

// TestTierSpecificFeatures tests features specific to each tier
func (suite *ContractTestSuite) TestTierSpecificFeatures() {
	for tier, server := range suite.servers {
		suite.T().Run(fmt.Sprintf("Tier_%s", tier), func(t *testing.T) {
			switch tier {
			case "intermediate", "advanced", "enterprise":
				// Test dependencies endpoint
				resp, err := http.Get(server.baseURL + "/health/dependencies")
				if err == nil {
					defer resp.Body.Close()
					assert.True(t, resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusServiceUnavailable)
				}

			case "advanced", "enterprise":
				// Test metrics endpoint (if available)
				resp, err := http.Get(server.baseURL + "/health/metrics")
				if err == nil {
					defer resp.Body.Close()
					// Metrics might not be available on all advanced tiers
					assert.True(t, resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNotFound)
				}

			case "enterprise":
				// Enterprise-specific tests would go here
				// For now, just verify the server is running
				resp, err := http.Get(server.baseURL + "/health")
				require.NoError(t, err)
				defer resp.Body.Close()
				assert.Equal(t, http.StatusOK, resp.StatusCode)
			}
		})
	}
}

// TestPerformanceBenchmarks runs performance tests on generated projects
func (suite *ContractTestSuite) TestPerformanceBenchmarks() {
	for tier, server := range suite.servers {
		suite.T().Run(fmt.Sprintf("Performance_Tier_%s", tier), func(t *testing.T) {
			// Test response time
			start := time.Now()
			resp, err := http.Get(server.baseURL + "/health")
			duration := time.Since(start)

			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusOK, resp.StatusCode)
			assert.Less(t, duration, 100*time.Millisecond, "Health endpoint should respond within 100ms")

			// Test concurrent requests
			suite.testConcurrentRequests(t, server.baseURL+"/health", 10, 1*time.Second)
		})
	}
}

// testConcurrentRequests tests handling of concurrent requests
func (suite *ContractTestSuite) testConcurrentRequests(t *testing.T, url string, concurrency int, timeout time.Duration) {
	type result struct {
		statusCode int
		duration   time.Duration
		err        error
	}

	results := make(chan result, concurrency)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Start concurrent requests
	for i := 0; i < concurrency; i++ {
		go func() {
			start := time.Now()
			req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
			if err != nil {
				results <- result{err: err}
				return
			}

			resp, err := http.DefaultClient.Do(req)
			duration := time.Since(start)
			
			if err != nil {
				results <- result{err: err, duration: duration}
				return
			}
			defer resp.Body.Close()

			results <- result{statusCode: resp.StatusCode, duration: duration}
		}()
	}

	// Collect results
	successCount := 0
	for i := 0; i < concurrency; i++ {
		select {
		case res := <-results:
			if res.err == nil && res.statusCode == http.StatusOK {
				successCount++
			}
		case <-ctx.Done():
			t.Fatal("Concurrent requests test timed out")
		}
	}

	assert.Equal(t, concurrency, successCount, "All concurrent requests should succeed")
}

// TestGeneratedProjectStructure validates the structure of generated projects
func (suite *ContractTestSuite) TestGeneratedProjectStructure() {
	expectedFiles := map[string][]string{
		"basic": {
			"README.md", "go.mod", "Makefile", "Dockerfile",
			"cmd/server/main.go",
			"internal/handlers/health.go",
			"internal/server/server.go",
			"internal/config/config.go",
		},
		"intermediate": {
			"README.md", "go.mod", "Makefile", "Dockerfile",
			"cmd/server/main.go",
			"internal/handlers/health.go",
			"internal/handlers/dependencies.go",
			"internal/server/server.go",
			"internal/config/config.go",
			"deployments/kubernetes/deployment.yaml",
		},
		"advanced": {
			"README.md", "go.mod", "Makefile", "Dockerfile",
			"cmd/server/main.go",
			"internal/handlers/health.go",
			"internal/observability/tracing.go",
			"internal/observability/metrics.go",
			"deployments/kubernetes/deployment.yaml",
		},
		"enterprise": {
			"README.md", "go.mod", "Makefile", "Dockerfile",
			"cmd/server/main.go",
			"internal/handlers/health.go",
			"internal/security/mtls.go",
			"internal/security/rbac.go",
			"internal/compliance/audit.go",
			"deployments/kubernetes/deployment.yaml",
		},
	}

	for tier, files := range expectedFiles {
		projectPath := suite.testProjects[tier]
		suite.T().Run(fmt.Sprintf("Structure_Tier_%s", tier), func(t *testing.T) {
			for _, file := range files {
				filePath := filepath.Join(projectPath, file)
				_, err := os.Stat(filePath)
				assert.NoError(t, err, "File %s should exist in %s tier", file, tier)
			}
		})
	}
}

// TestTypeScriptClientGeneration tests TypeScript client generation
func (suite *ContractTestSuite) TestTypeScriptClientGeneration() {
	for tier, projectPath := range suite.testProjects {
		suite.T().Run(fmt.Sprintf("TypeScript_Tier_%s", tier), func(t *testing.T) {
			clientPath := filepath.Join(projectPath, "client", "typescript")
			
			// Check if TypeScript files exist
			expectedFiles := []string{
				"src/client.ts",
				"src/types.ts",
				"package.json",
				"tsconfig.json",
			}

			for _, file := range expectedFiles {
				filePath := filepath.Join(clientPath, file)
				_, err := os.Stat(filePath)
				assert.NoError(t, err, "TypeScript file %s should exist", file)
			}

			// Try to compile TypeScript (if npm is available)
			if _, err := exec.LookPath("npm"); err == nil {
				originalDir, _ := os.Getwd()
				os.Chdir(clientPath)
				defer os.Chdir(originalDir)

				// Install dependencies and build
				cmd := exec.Command("npm", "install")
				err := cmd.Run()
				if err == nil {
					cmd = exec.Command("npm", "run", "build")
					err = cmd.Run()
					assert.NoError(t, err, "TypeScript client should compile successfully")
				}
			}
		})
	}
}

// TestContractTestSuite runs the complete contract test suite
func TestContractTestSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping contract tests in short mode")
	}

	suite.Run(t, new(ContractTestSuite))
}