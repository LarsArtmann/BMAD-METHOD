package features

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/LarsArtmann/BMAD-METHOD/pkg/config"
)

// FeatureType represents different types of features that can be composed
type FeatureType string

const (
	FeatureTypeCore         FeatureType = "core"
	FeatureTypeObservability FeatureType = "observability"
	FeatureTypeSecurity     FeatureType = "security"
	FeatureTypeStorage      FeatureType = "storage"
	FeatureTypeAPI          FeatureType = "api"
	FeatureTypeDeployment   FeatureType = "deployment"
	FeatureTypeMessaging    FeatureType = "messaging"
	FeatureTypeCaching      FeatureType = "caching"
)

// Feature represents a composable feature with metadata and dependencies
type Feature struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Type        FeatureType            `json:"type"`
	Version     string                 `json:"version"`
	
	// Dependencies and conflicts
	Dependencies []string               `json:"dependencies"`
	Conflicts    []string               `json:"conflicts"`
	Requires     []string               `json:"requires"`
	
	// Configuration
	ConfigSchema map[string]interface{} `json:"config_schema"`
	DefaultConfig map[string]interface{} `json:"default_config"`
	
	// Composition metadata
	Priority     int                    `json:"priority"`
	Category     string                 `json:"category"`
	Tags         []string               `json:"tags"`
	
	// Tier compatibility
	MinTier      string                 `json:"min_tier"`
	MaxTier      string                 `json:"max_tier"`
	
	// Implementation
	Generator    FeatureGenerator       `json:"-"`
	Validator    FeatureValidator       `json:"-"`
}

// FeatureGenerator defines the interface for generating feature components
type FeatureGenerator interface {
	Generate(ctx context.Context, config *config.ProjectConfig, featureConfig map[string]interface{}) (*GeneratedFeature, error)
	Validate(config *config.ProjectConfig, featureConfig map[string]interface{}) error
	GetTemplates() []string
	GetAssets() []string
}

// FeatureValidator validates feature configuration and compatibility
type FeatureValidator interface {
	ValidateConfig(featureConfig map[string]interface{}) error
	ValidateCompatibility(otherFeatures []Feature) error
	ValidateTier(tier string) error
}

// GeneratedFeature represents the output of a feature generator
type GeneratedFeature struct {
	Files       map[string]string      `json:"files"`
	Templates   []string               `json:"templates"`
	Assets      []string               `json:"assets"`
	Metadata    map[string]interface{} `json:"metadata"`
	PostActions []PostGenerationAction `json:"post_actions"`
}

// PostGenerationAction represents actions to perform after feature generation
type PostGenerationAction struct {
	Type        string                 `json:"type"`
	Description string                 `json:"description"`
	Command     string                 `json:"command"`
	Args        []string               `json:"args"`
	Env         map[string]string      `json:"env"`
	WorkingDir  string                 `json:"working_dir"`
}

// FeatureComposer orchestrates feature composition and conflict resolution
type FeatureComposer struct {
	registry     *FeatureRegistry
	resolver     *DependencyResolver
	validator    *CompositionValidator
	generator    *CompositionGenerator
	
	mu           sync.RWMutex
}

// FeatureRegistry manages available features and their metadata
type FeatureRegistry struct {
	features     map[string]*Feature
	byType       map[FeatureType][]*Feature
	byCategory   map[string][]*Feature
	
	mu           sync.RWMutex
}

// DependencyResolver handles feature dependency resolution
type DependencyResolver struct {
	registry     *FeatureRegistry
}

// CompositionValidator validates feature compositions for conflicts and compatibility
type CompositionValidator struct {
	registry     *FeatureRegistry
}

// CompositionGenerator generates the final composed output
type CompositionGenerator struct {
	registry     *FeatureRegistry
}

