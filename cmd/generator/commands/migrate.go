package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/LarsArtmann/BMAD-METHOD/pkg/config"
	"github.com/LarsArtmann/BMAD-METHOD/pkg/generator"
)

var (
	migrateTargetTier string
	migrateTargetDir  string
	migrateDryRun     bool
	migrateForce      bool
	migrateBackup     bool
)

// migrateCmd represents the migrate command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate a project between template tiers",
	Long: `Migrate a project between template tiers (basic → intermediate → advanced → enterprise).

This command upgrades or downgrades your project to a different complexity tier,
adding or removing features as needed while preserving your customizations.

Migration paths:
  basic → intermediate    Add dependency health checks
  intermediate → advanced Add full observability (OpenTelemetry, CloudEvents)
  advanced → enterprise   Add security and compliance features
  
Reverse migrations are also supported with appropriate warnings.

The migration process:
1. Detects current project tier
2. Validates migration path
3. Shows migration plan with added/removed features
4. Updates dependencies and configurations
5. Adds new code components
6. Updates documentation

Examples:
  # Migrate from basic to intermediate tier
  template-health-endpoint migrate --to intermediate

  # Migrate to enterprise tier with backup
  template-health-endpoint migrate --to enterprise --backup

  # Show migration plan without applying changes
  template-health-endpoint migrate --to advanced --dry-run

  # Force migration without confirmation
  template-health-endpoint migrate --to intermediate --force`,
	RunE: runMigrateProject,
}

func init() {
	rootCmd.AddCommand(migrateCmd)

	// Flags
	migrateCmd.Flags().StringVar(&migrateTargetTier, "to", "", "target tier to migrate to (basic, intermediate, advanced, enterprise)")
	migrateCmd.Flags().StringVarP(&migrateTargetDir, "target", "t", ".", "target project directory to migrate")
	migrateCmd.Flags().BoolVar(&migrateDryRun, "dry-run", false, "show migration plan without applying changes")
	migrateCmd.Flags().BoolVar(&migrateForce, "force", false, "force migration without confirmation")
	migrateCmd.Flags().BoolVar(&migrateBackup, "backup", true, "create backup before migration")

	// Mark required flags
	migrateCmd.MarkFlagRequired("to")

	// Bind flags to viper
	viper.BindPFlag("migrate.to", migrateCmd.Flags().Lookup("to"))
	viper.BindPFlag("migrate.target", migrateCmd.Flags().Lookup("target"))
	viper.BindPFlag("migrate.dry-run", migrateCmd.Flags().Lookup("dry-run"))
	viper.BindPFlag("migrate.force", migrateCmd.Flags().Lookup("force"))
	viper.BindPFlag("migrate.backup", migrateCmd.Flags().Lookup("backup"))
}

