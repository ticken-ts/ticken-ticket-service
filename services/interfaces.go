package services

import (
	"github.com/google/uuid"
	"ticken-ticket-service/models"
	"ticken-ticket-service/utils/money"
)

type IProvider interface {
	GetTicketIssuer() TicketIssuer
	GetTicketLinker() TicketLinker
	GetEventManager() IEventManager
	GetUserManager() UserManager
	GetTicketTrader() TicketTrader
}

type TicketTrader interface {
	SellTicket(ownerID, eventID, ticketID uuid.UUID, price *money.Money) (*models.Ticket, error)
	BuyResoldTicket(buyerID, eventID, ticketID, resellID uuid.UUID) (*models.Ticket, error)
}

type TicketIssuer interface {
	IssueTicket(eventID uuid.UUID, section string, owner uuid.UUID) (*models.Ticket, error)
	GetUserTickets(userID uuid.UUID) ([]*models.Ticket, error)
}

type TicketLinker interface {
	LinkTickets(ownerID uuid.UUID, eventContractAddress string) ([]*models.Ticket, error)
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
	GetUserPrivKey(uuid uuid.UUID) (string, error)
}