// CompositionRequest represents a request to compose features
type CompositionRequest struct {
	Features     []string               `json:"features"`
	Config       *config.ProjectConfig  `json:"config"`
	Options      CompositionOptions     `json:"options"`
	FeatureConfigs map[string]map[string]interface{} `json:"feature_configs"`
}

// CompositionOptions controls how features are composed
type CompositionOptions struct {
	AutoResolveDependencies bool              `json:"auto_resolve_dependencies"`
	FailOnConflicts        bool              `json:"fail_on_conflicts"`
	PreferredVersions      map[string]string `json:"preferred_versions"`
	ExcludeFeatures        []string          `json:"exclude_features"`
	IncludeOptional        bool              `json:"include_optional"`
	DryRun                 bool              `json:"dry_run"`
}

// CompositionResult represents the result of feature composition
type CompositionResult struct {
	ResolvedFeatures  []string                    `json:"resolved_features"`
	GeneratedOutput   *GeneratedFeature           `json:"generated_output"`
	Dependencies      map[string][]string         `json:"dependencies"`
	Conflicts         []ConflictInfo              `json:"conflicts"`
	Warnings          []string                    `json:"warnings"`
	PostActions       []PostGenerationAction      `json:"post_actions"`
	Metadata          map[string]interface{}      `json:"metadata"`
}

// ConflictInfo represents information about feature conflicts
type ConflictInfo struct {
	Feature1    string   `json:"feature1"`
	Feature2    string   `json:"feature2"`
	ConflictType string  `json:"conflict_type"`
	Description string   `json:"description"`
	Resolution  string   `json:"resolution"`
}

// NewFeatureComposer creates a new feature composer
func NewFeatureComposer() *FeatureComposer {
	registry := &FeatureRegistry{
		features:   make(map[string]*Feature),
		byType:     make(map[FeatureType][]*Feature),
		byCategory: make(map[string][]*Feature),
	}
	
	return &FeatureComposer{
		registry:  registry,
		resolver:  &DependencyResolver{registry: registry},
		validator: &CompositionValidator{registry: registry},
		generator: &CompositionGenerator{registry: registry},
	}
}

// RegisterFeature registers a new feature in the registry
func (fc *FeatureComposer) RegisterFeature(feature *Feature) error {
	fc.mu.Lock()
	defer fc.mu.Unlock()
	
	// Validate feature
	if feature.ID == "" {
		return fmt.Errorf("feature ID cannot be empty")
	}
	
	if feature.Generator == nil {
		return fmt.Errorf("feature %s must have a generator", feature.ID)
	}
	
	// Check for conflicts with existing features
	if existing, exists := fc.registry.features[feature.ID]; exists {
		return fmt.Errorf("feature %s already registered with version %s", feature.ID, existing.Version)
	}
	
	// Register feature
	fc.registry.mu.Lock()
	defer fc.registry.mu.Unlock()
	
	fc.registry.features[feature.ID] = feature
	
	// Index by type
	fc.registry.byType[feature.Type] = append(fc.registry.byType[feature.Type], feature)
	
	// Index by category
	if feature.Category != "" {
		fc.registry.byCategory[feature.Category] = append(fc.registry.byCategory[feature.Category], feature)
	}
	
	return nil
}

// GetFeature retrieves a feature by ID
func (fc *FeatureComposer) GetFeature(id string) (*Feature, error) {
	fc.registry.mu.RLock()
	defer fc.registry.mu.RUnlock()
	
	feature, exists := fc.registry.features[id]
	if !exists {
		return nil, fmt.Errorf("feature %s not found", id)
	}
	
	return feature, nil
}

// ListFeatures returns all registered features, optionally filtered by type
func (fc *FeatureComposer) ListFeatures(featureType *FeatureType) []*Feature {
	fc.registry.mu.RLock()
	defer fc.registry.mu.RUnlock()
	
	if featureType != nil {
		features := make([]*Feature, len(fc.registry.byType[*featureType]))
		copy(features, fc.registry.byType[*featureType])
		return features
	}
	
	features := make([]*Feature, 0, len(fc.registry.features))
	for _, feature := range fc.registry.features {
		features = append(features, feature)
	}
	
	// Sort by priority (higher priority first)
	sort.Slice(features, func(i, j int) bool {
		return features[i].Priority > features[j].Priority
	})
	
	return features
}

