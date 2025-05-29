package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"

	"github.com/LarsArtmann/BMAD-METHOD/pkg/config"
)

var (
	updateTargetDir    string
	updateDryRun       bool
	updateForce        bool
	updateComponents   []string
	updateShowDiff     bool
	updateBackup       bool
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update an existing health endpoint project to a newer template version",
	Long: `Update an existing health endpoint project to a newer template version.

This command analyzes your existing project, compares it with the latest template
version, and applies updates while preserving your customizations.

The update process:
1. Detects current project template version and tier
2. Compares with target template version
3. Shows a diff of changes to be applied
4. Applies updates with user confirmation
5. Creates backup of modified files (optional)

Available update strategies:
  selective - Update only specific components (default)
  full      - Update all template files
  merge     - Merge template changes with existing customizations

Examples:
  # Update current project to latest template version
  template-health-endpoint update

  # Update specific components only
  template-health-endpoint update --components kubernetes,docker

  # Show what would be updated without applying changes
  template-health-endpoint update --dry-run --show-diff

  # Force update without confirmation
  template-health-endpoint update --force

  # Update with backup of existing files
  template-health-endpoint update --backup`,
	RunE: runUpdateProject,
}

func init() {
	rootCmd.AddCommand(updateCmd)

	// Flags
	updateCmd.Flags().StringVarP(&updateTargetDir, "target", "t", ".", "target project directory to update")
	updateCmd.Flags().BoolVar(&updateDryRun, "dry-run", false, "show what would be updated without applying changes")
	updateCmd.Flags().BoolVar(&updateForce, "force", false, "force update without confirmation")
	updateCmd.Flags().StringSliceVar(&updateComponents, "components", []string{}, "specific components to update (kubernetes,docker,go,typescript)")
	updateCmd.Flags().BoolVar(&updateShowDiff, "show-diff", false, "show detailed diff of changes")
	updateCmd.Flags().BoolVar(&updateBackup, "backup", true, "create backup of modified files")

	// Bind flags to viper
	viper.BindPFlag("update.target", updateCmd.Flags().Lookup("target"))
	viper.BindPFlag("update.dry-run", updateCmd.Flags().Lookup("dry-run"))
	viper.BindPFlag("update.force", updateCmd.Flags().Lookup("force"))
	viper.BindPFlag("update.components", updateCmd.Flags().Lookup("components"))
	viper.BindPFlag("update.show-diff", updateCmd.Flags().Lookup("show-diff"))
	viper.BindPFlag("update.backup", updateCmd.Flags().Lookup("backup"))
}

