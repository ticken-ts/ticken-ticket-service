package services

import "ticken-ticket-service/models"

type Provider interface {
	GetTicketIssuer() TicketIssuer
	GetEventManager() EventManager
}

type TicketIssuer interface {
	IssueTicket(eventID string, section string, owner string) (*models.Ticket, error)
}

type EventManager interface {
	AddEvent(EventID string, OrganizerID string, PvtBCChannel string) (*models.Event, error)
}
