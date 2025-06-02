package domain

import (
	"time"
)

// BaseDomainEvent provides common functionality for domain events
type BaseDomainEvent struct {
	occurredAt   time.Time
	eventType    string
	aggregateID  string
	version      int
}

// NewBaseDomainEvent creates a new base domain event
func NewBaseDomainEvent(eventType, aggregateID string, version int) BaseDomainEvent {
	return BaseDomainEvent{
		occurredAt:  time.Now(),
		eventType:   eventType,
		aggregateID: aggregateID,
		version:     version,
	}
}

// OccurredAt returns when the event occurred
func (e BaseDomainEvent) OccurredAt() time.Time {
	return e.occurredAt
}

// EventType returns the event type
func (e BaseDomainEvent) EventType() string {
	return e.eventType
}

// AggregateID returns the aggregate ID
func (e BaseDomainEvent) AggregateID() string {
	return e.aggregateID
}

// Version returns the event version
func (e BaseDomainEvent) Version() int {
	return e.version
}

// Project Events

// ProjectCreatedEvent is fired when a project is created
type ProjectCreatedEvent struct {
	BaseDomainEvent
	Name  string
	Tier  string
	Owner string
}

// NewProjectCreatedEvent creates a new project created event
func NewProjectCreatedEvent(aggregateID, name, tier, owner string) *ProjectCreatedEvent {
	return &ProjectCreatedEvent{
		BaseDomainEvent: NewBaseDomainEvent("ProjectCreated", aggregateID, 1),
		Name:            name,
		Tier:            tier,
		Owner:           owner,
	}
}

// ProjectDescriptionUpdatedEvent is fired when a project description is updated
type ProjectDescriptionUpdatedEvent struct {
	BaseDomainEvent
	Description string
}

// NewProjectDescriptionUpdatedEvent creates a new project description updated event
func NewProjectDescriptionUpdatedEvent(aggregateID, description string) *ProjectDescriptionUpdatedEvent {
	return &ProjectDescriptionUpdatedEvent{
		BaseDomainEvent: NewBaseDomainEvent("ProjectDescriptionUpdated", aggregateID, 1),
		Description:     description,
	}
}

// ProjectConfigurationUpdatedEvent is fired when a project configuration is updated
type ProjectConfigurationUpdatedEvent struct {
	BaseDomainEvent
	Configuration ProjectConfiguration
}

// NewProjectConfigurationUpdatedEvent creates a new project configuration updated event
func NewProjectConfigurationUpdatedEvent(aggregateID string, config ProjectConfiguration) *ProjectConfigurationUpdatedEvent {
	return &ProjectConfigurationUpdatedEvent{
		BaseDomainEvent: NewBaseDomainEvent("ProjectConfigurationUpdated", aggregateID, 1),
		Configuration:   config,
	}
}

// ProjectFeatureAddedEvent is fired when a feature is added to a project
type ProjectFeatureAddedEvent struct {
	BaseDomainEvent
	FeatureID   string
	FeatureName string
}

// NewProjectFeatureAddedEvent creates a new project feature added event
func NewProjectFeatureAddedEvent(aggregateID, featureID, featureName string) *ProjectFeatureAddedEvent {
	return &ProjectFeatureAddedEvent{
		BaseDomainEvent: NewBaseDomainEvent("ProjectFeatureAdded", aggregateID, 1),
		FeatureID:       featureID,
		FeatureName:     featureName,
	}
}

// ProjectFeatureRemovedEvent is fired when a feature is removed from a project
type ProjectFeatureRemovedEvent struct {
	BaseDomainEvent
	FeatureID string
}

// NewProjectFeatureRemovedEvent creates a new project feature removed event
func NewProjectFeatureRemovedEvent(aggregateID, featureID string) *ProjectFeatureRemovedEvent {
	return &ProjectFeatureRemovedEvent{
		BaseDomainEvent: NewBaseDomainEvent("ProjectFeatureRemoved", aggregateID, 1),
		FeatureID:       featureID,
	}
}

// ProjectFeatureEnabledEvent is fired when a feature is enabled in a project
type ProjectFeatureEnabledEvent struct {
	BaseDomainEvent
	FeatureID string
}

// NewProjectFeatureEnabledEvent creates a new project feature enabled event
func NewProjectFeatureEnabledEvent(aggregateID, featureID string) *ProjectFeatureEnabledEvent {
	return &ProjectFeatureEnabledEvent{
		BaseDomainEvent: NewBaseDomainEvent("ProjectFeatureEnabled", aggregateID, 1),
		FeatureID:       featureID,
	}
}

// ProjectFeatureDisabledEvent is fired when a feature is disabled in a project
type ProjectFeatureDisabledEvent struct {
	BaseDomainEvent
	FeatureID string
}

