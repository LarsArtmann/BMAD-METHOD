package events

// Simplified events for intermediate tier
type EventEmitter struct{}

func NewEventEmitter(serviceName, sinkURL string) *EventEmitter {
	return &EventEmitter{}
}
