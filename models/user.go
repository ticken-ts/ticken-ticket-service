package models

import (
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	mongoID           primitive.ObjectID `bson:"_id"`
	UUID              uuid.UUID          `bson:"uuid"`
	AddressPKStoreKey string             `bson:"addressPKStore"`
	WalletAddress     string             `bson:"wallet_address"`
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
