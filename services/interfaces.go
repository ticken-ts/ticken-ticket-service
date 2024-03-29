package services

import (
	"github.com/google/uuid"
	"ticken-ticket-service/models"
	"ticken-ticket-service/utils/money"
)

type IProvider interface {
	GetTicketIssuer() ITicketIssuer
	GetTicketLinker() ITicketLinker
	GetEventManager() IEventManager
	GetUserManager() IUserManager
	GetTicketTrader() ITicketTrader
}

type ITicketTrader interface {
	SellTicket(ownerID, eventID, ticketID uuid.UUID, price *money.Money) (*models.Ticket, error)
	BuyResoldTicket(buyerID, eventID, ticketID, resellID uuid.UUID) (*models.Ticket, error)
	GetTicketsInResells(eventID uuid.UUID, section string) ([]*models.Ticket, error)
}

type ITicketIssuer interface {
	IssueTicket(eventID uuid.UUID, section string, owner uuid.UUID) (*models.Ticket, error)
	GetUserTickets(userID uuid.UUID) ([]*models.Ticket, error)
}

type ITicketLinker interface {
	LinkTickets(ownerID uuid.UUID, eventContractAddress string) ([]*models.Ticket, error)
}

type IEventManager interface {
	AddEvent(eventID, organizerID uuid.UUID, pvtBCChannel, pubBCAddress string) (*models.Event, error)
}

type IUserManager interface {
	// CreateUser creates a new user and returns it.
	// pubBCPrivateKey is the private key of
	// the user in the public blockchain if the
	// user provided one, if is an empty string,
	// a new key is generated
	CreateUser(uuid uuid.UUID, pubBCPrivateKey string) (*models.Attendant, error)

	// todo
	RegisterUser(email, password, firstname, lastname, providedPK string) (*models.Attendant, error)

	GetUser(uuid uuid.UUID) (*models.Attendant, error)

	// format -> hex / pem
	GetUserPrivKey(uuid uuid.UUID, format string) (string, error)
}
