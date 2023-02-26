package models

import (
	"github.com/google/uuid"
	"math/big"
)

type Ticket struct {
	TicketID uuid.UUID `bson:"ticket_id"`
	TokenID  *big.Int  `bson:"token_id"`
	OwnerID  uuid.UUID `bson:"owner"`
	Section  string    `bson:"section"`
	EventID  uuid.UUID `bson:"event_id"`
	Status   string    `bson:"status"`

	/************ blockchain **************/
	PubbcTxID string `bson:"pubbc_tx_id"`
	PvtbcTxID string `bson:"pvtbc_tx_id"`
	/**************************************/
}
