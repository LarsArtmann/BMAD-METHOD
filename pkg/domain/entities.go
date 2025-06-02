package domain

import (
	"time"
	"errors"
	"fmt"
)

// Entity represents a domain entity with identity
type Entity interface {
	ID() string
	CreatedAt() time.Time
	UpdatedAt() time.Time
}

// AggregateRoot represents the root of an aggregate in DDD
type AggregateRoot interface {
	Entity
	DomainEvents() []DomainEvent
	ClearDomainEvents()
	AddDomainEvent(event DomainEvent)
}

// DomainEvent represents something that happened in the domain
type DomainEvent interface {
	OccurredAt() time.Time
	EventType() string
	AggregateID() string
	Version() int
}

// BaseEntity provides common entity functionality
type BaseEntity struct {
	id        string
	createdAt time.Time
	updatedAt time.Time
}

// NewBaseEntity creates a new base entity
func NewBaseEntity(id string) BaseEntity {
	now := time.Now()
	return BaseEntity{
		id:        id,
		createdAt: now,
		updatedAt: now,
	}
}

// ID returns the entity ID
func (e BaseEntity) ID() string {
	return e.id
}

// CreatedAt returns when the entity was created
func (e BaseEntity) CreatedAt() time.Time {
	return e.createdAt
}

// UpdatedAt returns when the entity was last updated
func (e BaseEntity) UpdatedAt() time.Time {
	return e.updatedAt
}

// Touch updates the entity's updated timestamp
func (e *BaseEntity) Touch() {
	e.updatedAt = time.Now()
}

// BaseAggregateRoot provides common aggregate root functionality
type BaseAggregateRoot struct {
	BaseEntity
	domainEvents []DomainEvent
	version      int
}

// NewBaseAggregateRoot creates a new base aggregate root
func NewBaseAggregateRoot(id string) BaseAggregateRoot {
	return BaseAggregateRoot{
		BaseEntity:   NewBaseEntity(id),
		domainEvents: make([]DomainEvent, 0),
		version:      1,
	}
}

// DomainEvents returns all domain events for this aggregate
func (ar BaseAggregateRoot) DomainEvents() []DomainEvent {
	return ar.domainEvents
}

// ClearDomainEvents clears all domain events
func (ar *BaseAggregateRoot) ClearDomainEvents() {
	ar.domainEvents = make([]DomainEvent, 0)
}

// AddDomainEvent adds a domain event to the aggregate
func (ar *BaseAggregateRoot) AddDomainEvent(event DomainEvent) {
	ar.domainEvents = append(ar.domainEvents, event)
}

// Version returns the aggregate version
func (ar BaseAggregateRoot) Version() int {
	return ar.version
}

// IncrementVersion increments the aggregate version
func (ar *BaseAggregateRoot) IncrementVersion() {
	ar.version++
	ar.Touch()
}

// Project represents a code generation project aggregate
type Project struct {
	BaseAggregateRoot
	name          ProjectName
	description   string
	tier          ProjectTier
	configuration ProjectConfiguration
	features      []ProjectFeature
	status        ProjectStatus
	owner         string
}

// ProjectName value object
type ProjectName struct {
	value string
}

// NewProjectName creates a new project name
func NewProjectName(name string) (ProjectName, error) {
	if name == "" {
		return ProjectName{}, errors.New("project name cannot be empty")
	}
	if len(name) > 100 {
		return ProjectName{}, errors.New("project name cannot exceed 100 characters")
	}
	return ProjectName{value: name}, nil
}

// String returns the project name as a string
func (pn ProjectName) String() string {
	return pn.value
}

// ProjectTier value object
type ProjectTier struct {
	value string
}

// NewProjectTier creates a new project tier
func NewProjectTier(tier string) (ProjectTier, error) {
	validTiers := []string{"basic", "intermediate", "advanced", "enterprise"}
	for _, validTier := range validTiers {
		if tier == validTier {
			return ProjectTier{value: tier}, nil
		}
	}
	return ProjectTier{}, fmt.Errorf("invalid project tier: %s", tier)
}

