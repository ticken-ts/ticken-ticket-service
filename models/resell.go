package models

import (
	"github.com/google/uuid"
	"ticken-ticket-service/utils/money"
)

type Resell struct {
	ResellID uuid.UUID    `bson:"resell_id"`
	Price    *money.Money `bson:"money"`
	Active   bool         `bson:"active"`
}

func (resell *Resell) IsOnBlockchain() bool {
	return resell.Price.IsCrypto()
}

func (resell *Resell) Deactivate() {
	resell.Active = false
}
