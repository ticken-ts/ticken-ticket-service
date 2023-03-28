package models

import (
	"fmt"
	"github.com/google/uuid"
	"math/big"
	"ticken-ticket-service/utils/money"
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

	/***************** sales *******************/
	Resells []*Resell `bson:"resells"`
	/*******************************************/
}

func (ticket *Ticket) IsOnSale() bool {
	if ticket.Resells == nil || len(ticket.Resells) == 0 {
		return false
	}
	for _, announcement := range ticket.Resells {
		if announcement.Active {
			return true
		}
	}
	return false
}

func (ticket *Ticket) TransferTo(anotherAttendant *User) error {
	if ticket.IsOwnedBy(anotherAttendant) {
		return fmt.Errorf(
			"ticket is already ownerd by %s",
			anotherAttendant.UUID,
		)
	}

	ticket.OwnerID = anotherAttendant.UUID
	return nil
}

func (ticket *Ticket) IsOwnedBy(attendant *User) bool {
	return attendant.UUID == ticket.OwnerID
}

func (ticket *Ticket) CreateResell(price *money.Money) (*Resell, error) {
	if ticket.Resells == nil {
		ticket.Resells = make([]*Resell, 0)
	}

	for _, announcement := range ticket.Resells {
		if announcement.Active && announcement.Price.Currency == price.Currency {
			return nil,
				fmt.Errorf(
					"there is already a sale announcement for this ticket for the currency %s",
					price.Currency.Name,
				)
		}
	}

	newResell := &Resell{Price: price, Active: true, ResellID: uuid.New()}
	ticket.Resells = append(ticket.Resells, newResell)

	return newResell, nil
}

func (ticket *Ticket) GetResell(resellID uuid.UUID) *Resell {
	for _, resell := range ticket.Resells {
		if resell.ResellID == resellID {
			return resell
		}
	}
	return nil
}

func (ticket *Ticket) SellTo(buyer *User, resellID uuid.UUID) error {
	if !ticket.IsOnSale() {
		return fmt.Errorf("ticket is not on sale")
	}

	if ticket.GetResell(resellID) == nil {
		return fmt.Errorf("resell %s doest not exist", resellID.String())
	}

	// invalidate all resells
	for _, resell := range ticket.Resells {
		resell.Deactivate()
	}

	if err := ticket.TransferTo(buyer); err != nil {
		return err
	}

	return nil
}

func (ticket *Ticket) ToBatman() {
	ticket.OwnerID = uuid.Nil
}
