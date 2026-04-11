package service

import "github.com/google/uuid"

type OrderProviderCreatedPayload struct {
	DraftID         uuid.UUID
	UserID          uuid.UUID
	VenueID         uuid.UUID
	ProviderOrderID uuid.UUID
}

type OrderStatusChangedPayload struct {
	ProviderOrderID uuid.UUID
	Status          string
}
