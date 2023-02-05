package services

import (
	"github.com/google/uuid"
	"ticken-ticket-service/models"
)

type IProvider interface {
	GetTicketIssuer() TicketIssuer
	GetTicketSigner() TicketSigner
	GetEventManager() EventManager
}

type TicketIssuer interface {
	IssueTicket(eventID uuid.UUID, section string, owner uuid.UUID) (*models.Ticket, error)
}

type TicketSigner interface {
	SignTicket(eventID string, ticketID string, owner string) (*models.Ticket, error)
}

type EventManager interface {
	AddEvent(EventID string, OrganizerID string, PvtBCChannel string) (*models.Event, error)
}

type UserManager interface {
	CreateUser() (*models.User, error)
}
