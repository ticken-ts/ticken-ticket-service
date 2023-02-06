package services

import (
	"errors"
	"github.com/google/uuid"
	"ticken-ticket-service/infra"
	"ticken-ticket-service/infra/public_blockchain"
	"ticken-ticket-service/models"
	"ticken-ticket-service/repos"
)

type userManager struct {
	eventRepository  repos.EventRepository
	ticketRepository repos.TicketRepository
	userRepository   repos.UserRepository
	blockchain       public_blockchain.PublicBC
	hsm              infra.HSM
}

func NewUserManager(
	eventRepository repos.EventRepository,
	ticketRepository repos.TicketRepository,
	userRepository repos.UserRepository,
	blockchain public_blockchain.PublicBC,
	hsm infra.HSM,
) UserManager {
	return &userManager{
		ticketRepository: ticketRepository,
		eventRepository:  eventRepository,
		userRepository:   userRepository,
		blockchain:       blockchain,
		hsm:              hsm,
	}
}

func (userManager *userManager) CreateUser(uuid uuid.UUID, providedPK string) (*models.User, error) {
	newUser := models.NewUser(uuid)
	var pkStoreKey string
	var err error

	// check if user exists
	user := userManager.userRepository.FindUser(uuid)
	if user != nil {
		return nil, errors.New("user already exists")
	}

	if providedPK != "" {
		pkStoreKey, err = userManager.hsm.Store([]byte(providedPK))
		if err != nil {
			return nil, err
		}
	} else {
		newPK, err := userManager.blockchain.GeneratePrivateKey()
		if err != nil {
			return nil, err
		}
		pkStoreKey, err = userManager.hsm.Store([]byte(newPK))
		if err != nil {
			return nil, err
		}
	}

	newUser.SetAddressPKStoreKey(pkStoreKey)
	err = userManager.userRepository.AddUser(newUser)
	if err != nil {
		return nil, err
	}
	return newUser, nil
}
