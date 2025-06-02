package domain

import (
	"context"
	"errors"
	"time"
)

// Repository errors
var (
	ErrEntityNotFound      = errors.New("entity not found")
	ErrEntityAlreadyExists = errors.New("entity already exists")
	ErrConcurrencyConflict = errors.New("concurrency conflict")
	ErrInvalidQuery        = errors.New("invalid query")
)

// Repository represents a generic repository interface
type Repository[T Entity] interface {
	Save(ctx context.Context, entity T) error
	FindByID(ctx context.Context, id string) (T, error)
	Delete(ctx context.Context, id string) error
	Exists(ctx context.Context, id string) (bool, error)
}

// Specification represents a query specification
type Specification[T Entity] interface {
	IsSatisfiedBy(entity T) bool
	ToSQL() (string, []interface{}, error)
}

// QueryRepository extends Repository with query capabilities
type QueryRepository[T Entity] interface {
	Repository[T]
	FindBy(ctx context.Context, spec Specification[T]) ([]T, error)
	FindAll(ctx context.Context, options QueryOptions) ([]T, error)
	Count(ctx context.Context, spec Specification[T]) (int64, error)
}

// QueryOptions represents options for queries
type QueryOptions struct {
	Limit      int
	Offset     int
	SortBy     string
	SortOrder  SortOrder
	Filters    map[string]interface{}
}

// SortOrder represents sort order
type SortOrder string

const (
	SortOrderAsc  SortOrder = "ASC"
	SortOrderDesc SortOrder = "DESC"
)

// ProjectRepository manages Project aggregates
type ProjectRepository interface {
	QueryRepository[*Project]
	
	// Project-specific methods
	FindByName(ctx context.Context, name string) (*Project, error)
	FindByOwner(ctx context.Context, owner string) ([]*Project, error)
	FindByTier(ctx context.Context, tier string) ([]*Project, error)
	FindByStatus(ctx context.Context, status ProjectStatus) ([]*Project, error)
	
	// Advanced queries
	FindProjectsWithFeature(ctx context.Context, featureID string) ([]*Project, error)
	FindRecentProjects(ctx context.Context, limit int) ([]*Project, error)
	SearchProjects(ctx context.Context, query string, options QueryOptions) ([]*Project, error)
}

// GenerationRepository manages Generation aggregates
type GenerationRepository interface {
	QueryRepository[*Generation]
	
	// Generation-specific methods
	FindByProjectID(ctx context.Context, projectID string) ([]*Generation, error)
	FindByStatus(ctx context.Context, status GenerationStatus) ([]*Generation, error)
	FindByRequestedBy(ctx context.Context, requestedBy string) ([]*Generation, error)
	
	// Advanced queries
	FindActiveGenerations(ctx context.Context) ([]*Generation, error)
	FindFailedGenerations(ctx context.Context, since *time.Time) ([]*Generation, error)
	FindGenerationStats(ctx context.Context, projectID string) (*GenerationStats, error)
}

// GenerationStats represents generation statistics
type GenerationStats struct {
	TotalGenerations    int64
	SuccessfulGenerations int64
	FailedGenerations   int64
	AverageGenerationTime time.Duration
	LastGenerationTime  *time.Time
}

// FeatureRepository manages feature metadata (not aggregates, but supporting data)
type FeatureRepository interface {
	// Feature registry operations
	SaveFeatureMetadata(ctx context.Context, feature *FeatureMetadata) error
	FindFeatureMetadata(ctx context.Context, featureID string) (*FeatureMetadata, error)
	ListAvailableFeatures(ctx context.Context, tier string) ([]*FeatureMetadata, error)
	
	// Feature usage tracking
	RecordFeatureUsage(ctx context.Context, usage *FeatureUsage) error
	GetFeatureUsageStats(ctx context.Context, featureID string) (*FeatureUsageStats, error)
	GetPopularFeatures(ctx context.Context, limit int) ([]*FeatureMetadata, error)
}

