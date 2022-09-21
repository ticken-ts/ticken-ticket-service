package models

import (
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Ticket struct {
	mongoID  primitive.ObjectID `bson:"_id"`
	TicketID string             `json:"ticket_id" bson:"ticket_id"`
	Owner    string             `json:"owner" bson:"owner"`
	Section  string             `json:"section" bson:"section"`
	EventID  string             `json:"event_id" bson:"event_id"`
	Status   string             `json:"status" bson:"status"`
}

func NewTicket(eventID string, section string, owner string) *Ticket {
	return &Ticket{
		TicketID: uuid.NewString(),
		EventID:  eventID,
		Section:  section,
		Owner:    owner,
	}
}
