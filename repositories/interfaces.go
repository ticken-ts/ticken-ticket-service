package repositories

import (
	"ticken-ticket-service/models"
)

type EventRepository interface {
	AddEvent(event *models.Event) error
	FindEventByID(eventID string) *models.Event
}

type TicketRepository interface {
}

type Provider interface {
}
