# Domain-Driven Design and Clean Architecture Refactoring

## Objective
Refactor an existing codebase to implement Domain-Driven Design (DDD) principles with Clean Architecture patterns, establishing clear boundaries, proper domain modeling, and enterprise-grade architectural foundations.

## Context
This prompt focuses on transforming a procedural or service-oriented codebase into a well-structured domain-centric architecture that supports complex business logic, scalability, and maintainability through proper separation of concerns.

## Task Description

### Phase 1: Domain Analysis and Modeling
1. **Domain Discovery**
   - Identify core business concepts and entities within the application
   - Map business processes and workflows
   - Discover domain experts and stakeholder requirements
   - Establish ubiquitous language for consistent terminology

2. **Bounded Context Definition**
   - Define clear boundaries between different domain areas
   - Identify context maps and relationships between contexts
   - Establish integration patterns between bounded contexts
   - Define shared kernels and anti-corruption layers

3. **Aggregate Design**
   - Identify aggregate roots and their boundaries
   - Define entity relationships and value objects
   - Establish invariants and business rules
   - Design for consistency and transactional boundaries

### Phase 2: Core Domain Implementation
1. **Domain Entities and Value Objects**
   - Create `pkg/domain/entities.go` with proper entity design
   - Implement aggregate roots with domain event support
   - Build value objects with immutability and validation
   - Establish entity identity and lifecycle management

2. **Domain Events System**
   - Create `pkg/domain/events.go` with comprehensive event modeling
   - Implement event sourcing patterns where appropriate
   - Build event handlers and subscribers
   - Establish event versioning and backward compatibility

3. **Repository Patterns**
   - Create `pkg/domain/repositories.go` with repository interfaces
   - Define specifications for complex queries
   - Implement unit of work patterns for transactions
   - Build read model repositories for CQRS support

### Phase 3: Clean Architecture Layers
1. **Application Services Layer**
   - Create `pkg/application/` for use cases and application logic
   - Implement command and query handlers (CQRS)
   - Build application services for coordinating domain logic
   - Establish transaction boundaries and error handling

2. **Infrastructure Layer**
   - Create `pkg/infrastructure/` for external concerns
   - Implement repository implementations with database access
   - Build external service integrations and adapters
   - Establish configuration and dependency injection

3. **Interface Layer**
   - Create `pkg/interfaces/` for external interfaces
   - Implement REST API controllers and handlers
   - Build CLI interfaces and command processors
   - Establish serialization and validation layers

### Phase 4: Advanced Patterns Implementation
1. **CQRS (Command Query Responsibility Segregation)**
   - Separate read and write models for optimal performance
   - Implement command handlers for write operations
   - Build query handlers for read operations
   - Establish eventual consistency patterns

2. **Event Sourcing**
   - Store domain events as the source of truth
   - Implement event store with proper versioning
   - Build event replay and projection mechanisms
   - Establish snapshotting for performance optimization

3. **Saga Patterns**
   - Implement long-running business processes
   - Build compensation logic for distributed transactions
   - Establish timeout and retry mechanisms
   - Create process managers for complex workflows

## Architecture Structure

### Directory Layout
```
pkg/
├── domain/                 # Core domain logic
│   ├── entities.go        # Domain entities and aggregates
│   ├── events.go          # Domain events
│   ├── repositories.go    # Repository interfaces
│   └── services.go        # Domain services
├── application/           # Application layer
│   ├── commands/          # Command handlers
│   ├── queries/           # Query handlers
│   ├── services/          # Application services
│   └── usecases/          # Use case implementations
├── infrastructure/       # Infrastructure layer
│   ├── persistence/       # Database implementations
│   ├── messaging/         # Event bus implementations
│   ├── external/          # External service adapters
│   └── config/           # Configuration management
└── interfaces/           # Interface layer
    ├── http/             # REST API handlers
    ├── cli/              # Command-line interfaces
    └── events/           # Event handlers
```

### Core Patterns Implementation

