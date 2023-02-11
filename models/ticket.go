package models

import (
	"github.com/google/uuid"
	"ticken-ticket-service/utils"
)

type Ticket struct {
	TicketID uuid.UUID `bson:"ticket_id"`
	TokenID  int       `bson:"token_id"`
	OwnerID  uuid.UUID `bson:"owner"`
	Section  string    `bson:"section"`
	EventID  uuid.UUID `bson:"event_id"`
	Status   string    `bson:"status"`
	TxHash   string    `bson:"tx_hash"`
}

type ticketSignatureFields struct {
	ticketID uuid.UUID
	eventID  uuid.UUID
}

func NewTicket(eventID uuid.UUID, section string, ownerID uuid.UUID) *Ticket {
	return &Ticket{
		TicketID: uuid.New(),
		EventID:  eventID,
		Section:  section,
		OwnerID:  ownerID,
	}
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
