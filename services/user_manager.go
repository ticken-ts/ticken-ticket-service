package services

import (
	"github.com/google/uuid"
	pubbc "github.com/ticken-ts/ticken-pubbc-connector"
	"ticken-ticket-service/infra"
	"ticken-ticket-service/models"
	"ticken-ticket-service/repos"
	"ticken-ticket-service/tickenerr"
	"ticken-ticket-service/tickenerr/usererr"
)

type UserManager struct {
	eventRepository  repos.EventRepository
	ticketRepository repos.TicketRepository
	userRepository   repos.UserRepository
	pubbcAdmin       pubbc.Admin
	hsm              infra.HSM
}

func NewUserManager(repoProvider repos.IProvider, pubbcAdmin pubbc.Admin, hsm infra.HSM) IUserManager {
	return &UserManager{
		ticketRepository: repoProvider.GetTicketRepository(),
		eventRepository:  repoProvider.GetEventRepository(),
		userRepository:   repoProvider.GetUserRepository(),
		pubbcAdmin:       pubbcAdmin,
		hsm:              hsm,
	}
}

func (userManager *UserManager) CreateUser(attendantID uuid.UUID, providedPK string) (*models.User, error) {
	newAttendant := &models.User{UUID: attendantID}

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
	err = userManager.userRepository.AddOne(newAttendant)
	if err != nil {
		return nil, tickenerr.FromError(usererr.StoreUserInDatabaseErrorCode, err)
	}
	return newAttendant, nil
}

func (userManager *UserManager) GetUser(uuid uuid.UUID) (*models.User, error) {
	user := userManager.userRepository.FindUser(uuid)
	if user == nil {
		return nil, tickenerr.New(usererr.UserNotFoundErrorCode)
	}
	return user, nil
}

func (userManager *UserManager) GetUserPrivKey(uuid uuid.UUID) (string, error) {
	user := userManager.userRepository.FindUser(uuid)
	if user == nil {
		return "", tickenerr.New(usererr.UserNotFoundErrorCode)
	}

	userPrivKey, err := userManager.hsm.Retrieve(user.AddressPKStoreKey)
	if err != nil {
		return "", tickenerr.FromError(usererr.PrivateKeyRetrieveErrorCode, err)
	}

	return string(userPrivKey), err
}