// String returns the tier as a string
func (pt ProjectTier) String() string {
	return pt.value
}

// IsAtLeast checks if this tier is at least the specified tier
func (pt ProjectTier) IsAtLeast(other ProjectTier) bool {
	tierOrder := map[string]int{
		"basic":        1,
		"intermediate": 2,
		"advanced":     3,
		"enterprise":   4,
	}
	return tierOrder[pt.value] >= tierOrder[other.value]
}

// ProjectConfiguration value object
type ProjectConfiguration struct {
	outputDir     string
	namespace     string
	packageName   string
	enableFeatures map[string]bool
	customSettings map[string]interface{}
}

// NewProjectConfiguration creates a new project configuration
func NewProjectConfiguration(outputDir, namespace, packageName string) ProjectConfiguration {
	return ProjectConfiguration{
		outputDir:      outputDir,
		namespace:      namespace,
		packageName:    packageName,
		enableFeatures: make(map[string]bool),
		customSettings: make(map[string]interface{}),
	}
}

// OutputDir returns the output directory
func (pc ProjectConfiguration) OutputDir() string {
	return pc.outputDir
}

// Namespace returns the namespace
func (pc ProjectConfiguration) Namespace() string {
	return pc.namespace
}

// PackageName returns the package name
func (pc ProjectConfiguration) PackageName() string {
	return pc.packageName
}

// IsFeatureEnabled checks if a feature is enabled
func (pc ProjectConfiguration) IsFeatureEnabled(feature string) bool {
	return pc.enableFeatures[feature]
}

// EnableFeature enables a feature
func (pc *ProjectConfiguration) EnableFeature(feature string) {
	pc.enableFeatures[feature] = true
}

// DisableFeature disables a feature
func (pc *ProjectConfiguration) DisableFeature(feature string) {
	pc.enableFeatures[feature] = false
}

// SetCustomSetting sets a custom setting
func (pc *ProjectConfiguration) SetCustomSetting(key string, value interface{}) {
	pc.customSettings[key] = value
}

// GetCustomSetting gets a custom setting
func (pc ProjectConfiguration) GetCustomSetting(key string) (interface{}, bool) {
	value, exists := pc.customSettings[key]
	return value, exists
}

// ProjectFeature value object
type ProjectFeature struct {
	id          string
	name        string
	enabled     bool
	config      map[string]interface{}
	dependencies []string
}

// NewProjectFeature creates a new project feature
func NewProjectFeature(id, name string) ProjectFeature {
	return ProjectFeature{
		id:           id,
		name:         name,
		enabled:      false,
		config:       make(map[string]interface{}),
		dependencies: make([]string, 0),
	}
}

// ID returns the feature ID
func (pf ProjectFeature) ID() string {
	return pf.id
}

// Name returns the feature name
func (pf ProjectFeature) Name() string {
	return pf.name
}

// IsEnabled returns whether the feature is enabled
func (pf ProjectFeature) IsEnabled() bool {
	return pf.enabled
}

// Enable enables the feature
func (pf *ProjectFeature) Enable() {
	pf.enabled = true
}

// Disable disables the feature
func (pf *ProjectFeature) Disable() {
	pf.enabled = false
}

// Dependencies returns the feature dependencies
func (pf ProjectFeature) Dependencies() []string {
	return pf.dependencies
}

// AddDependency adds a feature dependency
func (pf *ProjectFeature) AddDependency(dependency string) {
	pf.dependencies = append(pf.dependencies, dependency)
}

// ProjectStatus enumeration
type ProjectStatus int

const (
	ProjectStatusDraft ProjectStatus = iota
	ProjectStatusConfigured
	ProjectStatusGenerating
	ProjectStatusGenerated
	ProjectStatusFailed
	ProjectStatusArchived
)

// String returns the string representation of the project status
func (ps ProjectStatus) String() string {
	switch ps {
	case ProjectStatusDraft:
		return "draft"
	case ProjectStatusConfigured:
		return "configured"
	case ProjectStatusGenerating:
		return "generating"
	case ProjectStatusGenerated:
		return "generated"
	case ProjectStatusFailed:
		return "failed"
	case ProjectStatusArchived:
		return "archived"
	default:
		return "unknown"
	}
}

