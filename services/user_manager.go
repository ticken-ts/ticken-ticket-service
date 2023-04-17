package services

import (
	"github.com/google/uuid"
	pubbc "github.com/ticken-ts/ticken-pubbc-connector"
	"ticken-ticket-service/infra"
	"ticken-ticket-service/models"
	"ticken-ticket-service/repos"
	"ticken-ticket-service/sync"
	"ticken-ticket-service/tickenerr"
	"ticken-ticket-service/tickenerr/usererr"
)

type UserManager struct {
	eventRepository  repos.EventRepository
	ticketRepository repos.TicketRepository
	userRepository   repos.UserRepository
	pubbcAdmin       pubbc.Admin
	hsm              infra.HSM

	keycloakClient *sync.KeycloakHTTPClient
}

func NewUserManager(repoProvider repos.IProvider, pubbcAdmin pubbc.Admin, hsm infra.HSM, keycloakClient *sync.KeycloakHTTPClient) IUserManager {
	return &UserManager{
		ticketRepository: repoProvider.GetTicketRepository(),
		eventRepository:  repoProvider.GetEventRepository(),
		userRepository:   repoProvider.GetUserRepository(),
		pubbcAdmin:       pubbcAdmin,
		hsm:              hsm,
		keycloakClient:   keycloakClient,
	}
}

func (userManager *UserManager) RegisterUser(email, password, firstname, lastname, providedPK string) (*models.Attendant, error) {
	attendantID := uuid.New()

	// todo -> this is to handle standalone executions
	// how we can do this more clean and beautiful?
	// I know this is not the best way to do this, but everybody
	// knows that life is difficult and led us to some difficult
	// decisions but i want to be engineer next month :) so time
	// is something priority right now
	if userManager.keycloakClient != nil {
		keycloakUser, err := userManager.keycloakClient.RegisterUser(
			firstname,
			lastname,
			password,
			email,
		)
		if err != nil {
			return nil, tickenerr.FromError(usererr.RegisterAttendantErrorCode, err)
		}
		attendantID = keycloakUser.ID
	}

	return userManager.CreateUser(attendantID, providedPK)
}

func (userManager *UserManager) CreateUser(attendantID uuid.UUID, providedPK string) (*models.Attendant, error) {
	newAttendant := &models.Attendant{UUID: attendantID}

	var pkStoreKey, walletAddr string
	var err error

	// check if user exists
	user := userManager.userRepository.FindUser(attendantID)
	if user != nil {
		return nil, tickenerr.New(usererr.UserAlreadyExistErrorCode)
	}

	if len(providedPK) > 0 {
		walletAddr, err = userManager.pubbcAdmin.GetWalletForKey(providedPK)
		if err != nil {
			return nil, tickenerr.FromError(usererr.CreateWallerErrorCode, err)
		}
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

func (userManager *UserManager) GetUser(uuid uuid.UUID) (*models.Attendant, error) {
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
