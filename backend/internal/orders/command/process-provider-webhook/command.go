package process_provider_webhook

import "github.com/google/uuid"

type Command struct {
	EventType       string
	CorrelationID   *uuid.UUID
	ProviderOrderID *uuid.UUID
	Status          string
	CreationStatus  string
	Error           *string
}