// NewProject creates a new project aggregate
func NewProject(id string, name ProjectName, tier ProjectTier, owner string) (*Project, error) {
	if id == "" {
		return nil, errors.New("project ID cannot be empty")
	}
	if owner == "" {
		return nil, errors.New("project owner cannot be empty")
	}

	project := &Project{
		BaseAggregateRoot: NewBaseAggregateRoot(id),
		name:              name,
		tier:              tier,
		configuration:     NewProjectConfiguration("", "", ""),
		features:          make([]ProjectFeature, 0),
		status:            ProjectStatusDraft,
		owner:             owner,
	}

	// Add domain event
	event := NewProjectCreatedEvent(id, name.String(), tier.String(), owner)
	project.AddDomainEvent(event)

	return project, nil
}

// Name returns the project name
func (p Project) Name() ProjectName {
	return p.name
}

// Tier returns the project tier
func (p Project) Tier() ProjectTier {
	return p.tier
}

// Configuration returns the project configuration
func (p Project) Configuration() ProjectConfiguration {
	return p.configuration
}

// Features returns the project features
func (p Project) Features() []ProjectFeature {
	return p.features
}

// Status returns the project status
func (p Project) Status() ProjectStatus {
	return p.status
}

// Owner returns the project owner
func (p Project) Owner() string {
	return p.owner
}

// UpdateDescription updates the project description
func (p *Project) UpdateDescription(description string) {
	p.description = description
	p.IncrementVersion()
	
	event := NewProjectDescriptionUpdatedEvent(p.ID(), description)
	p.AddDomainEvent(event)
}

// UpdateConfiguration updates the project configuration
func (p *Project) UpdateConfiguration(config ProjectConfiguration) {
	p.configuration = config
	p.IncrementVersion()
	
	event := NewProjectConfigurationUpdatedEvent(p.ID(), config)
	p.AddDomainEvent(event)
}

// AddFeature adds a feature to the project
func (p *Project) AddFeature(feature ProjectFeature) error {
	// Check if feature already exists
	for _, existingFeature := range p.features {
		if existingFeature.ID() == feature.ID() {
			return fmt.Errorf("feature %s already exists", feature.ID())
		}
	}

	p.features = append(p.features, feature)
	p.IncrementVersion()
	
	event := NewProjectFeatureAddedEvent(p.ID(), feature.ID(), feature.Name())
	p.AddDomainEvent(event)
	
	return nil
}

// RemoveFeature removes a feature from the project
func (p *Project) RemoveFeature(featureID string) error {
	for i, feature := range p.features {
		if feature.ID() == featureID {
			p.features = append(p.features[:i], p.features[i+1:]...)
			p.IncrementVersion()
			
			event := NewProjectFeatureRemovedEvent(p.ID(), featureID)
			p.AddDomainEvent(event)
			
			return nil
		}
	}
	return fmt.Errorf("feature %s not found", featureID)
}

// EnableFeature enables a feature in the project
func (p *Project) EnableFeature(featureID string) error {
	for i, feature := range p.features {
		if feature.ID() == featureID {
			p.features[i].Enable()
			p.IncrementVersion()
			
			event := NewProjectFeatureEnabledEvent(p.ID(), featureID)
			p.AddDomainEvent(event)
			
			return nil
		}
	}
	return fmt.Errorf("feature %s not found", featureID)
}

// DisableFeature disables a feature in the project
func (p *Project) DisableFeature(featureID string) error {
	for i, feature := range p.features {
		if feature.ID() == featureID {
			p.features[i].Disable()
			p.IncrementVersion()
			
			event := NewProjectFeatureDisabledEvent(p.ID(), featureID)
			p.AddDomainEvent(event)
			
			return nil
		}
	}
	return fmt.Errorf("feature %s not found", featureID)
}