// NewProjectFeatureDisabledEvent creates a new project feature disabled event
func NewProjectFeatureDisabledEvent(aggregateID, featureID string) *ProjectFeatureDisabledEvent {
	return &ProjectFeatureDisabledEvent{
		BaseDomainEvent: NewBaseDomainEvent("ProjectFeatureDisabled", aggregateID, 1),
		FeatureID:       featureID,
	}
}

// ProjectStatusChangedEvent is fired when a project status changes
type ProjectStatusChangedEvent struct {
	BaseDomainEvent
	FromStatus string
	ToStatus   string
}

// NewProjectStatusChangedEvent creates a new project status changed event
func NewProjectStatusChangedEvent(aggregateID, fromStatus, toStatus string) *ProjectStatusChangedEvent {
	return &ProjectStatusChangedEvent{
		BaseDomainEvent: NewBaseDomainEvent("ProjectStatusChanged", aggregateID, 1),
		FromStatus:      fromStatus,
		ToStatus:        toStatus,
	}
}

// Generation Events

// GenerationStartedEvent is fired when a generation is started
type GenerationStartedEvent struct {
	BaseDomainEvent
	ProjectID   string
	RequestedBy string
}

// NewGenerationStartedEvent creates a new generation started event
func NewGenerationStartedEvent(aggregateID, projectID, requestedBy string) *GenerationStartedEvent {
	return &GenerationStartedEvent{
		BaseDomainEvent: NewBaseDomainEvent("GenerationStarted", aggregateID, 1),
		ProjectID:       projectID,
		RequestedBy:     requestedBy,
	}
}

// GenerationProgressEvent is fired to report generation progress
type GenerationProgressEvent struct {
	BaseDomainEvent
	Message string
}

// NewGenerationProgressEvent creates a new generation progress event
func NewGenerationProgressEvent(aggregateID, message string) *GenerationProgressEvent {
	return &GenerationProgressEvent{
		BaseDomainEvent: NewBaseDomainEvent("GenerationProgress", aggregateID, 1),
		Message:         message,
	}
}

// GenerationCompletedEvent is fired when a generation completes successfully
type GenerationCompletedEvent struct {
	BaseDomainEvent
	ArtifactCount int
}

// NewGenerationCompletedEvent creates a new generation completed event
func NewGenerationCompletedEvent(aggregateID string, artifactCount int) *GenerationCompletedEvent {
	return &GenerationCompletedEvent{
		BaseDomainEvent: NewBaseDomainEvent("GenerationCompleted", aggregateID, 1),
		ArtifactCount:   artifactCount,
	}
}

// GenerationFailedEvent is fired when a generation fails
type GenerationFailedEvent struct {
	BaseDomainEvent
	ErrorCode    string
	ErrorMessage string
}

// NewGenerationFailedEvent creates a new generation failed event
func NewGenerationFailedEvent(aggregateID, errorCode, errorMessage string) *GenerationFailedEvent {
	return &GenerationFailedEvent{
		BaseDomainEvent: NewBaseDomainEvent("GenerationFailed", aggregateID, 1),
		ErrorCode:       errorCode,
		ErrorMessage:    errorMessage,
	}
}

// GenerationCancelledEvent is fired when a generation is cancelled
type GenerationCancelledEvent struct {
	BaseDomainEvent
}

// NewGenerationCancelledEvent creates a new generation cancelled event
func NewGenerationCancelledEvent(aggregateID string) *GenerationCancelledEvent {
	return &GenerationCancelledEvent{
		BaseDomainEvent: NewBaseDomainEvent("GenerationCancelled", aggregateID, 1),
	}
}

// GenerationArtifactCreatedEvent is fired when an artifact is created during generation
type GenerationArtifactCreatedEvent struct {
	BaseDomainEvent
	ArtifactPath string
	ArtifactSize int64
}

// NewGenerationArtifactCreatedEvent creates a new generation artifact created event
func NewGenerationArtifactCreatedEvent(aggregateID, artifactPath string, artifactSize int64) *GenerationArtifactCreatedEvent {
	return &GenerationArtifactCreatedEvent{
		BaseDomainEvent: NewBaseDomainEvent("GenerationArtifactCreated", aggregateID, 1),
		ArtifactPath:    artifactPath,
		ArtifactSize:    artifactSize,
	}
}

// Feature Events

// FeatureRegisteredEvent is fired when a new feature is registered
type FeatureRegisteredEvent struct {
	BaseDomainEvent
	FeatureID   string
	FeatureName string
	FeatureType string
	Version     string
}

