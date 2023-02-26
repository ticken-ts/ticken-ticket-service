package models

import (
	"github.com/google/uuid"
	"math/big"
	"ticken-ticket-service/utils"
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

type ticketSignatureFields struct {
	ticketID uuid.UUID
	eventID  uuid.UUID
}

func (ticket *Ticket) Sign(ownerPrivateKey string) ([]byte, error) {
	signerHelper, err := utils.NewSigner(ownerPrivateKey)
	if err != nil {
		return nil, err
	}

	signatureFields := &ticketSignatureFields{
		ticketID: ticket.TicketID,
		eventID:  ticket.EventID,
	}

	signature, err := signerHelper.Sign(signatureFields)
	if err != nil {
		return nil, err
	}

	return signature, nil
}