// ChangeStatus changes the project status
func (p *Project) ChangeStatus(newStatus ProjectStatus) error {
	if !p.isValidStatusTransition(p.status, newStatus) {
		return fmt.Errorf("invalid status transition from %s to %s", p.status, newStatus)
	}

	oldStatus := p.status
	p.status = newStatus
	p.IncrementVersion()
	
	event := NewProjectStatusChangedEvent(p.ID(), oldStatus.String(), newStatus.String())
	p.AddDomainEvent(event)
	
	return nil
}

// isValidStatusTransition validates status transitions
func (p *Project) isValidStatusTransition(from, to ProjectStatus) bool {
	validTransitions := map[ProjectStatus][]ProjectStatus{
		ProjectStatusDraft:       {ProjectStatusConfigured, ProjectStatusArchived},
		ProjectStatusConfigured:  {ProjectStatusGenerating, ProjectStatusDraft, ProjectStatusArchived},
		ProjectStatusGenerating:  {ProjectStatusGenerated, ProjectStatusFailed},
		ProjectStatusGenerated:   {ProjectStatusConfigured, ProjectStatusArchived},
		ProjectStatusFailed:      {ProjectStatusConfigured, ProjectStatusArchived},
		ProjectStatusArchived:    {}, // No transitions from archived
	}

	allowedTransitions, exists := validTransitions[from]
	if !exists {
		return false
	}

	for _, allowedTo := range allowedTransitions {
		if allowedTo == to {
			return true
		}
	}

	return false
}

// CanBeGenerated checks if the project can be generated
func (p *Project) CanBeGenerated() bool {
	return p.status == ProjectStatusConfigured && len(p.getEnabledFeatures()) > 0
}

// getEnabledFeatures returns only enabled features
func (p *Project) getEnabledFeatures() []ProjectFeature {
	var enabled []ProjectFeature
	for _, feature := range p.features {
		if feature.IsEnabled() {
			enabled = append(enabled, feature)
		}
	}
	return enabled
}

// Generation represents a project generation aggregate
type Generation struct {
	BaseAggregateRoot
	projectID    string
	status       GenerationStatus
	requestedBy  string
	startedAt    time.Time
	completedAt  *time.Time
	artifacts    []GenerationArtifact
	errors       []GenerationError
}

// GenerationStatus enumeration
type GenerationStatus int

const (
	GenerationStatusPending GenerationStatus = iota
	GenerationStatusRunning
	GenerationStatusCompleted
	GenerationStatusFailed
	GenerationStatusCancelled
)

// String returns the string representation of the generation status
func (gs GenerationStatus) String() string {
	switch gs {
	case GenerationStatusPending:
		return "pending"
	case GenerationStatusRunning:
		return "running"
	case GenerationStatusCompleted:
		return "completed"
	case GenerationStatusFailed:
		return "failed"
	case GenerationStatusCancelled:
		return "cancelled"
	default:
		return "unknown"
	}
}

// GenerationArtifact represents a generated artifact
type GenerationArtifact struct {
	path     string
	size     int64
	checksum string
	mimeType string
}

// NewGenerationArtifact creates a new generation artifact
func NewGenerationArtifact(path string, size int64, checksum, mimeType string) GenerationArtifact {
	return GenerationArtifact{
		path:     path,
		size:     size,
		checksum: checksum,
		mimeType: mimeType,
	}
}

// Path returns the artifact path
func (ga GenerationArtifact) Path() string {
	return ga.path
}

// Size returns the artifact size
func (ga GenerationArtifact) Size() int64 {
	return ga.size
}

// Checksum returns the artifact checksum
func (ga GenerationArtifact) Checksum() string {
	return ga.checksum
}

// MimeType returns the artifact MIME type
func (ga GenerationArtifact) MimeType() string {
	return ga.mimeType
}

// GenerationError represents an error during generation
type GenerationError struct {
	code      string
	message   string
	details   map[string]interface{}
	occurredAt time.Time
}

// NewGenerationError creates a new generation error
func NewGenerationError(code, message string, details map[string]interface{}) GenerationError {
	return GenerationError{
		code:       code,
		message:    message,
		details:    details,
		occurredAt: time.Now(),
	}
}

// Code returns the error code
func (ge GenerationError) Code() string {
	return ge.code
}

