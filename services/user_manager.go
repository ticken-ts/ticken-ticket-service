package services

import (
	"errors"
	"github.com/google/uuid"
	pubbc "github.com/ticken-ts/ticken-pubbc-connector"
	"ticken-ticket-service/infra"
	"ticken-ticket-service/models"
	"ticken-ticket-service/repos"
	"ticken-ticket-service/tickenerr"
	"ticken-ticket-service/tickenerr/usererr"
)

type userManager struct {
	eventRepository  repos.EventRepository
	ticketRepository repos.TicketRepository
	userRepository   repos.UserRepository
	pubbcAdmin       pubbc.Admin
	hsm              infra.HSM
}

func NewUserManager(
	eventRepository repos.EventRepository,
	ticketRepository repos.TicketRepository,
	userRepository repos.UserRepository,
	pubbcAdmin pubbc.Admin,
	hsm infra.HSM,
) UserManager {
	return &userManager{
		ticketRepository: ticketRepository,
		eventRepository:  eventRepository,
		userRepository:   userRepository,
		pubbcAdmin:       pubbcAdmin,
		hsm:              hsm,
	}
}

func (userManager *userManager) CreateUser(attendantID uuid.UUID, providedPK string) (*models.User, error) {
	newAttendant := models.NewUser(attendantID)

	var pkStoreKey, walletAddr string
	var err error

	// check if user exists
	user := userManager.userRepository.FindUser(attendantID)
	if user != nil {
		return nil, tickenerr.New(usererr.UserAlreadyExistErrorCode)
	}

	if len(providedPK) > 0 {
		pkStoreKey, err = userManager.hsm.Store([]byte(providedPK))
		if err != nil {
			return nil, tickenerr.FromError(usererr.PrivateKeyStoreErrorCode, err)
		}
	} else {
		newPK, newWalletAddr, err := userManager.pubbcAdmin.CreateWallet()
		if err != nil {
			return nil, tickenerr.FromError(usererr.CreateWallerErrorCode, err)
		}
		pkStoreKey, err = userManager.hsm.Store([]byte(newPK))
		if err != nil {
			return nil, tickenerr.FromError(usererr.PrivateKeyStoreErrorCode, err)
		}
		walletAddr = newWalletAddr
	}

	newAttendant.SetWallet(pkStoreKey, walletAddr)
	err = userManager.userRepository.AddUser(newAttendant)
	if err != nil {
		return nil, tickenerr.FromError(usererr.StoreUserInDatabase, err)
	}
	return newAttendant, nil
}

// GetUser returns the user with the given UUID
func (userManager *userManager) GetUser(uuid uuid.UUID) (*models.User, error) {
	user := userManager.userRepository.FindUser(uuid)
	if user == nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (userManager *userManager) GetUserPrivKey(uuid uuid.UUID) (string, error) {
	user := userManager.userRepository.FindUser(uuid)
	if user == nil {
		return "", errors.New("user not found")
	}

	userPrivKey, err := userManager.hsm.Retrieve(user.AddressPKStoreKey)
	if err != nil {
		return "", err
	}

	return string(userPrivKey), err
}
