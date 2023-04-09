package repos

import (
	"github.com/google/uuid"
	"math/big"
	"ticken-ticket-service/models"
)

type BaseRepository interface {
	Count() int64
	AddOne(element any) error
	AnyWithID(id uuid.UUID) bool
}

type EventRepository interface {
	BaseRepository
	FindEvent(eventID uuid.UUID) *models.Event
	GetActiveEvents() ([]*models.Event, error)
	FindEventByContractAddress(contractAddr string) *models.Event
}

type TicketRepository interface {
	BaseRepository
	UpdateTicketStatus(ticket *models.Ticket) error
	FindTicket(eventID uuid.UUID, ticketID uuid.UUID) *models.Ticket
	GetUserTickets(userID uuid.UUID) ([]*models.Ticket, error)
	UpdateTicketBlockchainData(ticket *models.Ticket) error
	UpdateTicketOwner(ticket *models.Ticket) error
	FindTicketByPUBBCToken(eventID uuid.UUID, token *big.Int) *models.Ticket
	UpdateResoldTicket(ticket *models.Ticket) error
	GetTicketsInResell(eventID uuid.UUID, section string) ([]*models.Ticket, error)
	AddTicketResell(eventID, ticketID uuid.UUID, saleAnnouncement *models.Resell) error
}

type UserRepository interface {
	BaseRepository
	FindUser(userID uuid.UUID) *models.User
}

type IProvider interface {
	GetUserRepository() UserRepository
	GetEventRepository() EventRepository
	GetTicketRepository() TicketRepository
}

type IFactory interface {
	BuildEventRepository() any
	BuildTicketRepository() any
	BuildUserRepository() any
}
