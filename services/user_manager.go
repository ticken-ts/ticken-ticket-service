package services

import (
	"github.com/google/uuid"
	"ticken-ticket-service/infra/public_blockchain"
	"ticken-ticket-service/models"
	"ticken-ticket-service/repos"
)

type userManager struct {
	eventRepository  repos.EventRepository
	ticketRepository repos.TicketRepository
	userRepository   repos.UserRepository
	blockchain       public_blockchain.PublicBC
}

func NewUserManager(
	eventRepository repos.EventRepository,
	ticketRepository repos.TicketRepository,
	userRepository repos.UserRepository,
	blockchain public_blockchain.PublicBC,
) UserManager {
	return &userManager{
		ticketRepository: ticketRepository,
		eventRepository:  eventRepository,
		userRepository:   userRepository,
		blockchain:       blockchain,
	}
}

func (userManager *userManager) CreateUser(uuid uuid.UUID) (*models.User, error) {
	newUser := models.NewUser(uuid)
	err := userManager.userRepository.AddUser(newUser)
	if err != nil {
		return nil, err
	}
	return newUser, nil
}