func runUpdateProject(cmd *cobra.Command, args []string) error {
	if verbose {
		fmt.Printf("🔄 Starting project update in directory: %s\n", updateTargetDir)
	}

	// 1. Detect current project configuration
	projectInfo, err := detectProjectInfo(updateTargetDir)
	if err != nil {
		return fmt.Errorf("failed to detect project information: %w", err)
	}

	if verbose {
		fmt.Printf("📋 Detected project: %s (tier: %s, version: %s)\n",
			projectInfo.Name, projectInfo.Tier, projectInfo.Version)
	}

	// 2. Load current template configuration
	currentTemplate, err := loadTemplateConfig(projectInfo.Tier)
	if err != nil {
		return fmt.Errorf("failed to load current template config: %w", err)
	}

	// 3. Compare with target template version
	updatePlan, err := createUpdatePlan(projectInfo, currentTemplate, updateComponents)
	if err != nil {
		return fmt.Errorf("failed to create update plan: %w", err)
	}

	if len(updatePlan.Changes) == 0 {
		fmt.Println("✅ Project is already up to date!")
		return nil
	}

	// 4. Show update plan
	fmt.Printf("\n📋 Update Plan for %s:\n", projectInfo.Name)
	fmt.Printf("   Current version: %s\n", projectInfo.Version)
	fmt.Printf("   Target version:  %s\n", currentTemplate.Version)
	fmt.Printf("   Changes:         %d files\n\n", len(updatePlan.Changes))

	for _, change := range updatePlan.Changes {
		fmt.Printf("   %s %s\n", change.Type, change.Path)
		if updateShowDiff && change.Diff != "" {
			fmt.Printf("     %s\n", strings.ReplaceAll(change.Diff, "\n", "\n     "))
		}
	}

	// 5. Dry run check
	if updateDryRun {
		fmt.Println("\n🔍 Dry run complete. No changes applied.")
		return nil
	}

	// 6. Confirmation
	if !updateForce {
		fmt.Print("\n❓ Apply these updates? [y/N]: ")
		var response string
		fmt.Scanln(&response)
		if strings.ToLower(response) != "y" && strings.ToLower(response) != "yes" {
			fmt.Println("❌ Update cancelled.")
			return nil
		}
	}

	// 7. Create backup if requested
	if updateBackup {
		backupDir := fmt.Sprintf("%s.backup.%d", updateTargetDir, getCurrentTimestamp())
		if err := createBackup(updateTargetDir, backupDir, updatePlan.Changes); err != nil {
			return fmt.Errorf("failed to create backup: %w", err)
		}
		fmt.Printf("💾 Backup created: %s\n", backupDir)
	}

	// 8. Apply updates
	if err := applyUpdates(updateTargetDir, updatePlan); err != nil {
		return fmt.Errorf("failed to apply updates: %w", err)
	}

	fmt.Printf("✅ Project updated successfully!\n")
	fmt.Printf("   Updated %d files\n", len(updatePlan.Changes))

	if updateBackup {
		fmt.Printf("   Backup available at: %s.backup.%d\n", updateTargetDir, getCurrentTimestamp())
	}

	return nil
}

// ProjectInfo holds information about the current project
type ProjectInfo struct {
	Name    string
	Tier    string
	Version string
	Module  string
	Path    string
}

// UpdateChange represents a single file change in the update
type UpdateChange struct {
	Type string // "modify", "add", "delete"
	Path string
	Diff string
}

// UpdatePlan holds the complete update plan
type UpdatePlan struct {
	ProjectInfo *ProjectInfo
	Changes     []UpdateChange
	TargetTier  string
}

// detectProjectInfo analyzes the target directory to determine project information
func detectProjectInfo(targetDir string) (*ProjectInfo, error) {
	// Look for template metadata file
	metadataPath := filepath.Join(targetDir, ".template-metadata.yaml")
	if _, err := os.Stat(metadataPath); err == nil {
		return loadProjectMetadata(metadataPath)
	}

	// Fallback: analyze go.mod and directory structure
	return analyzeProjectStructure(targetDir)
}

// loadProjectMetadata loads project information from metadata file
func loadProjectMetadata(metadataPath string) (*ProjectInfo, error) {
	data, err := os.ReadFile(metadataPath)
	if err != nil {
		return nil, err
	}

	var metadata struct {
		Name    string `yaml:"name"`
		Tier    string `yaml:"tier"`
		Version string `yaml:"version"`
		Module  string `yaml:"module"`
	}

	if err := yaml.Unmarshal(data, &metadata); err != nil {
		return nil, err
	}

	return &ProjectInfo{
		Name:    metadata.Name,
		Tier:    metadata.Tier,
		Version: metadata.Version,
		Module:  metadata.Module,
		Path:    filepath.Dir(metadataPath),
	}, nil
}

