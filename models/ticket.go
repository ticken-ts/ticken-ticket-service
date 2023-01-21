package models

import (
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"math/big"
	"ticken-ticket-service/utils"
)

type Ticket struct {
	mongoID  primitive.ObjectID `bson:"_id"`
	TicketID string             `json:"ticket_id" bson:"ticket_id"`
	TokenID  big.Int            `json:"token_id" bson:"token_id"`
	Owner    string             `json:"owner" bson:"owner"`
	Section  string             `json:"section" bson:"section"`
	EventID  string             `json:"event_id" bson:"event_id"`
	Status   string             `json:"status" bson:"status"`
}

type ticketSignatureFields struct {
	ticketID string
	eventID  string
}

func NewTicket(eventID string, section string, owner string) *Ticket {
	return &Ticket{
		TicketID: uuid.NewString(),
		EventID:  eventID,
		Section:  section,
		Owner:    owner,
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
