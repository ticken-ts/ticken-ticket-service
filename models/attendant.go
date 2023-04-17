package models

import (
	"github.com/google/uuid"
)

type Attendant struct {
	UUID              uuid.UUID `bson:"uuid"`
	AddressPKStoreKey string    `bson:"addressPKStore"`
	WalletAddress     string    `bson:"wallet_address"`
}

func (u *Attendant) SetWallet(privKeyStorageKey, address string) {
	u.AddressPKStoreKey = privKeyStorageKey
	u.WalletAddress = address
}