// FeatureMetadata represents metadata about available features
type FeatureMetadata struct {
	ID            string                 `json:"id"`
	Name          string                 `json:"name"`
	Description   string                 `json:"description"`
	Type          string                 `json:"type"`
	Version       string                 `json:"version"`
	Category      string                 `json:"category"`
	Tags          []string               `json:"tags"`
	MinTier       string                 `json:"min_tier"`
	MaxTier       string                 `json:"max_tier"`
	Dependencies  []string               `json:"dependencies"`
	Conflicts     []string               `json:"conflicts"`
	ConfigSchema  map[string]interface{} `json:"config_schema"`
	DefaultConfig map[string]interface{} `json:"default_config"`
	Deprecated    bool                   `json:"deprecated"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
}

// FeatureUsage represents feature usage data
type FeatureUsage struct {
	ID         string    `json:"id"`
	FeatureID  string    `json:"feature_id"`
	ProjectID  string    `json:"project_id"`
	UserID     string    `json:"user_id"`
	UsedAt     time.Time `json:"used_at"`
	Success    bool      `json:"success"`
	Duration   time.Duration `json:"duration"`
	ErrorCode  *string   `json:"error_code,omitempty"`
	Metadata   map[string]interface{} `json:"metadata"`
}

// FeatureUsageStats represents aggregated feature usage statistics
type FeatureUsageStats struct {
	FeatureID        string        `json:"feature_id"`
	TotalUsages      int64         `json:"total_usages"`
	SuccessfulUsages int64         `json:"successful_usages"`
	FailedUsages     int64         `json:"failed_usages"`
	AverageDuration  time.Duration `json:"average_duration"`
	LastUsed         *time.Time    `json:"last_used"`
	PopularityRank   int           `json:"popularity_rank"`
}

// EventStore manages domain events
type EventStore interface {
	// Event storage
	AppendEvents(ctx context.Context, aggregateID string, expectedVersion int, events []DomainEvent) error
	GetEvents(ctx context.Context, aggregateID string) ([]DomainEvent, error)
	GetEventsAfterVersion(ctx context.Context, aggregateID string, version int) ([]DomainEvent, error)
	
	// Event querying
	GetEventsByType(ctx context.Context, eventType string, limit int) ([]DomainEvent, error)
	GetEventsByTimeRange(ctx context.Context, start, end time.Time) ([]DomainEvent, error)
	GetAllEvents(ctx context.Context, options QueryOptions) ([]DomainEvent, error)
	
	// Snapshots
	SaveSnapshot(ctx context.Context, aggregateID string, version int, data []byte) error
	GetSnapshot(ctx context.Context, aggregateID string) (*Snapshot, error)
}

// Snapshot represents an aggregate snapshot
type Snapshot struct {
	AggregateID string    `json:"aggregate_id"`
	Version     int       `json:"version"`
	Data        []byte    `json:"data"`
	CreatedAt   time.Time `json:"created_at"`
}

// UnitOfWork manages transactional boundaries
type UnitOfWork interface {
	Begin(ctx context.Context) error
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
	
	// Repository access within transaction
	ProjectRepository() ProjectRepository
	GenerationRepository() GenerationRepository
	FeatureRepository() FeatureRepository
	EventStore() EventStore
}

// ReadModelRepository manages read models for CQRS
type ReadModelRepository interface {
	// Project read models
	GetProjectSummary(ctx context.Context, projectID string) (*ProjectSummary, error)
	GetProjectsList(ctx context.Context, options ProjectListOptions) (*ProjectListResult, error)
	GetProjectDashboard(ctx context.Context, userID string) (*ProjectDashboard, error)
	
	// Generation read models
	GetGenerationHistory(ctx context.Context, projectID string, options QueryOptions) (*GenerationHistoryResult, error)
	GetGenerationProgress(ctx context.Context, generationID string) (*GenerationProgress, error)
	
	// Feature read models
	GetFeatureCatalog(ctx context.Context, tier string) (*FeatureCatalog, error)
	GetFeatureRecommendations(ctx context.Context, projectID string) (*FeatureRecommendations, error)
	
	// Analytics read models
	GetUsageAnalytics(ctx context.Context, options AnalyticsOptions) (*UsageAnalytics, error)
	GetPerformanceMetrics(ctx context.Context, options MetricsOptions) (*PerformanceMetrics, error)
}

// Read model data structures

// ProjectSummary represents a project summary read model
type ProjectSummary struct {
	ID            string                 `json:"id"`
	Name          string                 `json:"name"`
	Description   string                 `json:"description"`
	Tier          string                 `json:"tier"`
	Status        string                 `json:"status"`
	Owner         string                 `json:"owner"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
	FeatureCount  int                    `json:"feature_count"`
	EnabledFeatures []string             `json:"enabled_features"`
	LastGeneration *time.Time           `json:"last_generation"`
	GenerationCount int                  `json:"generation_count"`
	Metadata      map[string]interface{} `json:"metadata"`
}

// ProjectListOptions represents options for listing projects
type ProjectListOptions struct {
	QueryOptions
	Owner    *string `json:"owner,omitempty"`
	Tier     *string `json:"tier,omitempty"`
	Status   *string `json:"status,omitempty"`
	HasFeature *string `json:"has_feature,omitempty"`
}