// Message returns the error message
func (ge GenerationError) Message() string {
	return ge.message
}

// Details returns the error details
func (ge GenerationError) Details() map[string]interface{} {
	return ge.details
}

// OccurredAt returns when the error occurred
func (ge GenerationError) OccurredAt() time.Time {
	return ge.occurredAt
}

// NewGeneration creates a new generation aggregate
func NewGeneration(id, projectID, requestedBy string) (*Generation, error) {
	if id == "" {
		return nil, errors.New("generation ID cannot be empty")
	}
	if projectID == "" {
		return nil, errors.New("project ID cannot be empty")
	}
	if requestedBy == "" {
		return nil, errors.New("requestedBy cannot be empty")
	}

	generation := &Generation{
		BaseAggregateRoot: NewBaseAggregateRoot(id),
		projectID:         projectID,
		status:            GenerationStatusPending,
		requestedBy:       requestedBy,
		startedAt:         time.Now(),
		artifacts:         make([]GenerationArtifact, 0),
		errors:            make([]GenerationError, 0),
	}

	event := NewGenerationStartedEvent(id, projectID, requestedBy)
	generation.AddDomainEvent(event)

	return generation, nil
}

// ProjectID returns the project ID
func (g Generation) ProjectID() string {
	return g.projectID
}

// Status returns the generation status
func (g Generation) Status() GenerationStatus {
	return g.status
}

// RequestedBy returns who requested the generation
func (g Generation) RequestedBy() string {
	return g.requestedBy
}

// StartedAt returns when the generation started
func (g Generation) StartedAt() time.Time {
	return g.startedAt
}

// CompletedAt returns when the generation completed
func (g Generation) CompletedAt() *time.Time {
	return g.completedAt
}

// Artifacts returns the generation artifacts
func (g Generation) Artifacts() []GenerationArtifact {
	return g.artifacts
}

// Errors returns the generation errors
func (g Generation) Errors() []GenerationError {
	return g.errors
}

// Start starts the generation
func (g *Generation) Start() error {
	if g.status != GenerationStatusPending {
		return fmt.Errorf("cannot start generation with status %s", g.status)
	}

	g.status = GenerationStatusRunning
	g.IncrementVersion()

	event := NewGenerationProgressEvent(g.ID(), "Generation started")
	g.AddDomainEvent(event)

	return nil
}

// Complete completes the generation successfully
func (g *Generation) Complete() error {
	if g.status != GenerationStatusRunning {
		return fmt.Errorf("cannot complete generation with status %s", g.status)
	}

	g.status = GenerationStatusCompleted
	now := time.Now()
	g.completedAt = &now
	g.IncrementVersion()

	event := NewGenerationCompletedEvent(g.ID(), len(g.artifacts))
	g.AddDomainEvent(event)

	return nil
}

// Fail marks the generation as failed
func (g *Generation) Fail(err GenerationError) error {
	if g.status != GenerationStatusRunning {
		return fmt.Errorf("cannot fail generation with status %s", g.status)
	}

	g.status = GenerationStatusFailed
	g.errors = append(g.errors, err)
	now := time.Now()
	g.completedAt = &now
	g.IncrementVersion()

	event := NewGenerationFailedEvent(g.ID(), err.Code(), err.Message())
	g.AddDomainEvent(event)

	return nil
}

// Cancel cancels the generation
func (g *Generation) Cancel() error {
	if g.status != GenerationStatusPending && g.status != GenerationStatusRunning {
		return fmt.Errorf("cannot cancel generation with status %s", g.status)
	}

	g.status = GenerationStatusCancelled
	now := time.Now()
	g.completedAt = &now
	g.IncrementVersion()

	event := NewGenerationCancelledEvent(g.ID())
	g.AddDomainEvent(event)

	return nil
}

// AddArtifact adds an artifact to the generation
func (g *Generation) AddArtifact(artifact GenerationArtifact) {
	g.artifacts = append(g.artifacts, artifact)
	g.IncrementVersion()

	event := NewGenerationArtifactCreatedEvent(g.ID(), artifact.Path(), artifact.Size())
	g.AddDomainEvent(event)
}