// ComposeFeatures composes the requested features into a final output
func (fc *FeatureComposer) ComposeFeatures(ctx context.Context, request *CompositionRequest) (*CompositionResult, error) {
	result := &CompositionResult{
		Dependencies: make(map[string][]string),
		Conflicts:    make([]ConflictInfo, 0),
		Warnings:     make([]string, 0),
		PostActions:  make([]PostGenerationAction, 0),
		Metadata:     make(map[string]interface{}),
	}
	
	// Step 1: Resolve dependencies
	resolvedFeatures, err := fc.resolver.ResolveDependencies(request.Features, request.Options)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve dependencies: %w", err)
	}
	result.ResolvedFeatures = resolvedFeatures
	
	// Step 2: Validate composition
	conflicts, warnings, err := fc.validator.ValidateComposition(resolvedFeatures, request.Config)
	if err != nil {
		return nil, fmt.Errorf("composition validation failed: %w", err)
	}
	result.Conflicts = conflicts
	result.Warnings = warnings
	
	// Step 3: Handle conflicts
	if len(conflicts) > 0 && request.Options.FailOnConflicts {
		return result, fmt.Errorf("composition has %d conflicts", len(conflicts))
	}
	
	// Step 4: Generate composed output
	if !request.Options.DryRun {
		output, postActions, err := fc.generator.GenerateComposition(ctx, resolvedFeatures, request.Config, request.FeatureConfigs)
		if err != nil {
			return nil, fmt.Errorf("failed to generate composition: %w", err)
		}
		result.GeneratedOutput = output
		result.PostActions = postActions
	}
	
	// Step 5: Build dependency graph
	for _, featureID := range resolvedFeatures {
		if feature, err := fc.GetFeature(featureID); err == nil {
			result.Dependencies[featureID] = feature.Dependencies
		}
	}
	
	return result, nil
}

// ResolveDependencies resolves feature dependencies recursively
func (dr *DependencyResolver) ResolveDependencies(features []string, options CompositionOptions) ([]string, error) {
	resolved := make(map[string]bool)
	var resolveOrder []string
	
	var resolve func(string) error
	resolve = func(featureID string) error {
		if resolved[featureID] {
			return nil
		}
		
		// Check if feature is excluded
		for _, excluded := range options.ExcludeFeatures {
			if excluded == featureID {
				return nil
			}
		}
		
		feature, exists := dr.registry.features[featureID]
		if !exists {
			return fmt.Errorf("feature %s not found", featureID)
		}
		
		// Resolve dependencies first
		for _, dep := range feature.Dependencies {
			if err := resolve(dep); err != nil {
				return fmt.Errorf("failed to resolve dependency %s for feature %s: %w", dep, featureID, err)
			}
		}
		
		// Add to resolved list
		resolved[featureID] = true
		resolveOrder = append(resolveOrder, featureID)
		
		return nil
	}
	
	// Resolve all requested features
	for _, featureID := range features {
		if err := resolve(featureID); err != nil {
			if options.AutoResolveDependencies {
				continue // Skip unresolvable features
			}
			return nil, err
		}
	}
	
	return resolveOrder, nil
}