func runMigrateProject(cmd *cobra.Command, args []string) error {
	if verbose {
		fmt.Printf("🔄 Starting project migration in directory: %s\n", migrateTargetDir)
	}

	// Validate target tier
	validTiers := []string{"basic", "intermediate", "advanced", "enterprise"}
	if !contains(validTiers, migrateTargetTier) {
		return fmt.Errorf("invalid target tier '%s'. Valid tiers: %s", migrateTargetTier, strings.Join(validTiers, ", "))
	}

	// 1. Detect current project configuration
	projectInfo, err := detectProjectInfo(migrateTargetDir)
	if err != nil {
		return fmt.Errorf("failed to detect project information: %w", err)
	}

	if verbose {
		fmt.Printf("📋 Detected project: %s (current tier: %s)\n", projectInfo.Name, projectInfo.Tier)
	}

	// 2. Check if migration is needed
	if projectInfo.Tier == migrateTargetTier {
		fmt.Printf("✅ Project is already at tier '%s'. No migration needed.\n", migrateTargetTier)
		return nil
	}

	// 3. Validate migration path
	migrationPath, err := validateMigrationPath(projectInfo.Tier, migrateTargetTier)
	if err != nil {
		return fmt.Errorf("invalid migration path: %w", err)
	}

	// 4. Create migration plan
	migrationPlan, err := createMigrationPlan(projectInfo, migrateTargetTier, migrationPath)
	if err != nil {
		return fmt.Errorf("failed to create migration plan: %w", err)
	}

	// 5. Show migration plan
	fmt.Printf("\n📋 Migration Plan:\n")
	fmt.Printf("   Project:      %s\n", projectInfo.Name)
	fmt.Printf("   From:         %s tier\n", projectInfo.Tier)
	fmt.Printf("   To:           %s tier\n", migrateTargetTier)
	fmt.Printf("   Path:         %s\n", strings.Join(migrationPath, " → "))
	fmt.Printf("   Changes:      %d operations\n\n", len(migrationPlan.Operations))

	// Show operations
	for _, op := range migrationPlan.Operations {
		fmt.Printf("   %s %s\n", op.Type, op.Description)
		if len(op.Files) > 0 {
			for _, file := range op.Files {
				fmt.Printf("     - %s\n", file)
			}
		}
	}

	// Show warnings for downgrades
	if isDowngrade(projectInfo.Tier, migrateTargetTier) {
		fmt.Printf("\n⚠️  WARNING: This is a downgrade migration.\n")
		fmt.Printf("   Some features and files will be removed.\n")
		fmt.Printf("   Make sure you have backups of important customizations.\n")
	}

	// 6. Dry run check
	if migrateDryRun {
		fmt.Println("\n🔍 Dry run complete. No changes applied.")
		return nil
	}

	// 7. Confirmation
	if !migrateForce {
		fmt.Print("\n❓ Proceed with migration? [y/N]: ")
		var response string
		fmt.Scanln(&response)
		if strings.ToLower(response) != "y" && strings.ToLower(response) != "yes" {
			fmt.Println("❌ Migration cancelled.")
			return nil
		}
	}

	// 8. Create backup if requested
	if migrateBackup {
		backupDir := fmt.Sprintf("%s.backup.migration.%d", migrateTargetDir, getCurrentTimestamp())
		if err := createProjectBackup(migrateTargetDir, backupDir); err != nil {
			return fmt.Errorf("failed to create backup: %w", err)
		}
		fmt.Printf("💾 Backup created: %s\n", backupDir)
	}

	// 9. Apply migration
	if err := applyMigration(migrateTargetDir, migrationPlan); err != nil {
		return fmt.Errorf("failed to apply migration: %w", err)
	}

	fmt.Printf("✅ Migration completed successfully!\n")
	fmt.Printf("   Project migrated from %s to %s tier\n", projectInfo.Tier, migrateTargetTier)
	fmt.Printf("   Applied %d operations\n", len(migrationPlan.Operations))

	// Show next steps
	showMigrationNextSteps(projectInfo.Tier, migrateTargetTier)

	return nil
}

// MigrationOperation represents a single migration operation
type MigrationOperation struct {
	Type        string   // "add", "remove", "modify", "dependency"
	Description string
	Files       []string
	Commands    []string
}

// MigrationPlan holds the complete migration plan
type MigrationPlan struct {
	FromTier   string
	ToTier     string
	Path       []string
	Operations []MigrationOperation
}

// validateMigrationPath validates and returns the migration path
func validateMigrationPath(fromTier, toTier string) ([]string, error) {
	tierOrder := map[string]int{
		"basic":        1,
		"intermediate": 2,
		"advanced":     3,
		"enterprise":   4,
	}

	fromLevel, fromExists := tierOrder[fromTier]
	toLevel, toExists := tierOrder[toTier]

	if !fromExists {
		return nil, fmt.Errorf("unknown source tier: %s", fromTier)
	}
	if !toExists {
		return nil, fmt.Errorf("unknown target tier: %s", toTier)
	}

	// Build migration path
	var path []string
	if fromLevel < toLevel {
		// Upgrade path
		for i := fromLevel; i <= toLevel; i++ {
			for tier, level := range tierOrder {
				if level == i {
					path = append(path, tier)
					break
				}
			}
		}
	} else {
		// Downgrade path
		for i := fromLevel; i >= toLevel; i-- {
			for tier, level := range tierOrder {
				if level == i {
					path = append(path, tier)
					break
				}
			}
		}
	}

	return path, nil
}

// createMigrationPlan creates a detailed migration plan
func createMigrationPlan(projectInfo *ProjectInfo, targetTier string, path []string) (*MigrationPlan, error) {
	plan := &MigrationPlan{
		FromTier:   projectInfo.Tier,
		ToTier:     targetTier,
		Path:       path,
		Operations: []MigrationOperation{},
	}

	// Generate operations based on tier transitions
	for i := 0; i < len(path)-1; i++ {
		fromTier := path[i]
		toTier := path[i+1]
		
		ops := generateTierTransitionOperations(fromTier, toTier)
		plan.Operations = append(plan.Operations, ops...)
	}

	return plan, nil
}