// ProjectListResult represents the result of listing projects
type ProjectListResult struct {
	Projects    []*ProjectSummary `json:"projects"`
	TotalCount  int64             `json:"total_count"`
	HasMore     bool              `json:"has_more"`
	NextOffset  int               `json:"next_offset"`
}

// ProjectDashboard represents a user's project dashboard
type ProjectDashboard struct {
	UserID           string             `json:"user_id"`
	RecentProjects   []*ProjectSummary  `json:"recent_projects"`
	ActiveGenerations []*GenerationProgress `json:"active_generations"`
	Stats            *DashboardStats    `json:"stats"`
	Recommendations  *FeatureRecommendations `json:"recommendations"`
	LastUpdated      time.Time          `json:"last_updated"`
}

// DashboardStats represents dashboard statistics
type DashboardStats struct {
	TotalProjects      int   `json:"total_projects"`
	ProjectsByTier     map[string]int `json:"projects_by_tier"`
	TotalGenerations   int   `json:"total_generations"`
	SuccessfulGenerations int `json:"successful_generations"`
	FailedGenerations  int   `json:"failed_generations"`
	MostUsedFeatures   []string `json:"most_used_features"`
}

// GenerationHistoryResult represents generation history
type GenerationHistoryResult struct {
	Generations []*GenerationSummary `json:"generations"`
	TotalCount  int64                `json:"total_count"`
	HasMore     bool                 `json:"has_more"`
	Stats       *GenerationStats     `json:"stats"`
}

