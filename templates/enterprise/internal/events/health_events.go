package events

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

// HealthEventPublisher handles publishing health-related CloudEvents
type HealthEventPublisher struct {
	client cloudevents.Client
}

// NewHealthEventPublisher creates a new health event publisher
func NewHealthEventPublisher() (*HealthEventPublisher, error) {
	// TODO: Configure CloudEvents client with actual transport
	// For now, create a simple client
	client, err := cloudevents.NewClientHTTP()
	if err != nil {
		return nil, fmt.Errorf("failed to create CloudEvents client: %w", err)
	}

	return &HealthEventPublisher{
		client: client,
	}, nil
}

// HealthStatusChangeEvent represents a health status change event
type HealthStatusChangeEvent struct {
	ServiceName    string    `json:"service_name"`
	PreviousStatus string    `json:"previous_status"`
	CurrentStatus  string    `json:"current_status"`
	Timestamp      time.Time `json:"timestamp"`
	Details        string    `json:"details,omitempty"`
}

// DependencyStatusChangeEvent represents a dependency status change event
type DependencyStatusChangeEvent struct {
	ServiceName      string    `json:"service_name"`
	DependencyName   string    `json:"dependency_name"`
	PreviousStatus   string    `json:"previous_status"`
	CurrentStatus    string    `json:"current_status"`
	Timestamp        time.Time `json:"timestamp"`
	ResponseTime     string    `json:"response_time"`
	ErrorMessage     string    `json:"error_message,omitempty"`
}

// PublishHealthStatusChange publishes a health status change event
func (p *HealthEventPublisher) PublishHealthStatusChange(ctx context.Context, serviceName, previousStatus, currentStatus, details string) error {
	event := cloudevents.NewEvent()
	event.SetType("{{.Config.GoModule}}.health.status.changed")
	event.SetSource("{{.Config.Name}}")
	event.SetID(fmt.Sprintf("health-%d", time.Now().UnixNano()))
	event.SetTime(time.Now())

	eventData := HealthStatusChangeEvent{
		ServiceName:    serviceName,
		PreviousStatus: previousStatus,
		CurrentStatus:  currentStatus,
		Timestamp:      time.Now(),
		Details:        details,
	}

	if err := event.SetData(cloudevents.ApplicationJSON, eventData); err != nil {
		return fmt.Errorf("failed to set event data: %w", err)
	}

	// TODO: Send to actual CloudEvents endpoint
	// For now, just log the event
	eventJSON, _ := json.MarshalIndent(eventData, "", "  ")
	log.Printf("CloudEvent: Health Status Change\n%s", eventJSON)

	return nil
}

// PublishDependencyStatusChange publishes a dependency status change event
func (p *HealthEventPublisher) PublishDependencyStatusChange(ctx context.Context, serviceName, dependencyName, previousStatus, currentStatus, responseTime, errorMessage string) error {
	event := cloudevents.NewEvent()
	event.SetType("{{.Config.GoModule}}.dependency.status.changed")
	event.SetSource("{{.Config.Name}}")
	event.SetID(fmt.Sprintf("dependency-%d", time.Now().UnixNano()))
	event.SetTime(time.Now())

	eventData := DependencyStatusChangeEvent{
		ServiceName:      serviceName,
		DependencyName:   dependencyName,
		PreviousStatus:   previousStatus,
		CurrentStatus:    currentStatus,
		Timestamp:        time.Now(),
		ResponseTime:     responseTime,
		ErrorMessage:     errorMessage,
	}

	if err := event.SetData(cloudevents.ApplicationJSON, eventData); err != nil {
		return fmt.Errorf("failed to set event data: %w", err)
	}

	// TODO: Send to actual CloudEvents endpoint
	// For now, just log the event
	eventJSON, _ := json.MarshalIndent(eventData, "", "  ")
	log.Printf("CloudEvent: Dependency Status Change\n%s", eventJSON)

	return nil
}

// PublishStartupComplete publishes a service startup completion event
func (p *HealthEventPublisher) PublishStartupComplete(ctx context.Context, serviceName string, startupTime time.Duration) error {
	event := cloudevents.NewEvent()
	event.SetType("{{.Config.GoModule}}.service.startup.completed")
	event.SetSource("{{.Config.Name}}")
	event.SetID(fmt.Sprintf("startup-%d", time.Now().UnixNano()))
	event.SetTime(time.Now())

	eventData := map[string]interface{}{
		"service_name":   serviceName,
		"startup_time":   startupTime.String(),
		"timestamp":      time.Now(),
		"version":        "{{.Version}}",
	}

	if err := event.SetData(cloudevents.ApplicationJSON, eventData); err != nil {
		return fmt.Errorf("failed to set event data: %w", err)
	}

	// TODO: Send to actual CloudEvents endpoint
	// For now, just log the event
	eventJSON, _ := json.MarshalIndent(eventData, "", "  ")
	log.Printf("CloudEvent: Service Startup Complete\n%s", eventJSON)

	return nil
}
