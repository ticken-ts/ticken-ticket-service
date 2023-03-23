package models

import (
	"fmt"
	"github.com/google/uuid"
	"math/big"
)

type Ticket struct {
	/*************** ticket key ****************/
	EventID  uuid.UUID `bson:"event_id"`
	TicketID uuid.UUID `bson:"ticket_id"`
	/*******************************************/

	/****************** owner ******************/
	OwnerID uuid.UUID `bson:"owner"`
	/*******************************************/

	/****************** info *******************/
	Section string `bson:"section"`
	Status  string `bson:"status"`
	/*******************************************/

	/************** blockchain *****************/
	TokenID   *big.Int `bson:"token_id"`
	PubbcTxID string   `bson:"pubbc_tx_id"`
	PvtbcTxID string   `bson:"pvtbc_tx_id"`
	/*******************************************/
}

func (ticket *Ticket) TransferTo(anotherAttendant *User) error {
	if ticket.IsOwnedBy(anotherAttendant) {
		return fmt.Errorf("ticket is already ownerd by %s", anotherAttendant.UUID)
	}

	ticket.OwnerID = anotherAttendant.UUID

	return nil
}

func (ticket *Ticket) IsOwnedBy(attendant *User) bool {
	return attendant.UUID == ticket.OwnerID
}
