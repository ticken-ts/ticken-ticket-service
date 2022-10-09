package repos

import (
	"ticken-ticket-service/models"
)

type EventRepository interface {
	AddEvent(event *models.Event) error
	FindEvent(eventID string) *models.Event
}

type TicketRepository interface {
	AddTicket(ticket *models.Ticket) error
	UpdateTicketStatus(ticket *models.Ticket) error
	FindTicket(eventID string, ticketID string) *models.Ticket
}

type Provider interface {
	GetEventRepository() EventRepository
	GetTicketRepository() TicketRepository
}

type Factory interface {
	BuildEventRepository() any
	BuildTicketRepository() any
}