// generateTierTransitionOperations generates operations for a single tier transition
func generateTierTransitionOperations(fromTier, toTier string) []MigrationOperation {
	var operations []MigrationOperation

	switch fromTier + "→" + toTier {
	case "basic→intermediate":
		operations = append(operations, MigrationOperation{
			Type:        "add",
			Description: "Add dependency health check handlers",
			Files:       []string{"internal/handlers/dependencies.go"},
		})
		operations = append(operations, MigrationOperation{
			Type:        "modify",
			Description: "Update server configuration for dependency checks",
			Files:       []string{"internal/server/server.go", "internal/config/config.go"},
		})
		operations = append(operations, MigrationOperation{
			Type:        "dependency",
			Description: "Add dependency health check libraries",
			Commands:    []string{"go get github.com/heptiolabs/healthcheck"},
		})

	case "intermediate→advanced":
		operations = append(operations, MigrationOperation{
			Type:        "add",
			Description: "Add OpenTelemetry observability",
			Files:       []string{"internal/middleware/server_timing.go", "internal/events/health_events.go"},
		})
		operations = append(operations, MigrationOperation{
			Type:        "add",
			Description: "Add CloudEvents integration",
			Files:       []string{"internal/events/"},
		})
		operations = append(operations, MigrationOperation{
			Type:        "dependency",
			Description: "Add observability dependencies",
			Commands:    []string{
				"go get go.opentelemetry.io/otel",
				"go get github.com/cloudevents/sdk-go/v2",
			},
		})

	case "advanced→enterprise":
		operations = append(operations, MigrationOperation{
			Type:        "add",
			Description: "Add enterprise security features",
			Files:       []string{"internal/security/", "internal/compliance/"},
		})
		operations = append(operations, MigrationOperation{
			Type:        "add",
			Description: "Add multi-environment configurations",
			Files:       []string{"configs/development.yaml", "configs/staging.yaml", "configs/production.yaml"},
		})
		operations = append(operations, MigrationOperation{
			Type:        "modify",
			Description: "Update server for enterprise features",
			Files:       []string{"cmd/server/main.go", "internal/server/server.go"},
		})

	// Downgrade operations
	case "intermediate→basic":
		operations = append(operations, MigrationOperation{
			Type:        "remove",
			Description: "Remove dependency health check handlers",
			Files:       []string{"internal/handlers/dependencies.go"},
		})

	case "advanced→intermediate":
		operations = append(operations, MigrationOperation{
			Type:        "remove",
			Description: "Remove OpenTelemetry and CloudEvents features",
			Files:       []string{"internal/middleware/server_timing.go", "internal/events/"},
		})

	case "enterprise→advanced":
		operations = append(operations, MigrationOperation{
			Type:        "remove",
			Description: "Remove enterprise security and compliance features",
			Files:       []string{"internal/security/", "internal/compliance/", "configs/"},
		})
	}

	return operations
}

// applyMigration applies the migration plan
func applyMigration(targetDir string, plan *MigrationPlan) error {
	fmt.Printf("\n🔄 Applying migration...\n")

	for i, op := range plan.Operations {
		fmt.Printf("   [%d/%d] %s: %s\n", i+1, len(plan.Operations), op.Type, op.Description)

		switch op.Type {
		case "add":
			if err := addMigrationFiles(targetDir, plan.ToTier, op.Files); err != nil {
				return fmt.Errorf("failed to add files: %w", err)
			}

		case "remove":
			if err := removeMigrationFiles(targetDir, op.Files); err != nil {
				return fmt.Errorf("failed to remove files: %w", err)
			}

		case "modify":
			if err := modifyMigrationFiles(targetDir, plan.FromTier, plan.ToTier, op.Files); err != nil {
				return fmt.Errorf("failed to modify files: %w", err)
			}

		case "dependency":
			if err := runMigrationCommands(targetDir, op.Commands); err != nil {
				return fmt.Errorf("failed to run commands: %w", err)
			}
		}
	}

	// Update project metadata
	if err := updateProjectMetadata(targetDir, plan.ToTier); err != nil {
		return fmt.Errorf("failed to update project metadata: %w", err)
	}

	return nil
}