// ValidateComposition validates a feature composition for conflicts and compatibility
func (cv *CompositionValidator) ValidateComposition(features []string, config *config.ProjectConfig) ([]ConflictInfo, []string, error) {
	var conflicts []ConflictInfo
	var warnings []string
	
	featureMap := make(map[string]*Feature)
	for _, featureID := range features {
		feature, exists := cv.registry.features[featureID]
		if !exists {
			return nil, nil, fmt.Errorf("feature %s not found", featureID)
		}
		featureMap[featureID] = feature
	}
	
	// Check for conflicts between features
	for i, feature1ID := range features {
		feature1 := featureMap[feature1ID]
		
		// Validate tier compatibility
		if !cv.isCompatibleWithTier(feature1, config.Tier) {
			warnings = append(warnings, fmt.Sprintf("Feature %s may not be compatible with tier %s", feature1ID, config.Tier))
		}
		
		for j := i + 1; j < len(features); j++ {
			feature2ID := features[j]
			feature2 := featureMap[feature2ID]
			
			// Check explicit conflicts
			for _, conflict := range feature1.Conflicts {
				if conflict == feature2ID {
					conflicts = append(conflicts, ConflictInfo{
						Feature1:     feature1ID,
						Feature2:     feature2ID,
						ConflictType: "explicit",
						Description:  fmt.Sprintf("Feature %s explicitly conflicts with %s", feature1ID, feature2ID),
						Resolution:   "Remove one of the conflicting features",
					})
				}
			}
			
			// Check type conflicts (e.g., multiple storage backends)
			if cv.areTypeConflicts(feature1, feature2) {
				conflicts = append(conflicts, ConflictInfo{
					Feature1:     feature1ID,
					Feature2:     feature2ID,
					ConflictType: "type",
					Description:  fmt.Sprintf("Features %s and %s are of conflicting types", feature1ID, feature2ID),
					Resolution:   "Choose only one feature of this type",
				})
			}
		}
	}
	
	return conflicts, warnings, nil
}

// isCompatibleWithTier checks if a feature is compatible with the given tier
func (cv *CompositionValidator) isCompatibleWithTier(feature *Feature, tier string) bool {
	tierOrder := map[string]int{
		"basic":        1,
		"intermediate": 2,
		"advanced":     3,
		"enterprise":   4,
	}
	
	currentTierLevel := tierOrder[tier]
	
	if feature.MinTier != "" {
		minLevel := tierOrder[feature.MinTier]
		if currentTierLevel < minLevel {
			return false
		}
	}
	
	if feature.MaxTier != "" {
		maxLevel := tierOrder[feature.MaxTier]
		if currentTierLevel > maxLevel {
			return false
		}
	}
	
	return true
}

// areTypeConflicts checks if two features have conflicting types
func (cv *CompositionValidator) areTypeConflicts(feature1, feature2 *Feature) bool {
	// Define conflicting feature combinations
	conflictingTypes := map[FeatureType][]FeatureType{
		FeatureTypeStorage: {FeatureTypeStorage}, // Only one storage backend allowed
		FeatureTypeCaching: {FeatureTypeCaching}, // Only one caching solution allowed
	}
	
	if conflicts, exists := conflictingTypes[feature1.Type]; exists {
		for _, conflictType := range conflicts {
			if feature2.Type == conflictType && feature1.ID != feature2.ID {
				return true
			}
		}
	}
	
	return false
}

// GenerateComposition generates the final composed output from resolved features
func (cg *CompositionGenerator) GenerateComposition(ctx context.Context, features []string, config *config.ProjectConfig, featureConfigs map[string]map[string]interface{}) (*GeneratedFeature, []PostGenerationAction, error) {
	result := &GeneratedFeature{
		Files:     make(map[string]string),
		Templates: make([]string, 0),
		Assets:    make([]string, 0),
		Metadata:  make(map[string]interface{}),
	}
	
	var allPostActions []PostGenerationAction
	
	// Generate each feature in dependency order
	for _, featureID := range features {
		feature, exists := cg.registry.features[featureID]
		if !exists {
			return nil, nil, fmt.Errorf("feature %s not found", featureID)
		}
		
		// Get feature-specific configuration
		featureConfig := featureConfigs[featureID]
		if featureConfig == nil {
			featureConfig = feature.DefaultConfig
		}
		
		// Generate feature output
		generated, err := feature.Generator.Generate(ctx, config, featureConfig)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to generate feature %s: %w", featureID, err)
		}
		
		// Merge files (later features can override earlier ones)
		for path, content := range generated.Files {
			result.Files[path] = content
		}
		
		// Accumulate templates and assets
		result.Templates = append(result.Templates, generated.Templates...)
		result.Assets = append(result.Assets, generated.Assets...)
		
		// Merge metadata
		for key, value := range generated.Metadata {
			result.Metadata[key] = value
		}
		
		// Collect post-generation actions
		allPostActions = append(allPostActions, generated.PostActions...)
	}
	
	// Remove duplicates from templates and assets
	result.Templates = removeDuplicates(result.Templates)
	result.Assets = removeDuplicates(result.Assets)
	
	return result, allPostActions, nil
}