// NewFeatureRegisteredEvent creates a new feature registered event
func NewFeatureRegisteredEvent(featureID, featureName, featureType, version string) *FeatureRegisteredEvent {
	return &FeatureRegisteredEvent{
		BaseDomainEvent: NewBaseDomainEvent("FeatureRegistered", featureID, 1),
		FeatureID:       featureID,
		FeatureName:     featureName,
		FeatureType:     featureType,
		Version:         version,
	}
}

// FeatureCompositionRequestedEvent is fired when a feature composition is requested
type FeatureCompositionRequestedEvent struct {
	BaseDomainEvent
	ProjectID      string
	RequestedBy    string
	FeatureIDs     []string
	RequestOptions map[string]interface{}
}

// NewFeatureCompositionRequestedEvent creates a new feature composition requested event
func NewFeatureCompositionRequestedEvent(aggregateID, projectID, requestedBy string, featureIDs []string, options map[string]interface{}) *FeatureCompositionRequestedEvent {
	return &FeatureCompositionRequestedEvent{
		BaseDomainEvent: NewBaseDomainEvent("FeatureCompositionRequested", aggregateID, 1),
		ProjectID:       projectID,
		RequestedBy:     requestedBy,
		FeatureIDs:      featureIDs,
		RequestOptions:  options,
	}
}

// FeatureCompositionCompletedEvent is fired when a feature composition completes
type FeatureCompositionCompletedEvent struct {
	BaseDomainEvent
	ProjectID         string
	ResolvedFeatures  []string
	ConflictCount     int
	WarningCount      int
}

// NewFeatureCompositionCompletedEvent creates a new feature composition completed event
func NewFeatureCompositionCompletedEvent(aggregateID, projectID string, resolvedFeatures []string, conflictCount, warningCount int) *FeatureCompositionCompletedEvent {
	return &FeatureCompositionCompletedEvent{
		BaseDomainEvent:  NewBaseDomainEvent("FeatureCompositionCompleted", aggregateID, 1),
		ProjectID:        projectID,
		ResolvedFeatures: resolvedFeatures,
		ConflictCount:    conflictCount,
		WarningCount:     warningCount,
	}
}

// FeatureCompositionFailedEvent is fired when a feature composition fails
type FeatureCompositionFailedEvent struct {
	BaseDomainEvent
	ProjectID    string
	ErrorCode    string
	ErrorMessage string
}

// NewFeatureCompositionFailedEvent creates a new feature composition failed event
func NewFeatureCompositionFailedEvent(aggregateID, projectID, errorCode, errorMessage string) *FeatureCompositionFailedEvent {
	return &FeatureCompositionFailedEvent{
		BaseDomainEvent: NewBaseDomainEvent("FeatureCompositionFailed", aggregateID, 1),
		ProjectID:       projectID,
		ErrorCode:       errorCode,
		ErrorMessage:    errorMessage,
	}
}

// Template Events

// TemplateCompiledEvent is fired when a template is compiled
type TemplateCompiledEvent struct {
	BaseDomainEvent
	TemplatePath   string
	OutputPath     string
	CompileTime    time.Duration
	Success        bool
}

// NewTemplateCompiledEvent creates a new template compiled event
func NewTemplateCompiledEvent(aggregateID, templatePath, outputPath string, compileTime time.Duration, success bool) *TemplateCompiledEvent {
	return &TemplateCompiledEvent{
		BaseDomainEvent: NewBaseDomainEvent("TemplateCompiled", aggregateID, 1),
		TemplatePath:    templatePath,
		OutputPath:      outputPath,
		CompileTime:     compileTime,
		Success:         success,
	}
}

// TypeSpecGenerationEvent is fired when TypeSpec code generation occurs
type TypeSpecGenerationEvent struct {
	BaseDomainEvent
	SchemaPath     string
	TargetLanguage string
	OutputDir      string
	Success        bool
	GeneratedFiles []string
}

// NewTypeSpecGenerationEvent creates a new TypeSpec generation event
func NewTypeSpecGenerationEvent(aggregateID, schemaPath, targetLanguage, outputDir string, success bool, generatedFiles []string) *TypeSpecGenerationEvent {
	return &TypeSpecGenerationEvent{
		BaseDomainEvent: NewBaseDomainEvent("TypeSpecGeneration", aggregateID, 1),
		SchemaPath:      schemaPath,
		TargetLanguage:  targetLanguage,
		OutputDir:       outputDir,
		Success:         success,
		GeneratedFiles:  generatedFiles,
	}
}

// Validation Events

// ValidationStartedEvent is fired when validation starts
type ValidationStartedEvent struct {
	BaseDomainEvent
	ValidationType string
	TargetID       string
}

// NewValidationStartedEvent creates a new validation started event
func NewValidationStartedEvent(aggregateID, validationType, targetID string) *ValidationStartedEvent {
	return &ValidationStartedEvent{
		BaseDomainEvent: NewBaseDomainEvent("ValidationStarted", aggregateID, 1),
		ValidationType:  validationType,
		TargetID:        targetID,
	}
}

