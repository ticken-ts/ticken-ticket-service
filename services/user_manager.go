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

func (userManager *userManager) CreateUser(uuid uuid.UUID, providedPK string) (*models.User, error) {
	newUser := models.NewUser(uuid)
	if providedPK != "" {
		newUser.SetPubBCPrivateKey(providedPK)
	} else {
		newPK, err := userManager.blockchain.GeneratePrivateKey()
		if err != nil {
			return nil, err
		}
		newUser.SetPubBCPrivateKey(newPK)
	}
	err := userManager.userRepository.AddUser(newUser)
	if err != nil {
		return nil, err
	}
	return newUser, nil
}
