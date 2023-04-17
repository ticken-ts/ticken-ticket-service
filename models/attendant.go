package models

import (
	"github.com/google/uuid"
)

type Attendant struct {
	UUID          uuid.UUID `bson:"uuid"`
	WalletAddress string    `bson:"wallet_address"`

	PrivStoreKey string `bson:"priv_store_key"`
	PubKey       string `bson:"pub_key"`
}

func (u *Attendant) SetWallet(privStoreKey, pubKey, walletAddr string) {
	u.PrivStoreKey = privStoreKey
	u.PubKey = pubKey
	u.WalletAddress = walletAddr
}
