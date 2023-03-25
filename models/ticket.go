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
	SaleAnnouncements []*SaleAnnouncement `bson:"sale_announcements"`
	/*******************************************/
}

type SaleAnnouncement struct {
	Price  *money.Money `bson:"money"`
	Active bool         `bson:"active"`
}

func (announcement *SaleAnnouncement) IsOnBlockchain() bool {
	return announcement.Price.IsCrypto()
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

func (ticket *Ticket) CreateSaleAnnouncement(price *money.Money) (*SaleAnnouncement, error) {
	if ticket.SaleAnnouncements == nil {
		ticket.SaleAnnouncements = make([]*SaleAnnouncement, 0)
	}

	for _, announcement := range ticket.SaleAnnouncements {
		if announcement.Active && announcement.Price.Currency == price.Currency {
			return nil,
				fmt.Errorf(
					"there is already a sale announcement for this ticket for the currency %s",
					price.Currency.Name,
				)
		}
	}

	newSaleAnnouncement := &SaleAnnouncement{Price: price, Active: true}
	ticket.SaleAnnouncements = append(ticket.SaleAnnouncements, newSaleAnnouncement)

	return newSaleAnnouncement, nil
}
