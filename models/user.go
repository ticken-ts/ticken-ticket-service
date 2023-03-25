package models

import (
	"github.com/google/uuid"
)

type User struct {
	UUID              uuid.UUID `bson:"uuid"`
	AddressPKStoreKey string    `bson:"addressPKStore"`
	WalletAddress     string    `bson:"wallet_address"`
}

func NewUser(id uuid.UUID) *User {
	return &User{
		UUID: id,
	}
}

func (u *User) SetWallet(privKeyStorageKey, address string) {
	u.AddressPKStoreKey = privKeyStorageKey
	u.WalletAddress = address
}
