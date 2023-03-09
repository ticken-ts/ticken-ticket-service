package services

import (
	"github.com/google/uuid"
	"ticken-ticket-service/models"
)

type IProvider interface {
	GetTicketIssuer() TicketIssuer
	GetTicketSigner() TicketSigner
	GetEventManager() IEventManager
	GetUserManager() UserManager
}

type TicketIssuer interface {
	IssueTicket(eventID uuid.UUID, section string, owner uuid.UUID) (*models.Ticket, error)
	GetUserTickets(userID uuid.UUID) ([]*models.Ticket, error)
}

type TicketSigner interface {
	SignTicket(eventID string, ticketID string, owner string) (*models.Ticket, error)
}

type IEventManager interface {
	AddEvent(eventID, organizerID uuid.UUID, pvtBCChannel, pubBCAddress string) (*models.Event, error)
}

type UserManager interface {
	// CreateUser creates a new user and returns it.
	// pubBCPrivateKey is the private key of
	// the user in the public blockchain if the
	// user provided one, if is an empty string,
	// a new key is generated
	CreateUser(uuid uuid.UUID, pubBCPrivateKey string) (*models.User, error)
	GetUser(uuid uuid.UUID) (*models.User, error)
}