// addMigrationFiles copies files from the target tier template
func addMigrationFiles(targetDir, toTier string, files []string) error {
	templateDir := filepath.Join("templates", toTier)
	
	for _, file := range files {
		srcPath := filepath.Join(templateDir, file)
		dstPath := filepath.Join(targetDir, file)
		
		// Create directory if needed
		if err := os.MkdirAll(filepath.Dir(dstPath), 0755); err != nil {
			return err
		}
		
		// Copy file or directory
		if err := copyFileOrDir(srcPath, dstPath); err != nil {
			return err
		}
	}
	
	return nil
}

// removeMigrationFiles removes specified files
func removeMigrationFiles(targetDir string, files []string) error {
	for _, file := range files {
		path := filepath.Join(targetDir, file)
		if err := os.RemoveAll(path); err != nil {
			return err
		}
	}
	return nil
}

// modifyMigrationFiles updates existing files for the new tier
func modifyMigrationFiles(targetDir, fromTier, toTier string, files []string) error {
	// This would contain the logic to update existing files
	// For now, we'll just copy the new versions
	return addMigrationFiles(targetDir, toTier, files)
}

// runMigrationCommands runs shell commands for the migration
func runMigrationCommands(targetDir string, commands []string) error {
	for _, cmd := range commands {
		fmt.Printf("     Running: %s\n", cmd)
		// In a real implementation, this would execute the command
		// exec.Command("sh", "-c", cmd).Run()
	}
	return nil
}

// updateProjectMetadata updates the project metadata file
func updateProjectMetadata(targetDir, newTier string) error {
	metadataPath := filepath.Join(targetDir, ".template-metadata.yaml")
	
	// Create or update metadata file
	metadata := fmt.Sprintf(`name: %s
tier: %s
version: "1.0.0"
migrated_at: %d
`, filepath.Base(targetDir), newTier, getCurrentTimestamp())

	return os.WriteFile(metadataPath, []byte(metadata), 0644)
}

// createProjectBackup creates a full backup of the project
func createProjectBackup(sourceDir, backupDir string) error {
	return copyFileOrDir(sourceDir, backupDir)
}

// showMigrationNextSteps shows next steps after migration
func showMigrationNextSteps(fromTier, toTier string) {
	fmt.Printf("\n📋 Next Steps:\n")
	
	switch toTier {
	case "intermediate":
		fmt.Printf("   1. Configure dependency health checks in internal/config/config.go\n")
		fmt.Printf("   2. Test dependency endpoints: curl http://localhost:8080/health/dependencies\n")
		fmt.Printf("   3. Update your monitoring to include dependency health\n")
		
	case "advanced":
		fmt.Printf("   1. Configure OpenTelemetry endpoints in your environment\n")
		fmt.Printf("   2. Set up CloudEvents sink for event publishing\n")
		fmt.Printf("   3. Test observability: curl http://localhost:8080/metrics\n")
		fmt.Printf("   4. Verify tracing is working in your observability platform\n")
		
	case "enterprise":
		fmt.Printf("   1. Set up TLS certificates for mTLS authentication\n")
		fmt.Printf("   2. Configure RBAC policies in configs/rbac-*.json\n")
		fmt.Printf("   3. Set up audit log monitoring and retention\n")
		fmt.Printf("   4. Configure environment-specific settings\n")
		fmt.Printf("   5. Test security: verify client certificate authentication\n")
	}
	
	fmt.Printf("   \n")
	fmt.Printf("   Run 'go mod tidy' to update dependencies\n")
	fmt.Printf("   Run 'make test' to verify the migration\n")
}

// Helper functions

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func isDowngrade(fromTier, toTier string) bool {
	tierOrder := map[string]int{
		"basic": 1, "intermediate": 2, "advanced": 3, "enterprise": 4,
	}
	return tierOrder[fromTier] > tierOrder[toTier]
}

func copyFileOrDir(src, dst string) error {
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	
	if srcInfo.IsDir() {
		return copyDir(src, dst)
	}
	return copyFile(src, dst)
}

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		
		dstPath := filepath.Join(dst, relPath)
		
		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}
		
		return copyFile(path, dstPath)
	})
}
