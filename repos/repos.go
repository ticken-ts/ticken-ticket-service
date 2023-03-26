package repos

import (
	"github.com/google/uuid"
	"math/big"
	"ticken-ticket-service/models"
)

type EventRepository interface {
	AddEvent(event *models.Event) error
	FindEvent(eventID uuid.UUID) *models.Event
	GetActiveEvents() ([]*models.Event, error)
	FindEventByContractAddress(contractAddr string) *models.Event
}

type TicketRepository interface {
	AddTicket(ticket *models.Ticket) error
	UpdateTicketStatus(ticket *models.Ticket) error
	FindTicket(eventID uuid.UUID, ticketID uuid.UUID) *models.Ticket
	GetUserTickets(userID uuid.UUID) ([]*models.Ticket, error)
	UpdateTicketBlockchainData(ticket *models.Ticket) error
	UpdateTicketOwner(ticket *models.Ticket) error
	FindTicketByPUBBCToken(eventID uuid.UUID, token *big.Int) *models.Ticket
	UpdateResoldTicket(ticket *models.Ticket) error
	AddTicketResell(eventID, ticketID uuid.UUID, saleAnnouncement *models.Resell) error
}

type UserRepository interface {
	AddUser(user *models.User) error
	FindUser(userID uuid.UUID) *models.User
}

type IProvider interface {
	GetEventRepository() EventRepository
	GetTicketRepository() TicketRepository
	GetUserRepository() UserRepository
}

type IFactory interface {
	BuildEventRepository() any
	BuildTicketRepository() any
	BuildUserRepository() any
}