// GenerationSummary represents a generation summary
type GenerationSummary struct {
	ID          string                 `json:"id"`
	ProjectID   string                 `json:"project_id"`
	Status      string                 `json:"status"`
	RequestedBy string                 `json:"requested_by"`
	StartedAt   time.Time              `json:"started_at"`
	CompletedAt *time.Time             `json:"completed_at"`
	Duration    *time.Duration         `json:"duration"`
	ArtifactCount int                  `json:"artifact_count"`
	ErrorCount  int                    `json:"error_count"`
	Features    []string               `json:"features"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// GenerationProgress represents real-time generation progress
type GenerationProgress struct {
	GenerationID string                 `json:"generation_id"`
	ProjectID    string                 `json:"project_id"`
	Status       string                 `json:"status"`
	Progress     float64                `json:"progress"` // 0.0 to 1.0
	CurrentStep  string                 `json:"current_step"`
	TotalSteps   int                    `json:"total_steps"`
	CompletedSteps int                  `json:"completed_steps"`
	EstimatedTimeRemaining *time.Duration `json:"estimated_time_remaining"`
	StartedAt    time.Time              `json:"started_at"`
	LastUpdate   time.Time              `json:"last_update"`
}

// FeatureCatalog represents the feature catalog
type FeatureCatalog struct {
	Categories map[string][]*FeatureMetadata `json:"categories"`
	TotalCount int                           `json:"total_count"`
	TierStats  map[string]int                `json:"tier_stats"`
	NewFeatures []*FeatureMetadata           `json:"new_features"`
	PopularFeatures []*FeatureMetadata       `json:"popular_features"`
	LastUpdated time.Time                   `json:"last_updated"`
}

// FeatureRecommendations represents feature recommendations for a project
type FeatureRecommendations struct {
	ProjectID    string             `json:"project_id"`
	Recommended  []*FeatureMetadata `json:"recommended"`
	Complementary []*FeatureMetadata `json:"complementary"`
	Trending     []*FeatureMetadata `json:"trending"`
	Reasoning    map[string]string  `json:"reasoning"`
	Confidence   map[string]float64 `json:"confidence"`
	GeneratedAt  time.Time          `json:"generated_at"`
}

// Analytics and metrics data structures

// AnalyticsOptions represents options for analytics queries
type AnalyticsOptions struct {
	TimeRange   TimeRange `json:"time_range"`
	Granularity string    `json:"granularity"` // hour, day, week, month
	Filters     map[string]interface{} `json:"filters"`
	GroupBy     []string  `json:"group_by"`
}

// TimeRange represents a time range for analytics
type TimeRange struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

// UsageAnalytics represents usage analytics data
type UsageAnalytics struct {
	TimeRange        TimeRange              `json:"time_range"`
	ProjectCreations []DataPoint            `json:"project_creations"`
	Generations      []DataPoint            `json:"generations"`
	FeatureUsage     map[string][]DataPoint `json:"feature_usage"`
	UserActivity     []DataPoint            `json:"user_activity"`
	PopularTiers     map[string]int         `json:"popular_tiers"`
	TrendingFeatures []*FeatureMetadata     `json:"trending_features"`
}

// MetricsOptions represents options for performance metrics
type MetricsOptions struct {
	AnalyticsOptions
	MetricTypes []string `json:"metric_types"`
	Percentiles []float64 `json:"percentiles"`
}

// PerformanceMetrics represents performance metrics data
type PerformanceMetrics struct {
	TimeRange         TimeRange              `json:"time_range"`
	GenerationTimes   []DataPoint            `json:"generation_times"`
	SuccessRates      []DataPoint            `json:"success_rates"`
	ErrorRates        []DataPoint            `json:"error_rates"`
	ThroughputMetrics []DataPoint            `json:"throughput_metrics"`
	ResourceUsage     map[string][]DataPoint `json:"resource_usage"`
	Percentiles       map[string]map[float64]float64 `json:"percentiles"`
}

// DataPoint represents a single data point in time series data
type DataPoint struct {
	Timestamp time.Time   `json:"timestamp"`
	Value     float64     `json:"value"`
	Labels    map[string]string `json:"labels,omitempty"`
}

// Audit and compliance repositories

// AuditRepository manages audit trails
type AuditRepository interface {
	RecordAction(ctx context.Context, action *AuditAction) error
	GetAuditTrail(ctx context.Context, options AuditQueryOptions) (*AuditTrail, error)
	GetUserActions(ctx context.Context, userID string, options QueryOptions) ([]*AuditAction, error)
	GetResourceActions(ctx context.Context, resourceType, resourceID string, options QueryOptions) ([]*AuditAction, error)
}

// AuditAction represents an auditable action
type AuditAction struct {
	ID          string                 `json:"id"`
	UserID      string                 `json:"user_id"`
	Action      string                 `json:"action"`
	ResourceType string                `json:"resource_type"`
	ResourceID  string                 `json:"resource_id"`
	IPAddress   string                 `json:"ip_address"`
	UserAgent   string                 `json:"user_agent"`
	Timestamp   time.Time              `json:"timestamp"`
	Success     bool                   `json:"success"`
	ErrorCode   *string                `json:"error_code,omitempty"`
	Details     map[string]interface{} `json:"details"`
	Risk        RiskLevel              `json:"risk"`
}

// RiskLevel represents the risk level of an action
type RiskLevel string

const (
	RiskLevelLow    RiskLevel = "low"
	RiskLevelMedium RiskLevel = "medium"
	RiskLevelHigh   RiskLevel = "high"
	RiskLevelCritical RiskLevel = "critical"
)

// AuditQueryOptions represents options for audit queries
type AuditQueryOptions struct {
	QueryOptions
	UserID       *string    `json:"user_id,omitempty"`
	Action       *string    `json:"action,omitempty"`
	ResourceType *string    `json:"resource_type,omitempty"`
	ResourceID   *string    `json:"resource_id,omitempty"`
	RiskLevel    *RiskLevel `json:"risk_level,omitempty"`
	Success      *bool      `json:"success,omitempty"`
	TimeRange    *TimeRange `json:"time_range,omitempty"`
}

// AuditTrail represents an audit trail result
type AuditTrail struct {
	Actions     []*AuditAction `json:"actions"`
	TotalCount  int64          `json:"total_count"`
	HasMore     bool           `json:"has_more"`
	RiskSummary map[RiskLevel]int `json:"risk_summary"`
}

// Cache repository for performance optimization

// CacheRepository manages caching for read models and frequently accessed data
type CacheRepository interface {
	// Generic cache operations
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Get(ctx context.Context, key string, dest interface{}) error
	Delete(ctx context.Context, key string) error
	DeletePattern(ctx context.Context, pattern string) error
	
	// Specific cache methods
	CacheProjectSummary(ctx context.Context, projectID string, summary *ProjectSummary, ttl time.Duration) error
	GetProjectSummary(ctx context.Context, projectID string) (*ProjectSummary, error)
	
	CacheFeatureCatalog(ctx context.Context, tier string, catalog *FeatureCatalog, ttl time.Duration) error
	GetFeatureCatalog(ctx context.Context, tier string) (*FeatureCatalog, error)
	
	InvalidateProjectCache(ctx context.Context, projectID string) error
	InvalidateFeatureCache(ctx context.Context) error
	
	// Cache statistics
	GetCacheStats(ctx context.Context) (*CacheStats, error)
}

// CacheStats represents cache performance statistics
type CacheStats struct {
	HitRate     float64 `json:"hit_rate"`
	MissRate    float64 `json:"miss_rate"`
	TotalHits   int64   `json:"total_hits"`
	TotalMisses int64   `json:"total_misses"`
	KeyCount    int64   `json:"key_count"`
	MemoryUsage int64   `json:"memory_usage"`
}