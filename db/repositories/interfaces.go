package repositories

import "ticken-ticket-service/models/event"

type EventRepository interface {
	AddEvent(event *event.Event) error
	FindEventByID(eventID string) *event.Event
}

type TicketRepository interface {
}
