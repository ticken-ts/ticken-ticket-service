package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Event struct {
	mongoID      primitive.ObjectID `bson:"_id"`
	EventID      string             `json:"event_id" bson:"event_id"`
	OrganizerID  string             `json:"organizer_id" bson:"organizer_id"`
	PvtBCChannel string             `json:"pvt_bc_channel" bson:"pvt_bc_channel"`
	PubBcAddress string             `json:"pub_bc_address" bson:"pub_bc_address"`
}

func NewEvent(EventID string, OrganizerID string, PvtBCChannel string, PubBcAddress string) *Event {
	return &Event{
		EventID:      EventID,
		OrganizerID:  OrganizerID,
		PvtBCChannel: PvtBCChannel,
		PubBcAddress: PubBcAddress,
	}
}
