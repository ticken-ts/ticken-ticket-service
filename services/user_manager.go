package services

import (
	"ticken-ticket-service/infra/public_blockchain"
	"ticken-ticket-service/models"
	"ticken-ticket-service/repos"
)

type userManager struct {
	eventRepository  repos.EventRepository
	ticketRepository repos.TicketRepository
	blockchain       public_blockchain.PublicBC
}

func NewUserManager(
	eventRepository repos.EventRepository,
	ticketRepository repos.TicketRepository,
	blockchain public_blockchain.PublicBC,
) UserManager {
	return &userManager{
		ticketRepository: ticketRepository,
		eventRepository:  eventRepository,
		blockchain:       blockchain,
	}
}

func (userManager *userManager) CreateUser() (*models.User, error) {
	return new(models.User), nil
}