// ValidationCompletedEvent is fired when validation completes
type ValidationCompletedEvent struct {
	BaseDomainEvent
	ValidationType string
	TargetID       string
	Success        bool
	ErrorCount     int
	WarningCount   int
}

// NewValidationCompletedEvent creates a new validation completed event
func NewValidationCompletedEvent(aggregateID, validationType, targetID string, success bool, errorCount, warningCount int) *ValidationCompletedEvent {
	return &ValidationCompletedEvent{
		BaseDomainEvent: NewBaseDomainEvent("ValidationCompleted", aggregateID, 1),
		ValidationType:  validationType,
		TargetID:        targetID,
		Success:         success,
		ErrorCount:      errorCount,
		WarningCount:    warningCount,
	}
}

// Cache Events

// CacheHitEvent is fired when a cache hit occurs
type CacheHitEvent struct {
	BaseDomainEvent
	CacheKey   string
	CacheType  string
	HitRate    float64
}

// NewCacheHitEvent creates a new cache hit event
func NewCacheHitEvent(aggregateID, cacheKey, cacheType string, hitRate float64) *CacheHitEvent {
	return &CacheHitEvent{
		BaseDomainEvent: NewBaseDomainEvent("CacheHit", aggregateID, 1),
		CacheKey:        cacheKey,
		CacheType:       cacheType,
		HitRate:         hitRate,
	}
}

// CacheMissEvent is fired when a cache miss occurs
type CacheMissEvent struct {
	BaseDomainEvent
	CacheKey   string
	CacheType  string
	MissRate   float64
}

// NewCacheMissEvent creates a new cache miss event
func NewCacheMissEvent(aggregateID, cacheKey, cacheType string, missRate float64) *CacheMissEvent {
	return &CacheMissEvent{
		BaseDomainEvent: NewBaseDomainEvent("CacheMiss", aggregateID, 1),
		CacheKey:        cacheKey,
		CacheType:       cacheType,
		MissRate:        missRate,
	}
}

// Performance Events

// PerformanceMetricRecordedEvent is fired when a performance metric is recorded
type PerformanceMetricRecordedEvent struct {
	BaseDomainEvent
	MetricName  string
	MetricValue float64
	MetricUnit  string
	Tags        map[string]string
}

// NewPerformanceMetricRecordedEvent creates a new performance metric recorded event
func NewPerformanceMetricRecordedEvent(aggregateID, metricName string, metricValue float64, metricUnit string, tags map[string]string) *PerformanceMetricRecordedEvent {
	return &PerformanceMetricRecordedEvent{
		BaseDomainEvent: NewBaseDomainEvent("PerformanceMetricRecorded", aggregateID, 1),
		MetricName:      metricName,
		MetricValue:     metricValue,
		MetricUnit:      metricUnit,
		Tags:            tags,
	}
}

// Error Events

// ErrorOccurredEvent is fired when an error occurs in the system
type ErrorOccurredEvent struct {
	BaseDomainEvent
	ErrorCode    string
	ErrorMessage string
	ErrorType    string
	StackTrace   string
	Context      map[string]interface{}
}

// NewErrorOccurredEvent creates a new error occurred event
func NewErrorOccurredEvent(aggregateID, errorCode, errorMessage, errorType, stackTrace string, context map[string]interface{}) *ErrorOccurredEvent {
	return &ErrorOccurredEvent{
		BaseDomainEvent: NewBaseDomainEvent("ErrorOccurred", aggregateID, 1),
		ErrorCode:       errorCode,
		ErrorMessage:    errorMessage,
		ErrorType:       errorType,
		StackTrace:      stackTrace,
		Context:         context,
	}
}

// SecurityEvents

// SecurityEventDetectedEvent is fired when a security event is detected
type SecurityEventDetectedEvent struct {
	BaseDomainEvent
	EventType    string
	Severity     string
	UserID       string
	IPAddress    string
	UserAgent    string
	Resource     string
	Action       string
	Result       string
	Details      map[string]interface{}
}

// NewSecurityEventDetectedEvent creates a new security event detected event
func NewSecurityEventDetectedEvent(aggregateID, eventType, severity, userID, ipAddress, userAgent, resource, action, result string, details map[string]interface{}) *SecurityEventDetectedEvent {
	return &SecurityEventDetectedEvent{
		BaseDomainEvent: NewBaseDomainEvent("SecurityEventDetected", aggregateID, 1),
		EventType:       eventType,
		Severity:        severity,
		UserID:          userID,
		IPAddress:       ipAddress,
		UserAgent:       userAgent,
		Resource:        resource,
		Action:          action,
		Result:          result,
		Details:         details,
	}
}