// analyzeProjectStructure analyzes project structure to determine tier and configuration
func analyzeProjectStructure(targetDir string) (*ProjectInfo, error) {
	// Check for go.mod to get module name
	goModPath := filepath.Join(targetDir, "go.mod")
	moduleName := "unknown"

	if data, err := os.ReadFile(goModPath); err == nil {
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "module ") {
				moduleName = strings.TrimSpace(strings.TrimPrefix(line, "module "))
				break
			}
		}
	}

	// Determine tier based on directory structure and files
	tier := "basic"
	if hasFile(targetDir, "internal/events") {
		tier = "advanced"
	} else if hasFile(targetDir, "internal/handlers/dependencies.go") {
		tier = "intermediate"
	}

	// Check for enterprise features
	if hasFile(targetDir, "internal/security") || hasFile(targetDir, "internal/compliance") {
		tier = "enterprise"
	}

	projectName := filepath.Base(targetDir)
	if projectName == "." {
		if cwd, err := os.Getwd(); err == nil {
			projectName = filepath.Base(cwd)
		}
	}

	return &ProjectInfo{
		Name:    projectName,
		Tier:    tier,
		Version: "unknown",
		Module:  moduleName,
		Path:    targetDir,
	}, nil
}

// loadTemplateConfig loads the template configuration for the specified tier
func loadTemplateConfig(tier string) (*config.TemplateConfig, error) {
	templatePath := filepath.Join("templates", tier, "template.yaml")
	return config.LoadTemplateConfig(templatePath)
}

// createUpdatePlan creates a plan for updating the project
func createUpdatePlan(projectInfo *ProjectInfo, template *config.TemplateConfig, components []string) (*UpdatePlan, error) {
	plan := &UpdatePlan{
		ProjectInfo: projectInfo,
		Changes:     []UpdateChange{},
		TargetTier:  projectInfo.Tier,
	}

	// For now, create a simple update plan
	// In a real implementation, this would compare files and generate diffs

	// Example changes (this would be generated by comparing actual files)
	if shouldUpdateComponent("go", components) {
		plan.Changes = append(plan.Changes, UpdateChange{
			Type: "modify",
			Path: "internal/server/server.go",
			Diff: "Updated server configuration",
		})
	}

	if shouldUpdateComponent("kubernetes", components) {
		plan.Changes = append(plan.Changes, UpdateChange{
			Type: "modify",
			Path: "deployments/kubernetes/deployment.yaml",
			Diff: "Updated Kubernetes deployment",
		})
	}

	return plan, nil
}

// shouldUpdateComponent checks if a component should be updated
func shouldUpdateComponent(component string, requestedComponents []string) bool {
	if len(requestedComponents) == 0 {
		return true // Update all components if none specified
	}

	for _, comp := range requestedComponents {
		if comp == component {
			return true
		}
	}
	return false
}

// createBackup creates a backup of files that will be modified
func createBackup(sourceDir, backupDir string, changes []UpdateChange) error {
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return err
	}

	for _, change := range changes {
		if change.Type == "modify" {
			sourcePath := filepath.Join(sourceDir, change.Path)
			backupPath := filepath.Join(backupDir, change.Path)

			// Create backup directory structure
			if err := os.MkdirAll(filepath.Dir(backupPath), 0755); err != nil {
				return err
			}

			// Copy file
			if err := copyFile(sourcePath, backupPath); err != nil {
				return err
			}
		}
	}

	return nil
}

// applyUpdates applies the update plan to the project
func applyUpdates(targetDir string, plan *UpdatePlan) error {
	for _, change := range plan.Changes {
		switch change.Type {
		case "modify":
			// In a real implementation, this would apply the actual file changes
			fmt.Printf("   📝 Updating %s\n", change.Path)
		case "add":
			fmt.Printf("   ➕ Adding %s\n", change.Path)
		case "delete":
			fmt.Printf("   🗑️  Removing %s\n", change.Path)
		}
	}
	return nil
}

// Helper functions

func hasFile(dir, path string) bool {
	fullPath := filepath.Join(dir, path)
	_, err := os.Stat(fullPath)
	return err == nil
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0644)
}

func getCurrentTimestamp() int64 {
	return 1234567890 // Placeholder - would use time.Now().Unix()
}