// GetFeatureRecommendations suggests features based on project configuration
func (fc *FeatureComposer) GetFeatureRecommendations(config *config.ProjectConfig) ([]string, error) {
	var recommendations []string
	
	// Core recommendations based on tier
	switch config.Tier {
	case "basic":
		recommendations = append(recommendations, "health-basic", "logging-basic")
	case "intermediate":
		recommendations = append(recommendations, "health-intermediate", "logging-structured", "metrics-basic")
	case "advanced":
		recommendations = append(recommendations, "health-advanced", "observability-full", "security-rbac")
	case "enterprise":
		recommendations = append(recommendations, "health-enterprise", "observability-enterprise", "security-enterprise", "compliance-audit")
	}
	
	// Feature-specific recommendations
	if config.Features.Database {
		recommendations = append(recommendations, "storage-database", "migrations")
	}
	
	if config.Features.Cache {
		recommendations = append(recommendations, "caching-redis")
	}
	
	if config.Features.MessageQueue {
		recommendations = append(recommendations, "messaging-queue")
	}
	
	if config.Features.TypeScript {
		recommendations = append(recommendations, "typescript-client", "api-sdk")
	}
	
	return recommendations, nil
}

// removeDuplicates removes duplicate strings from a slice
func removeDuplicates(items []string) []string {
	seen := make(map[string]bool)
	var result []string
	
	for _, item := range items {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}
	
	return result
}

// GetCompositionSummary returns a human-readable summary of the composition
func (result *CompositionResult) GetCompositionSummary() string {
	var summary strings.Builder
	
	summary.WriteString(fmt.Sprintf("Feature Composition Summary:\n"))
	summary.WriteString(fmt.Sprintf("- Resolved Features: %d\n", len(result.ResolvedFeatures)))
	summary.WriteString(fmt.Sprintf("- Generated Files: %d\n", len(result.GeneratedOutput.Files)))
	summary.WriteString(fmt.Sprintf("- Post Actions: %d\n", len(result.PostActions)))
	summary.WriteString(fmt.Sprintf("- Conflicts: %d\n", len(result.Conflicts)))
	summary.WriteString(fmt.Sprintf("- Warnings: %d\n", len(result.Warnings)))
	
	if len(result.ResolvedFeatures) > 0 {
		summary.WriteString("\nIncluded Features:\n")
		for _, feature := range result.ResolvedFeatures {
			summary.WriteString(fmt.Sprintf("  - %s\n", feature))
		}
	}
	
	if len(result.Conflicts) > 0 {
		summary.WriteString("\nConflicts:\n")
		for _, conflict := range result.Conflicts {
			summary.WriteString(fmt.Sprintf("  - %s vs %s: %s\n", conflict.Feature1, conflict.Feature2, conflict.Description))
		}
	}
	
	if len(result.Warnings) > 0 {
		summary.WriteString("\nWarnings:\n")
		for _, warning := range result.Warnings {
			summary.WriteString(fmt.Sprintf("  - %s\n", warning))
		}
	}
	
	return summary.String()
}