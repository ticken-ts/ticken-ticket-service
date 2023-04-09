package fakes

import (
	"github.com/google/uuid"
	"ticken-ticket-service/config"
	"ticken-ticket-service/env"
	"ticken-ticket-service/infra"
	"ticken-ticket-service/models"
	"ticken-ticket-service/repos"
)

type FakeUsersPopulator struct {
	HSM    infra.HSM
	Config config.DevUser
	Repo   repos.UserRepository
}

func (populator *FakeUsersPopulator) Populate() error {
	if !env.TickenEnv.IsDev() {
		return nil
	}

	userID := uuid.MustParse(populator.Config.UserID)

	user := populator.Repo.FindUser(userID)
	if user != nil {
		return nil
	}

	walletPrivateKey := populator.Config.WalletPrivateKey

	hsmKey, _ := populator.HSM.Store([]byte(walletPrivateKey))

	fakeUser := &models.User{
		UUID:              userID,
		WalletAddress:     populator.Config.WalletAddress,
		AddressPKStoreKey: hsmKey,
	}

	return populator.Repo.AddOne(fakeUser)
}