#### 1. Aggregate Root Pattern
```go
// BaseAggregateRoot provides common aggregate functionality
type BaseAggregateRoot struct {
    BaseEntity
    domainEvents []DomainEvent
    version      int
}

func (ar *BaseAggregateRoot) AddDomainEvent(event DomainEvent) {
    ar.domainEvents = append(ar.domainEvents, event)
}

func (ar *BaseAggregateRoot) DomainEvents() []DomainEvent {
    return ar.domainEvents
}
```

#### 2. Repository Pattern
```go
// Repository interface for aggregate persistence
type Repository[T Entity] interface {
    Save(ctx context.Context, entity T) error
    FindByID(ctx context.Context, id string) (T, error)
    Delete(ctx context.Context, id string) error
}

// Specification pattern for complex queries
type Specification[T Entity] interface {
    IsSatisfiedBy(entity T) bool
    ToSQL() (string, []interface{}, error)
}
```

#### 3. Domain Events
```go
// DomainEvent interface for all domain events
type DomainEvent interface {
    OccurredAt() time.Time
    EventType() string
    AggregateID() string
    Version() int
}

// Event store for persistence
type EventStore interface {
    AppendEvents(ctx context.Context, aggregateID string, 
                expectedVersion int, events []DomainEvent) error
    GetEvents(ctx context.Context, aggregateID string) ([]DomainEvent, error)
}
```

## Implementation Guidelines

### Domain Layer Rules
1. **Pure Domain Logic**: No infrastructure dependencies
2. **Rich Domain Models**: Behavior-rich entities and value objects
3. **Invariant Protection**: Enforce business rules at the domain level
4. **Event-Driven**: Use domain events for loose coupling
5. **Ubiquitous Language**: Consistent terminology throughout

### Application Layer Principles
1. **Orchestration**: Coordinate domain objects without business logic
2. **Transaction Management**: Define transaction boundaries
3. **DTO Mapping**: Transform between domain and external models
4. **Validation**: Input validation and authorization
5. **Error Handling**: Consistent error processing and logging

### Infrastructure Layer Guidelines
1. **Inversion of Control**: Implement domain interfaces
2. **Configuration**: Externalize all configuration
3. **Logging and Monitoring**: Comprehensive observability
4. **Performance**: Optimize for specific use cases
5. **Security**: Implement security concerns

## Expected Deliverables
1. **Domain Model Implementation**
   - Complete entity and value object design
   - Domain events system with proper versioning
   - Repository interfaces with specifications
   - Domain services for complex business logic

2. **Clean Architecture Structure**
   - Properly layered application with clear boundaries
   - Dependency inversion with interface-based design
   - CQRS implementation for read/write separation
   - Event sourcing for audit and temporal queries

3. **Infrastructure Implementation**
   - Database repository implementations
   - Event store and message bus integration
   - External service adapters and anti-corruption layers
   - Configuration and dependency injection framework

4. **Testing Strategy**
   - Unit tests for domain logic with >95% coverage
   - Integration tests for repository implementations
   - Acceptance tests for use cases and scenarios
   - Performance tests for critical paths

## Success Criteria
- **Separation of Concerns**: Clear architectural boundaries with minimal coupling
- **Domain Richness**: Business logic properly encapsulated in domain layer
- **Testability**: High test coverage with fast, reliable tests
- **Maintainability**: Easy to understand and modify codebase
- **Performance**: Meets or exceeds existing performance benchmarks
- **Scalability**: Architecture supports horizontal and vertical scaling

## Quality Metrics
- **Cyclomatic Complexity**: <10 for individual methods
- **Test Coverage**: >90% for domain layer, >80% overall
- **Dependency Metrics**: Minimal coupling between layers
- **Performance**: <5% degradation from baseline
- **Documentation**: Complete API documentation and architectural decision records

## Migration Strategy
1. **Strangler Fig Pattern**: Gradually replace existing code
2. **Feature Toggles**: Safe deployment of new implementations
3. **Data Migration**: Preserve existing data during transition
4. **Rollback Plan**: Ability to revert to previous implementation
5. **Monitoring**: Comprehensive observability during migration

## Related Patterns
- **Event Sourcing**: Store events as source of truth
- **CQRS**: Separate read and write models
- **Saga Pattern**: Manage distributed transactions
- **Anti-Corruption Layer**: Protect domain from external models
- **Published Language**: Standardize integration contracts