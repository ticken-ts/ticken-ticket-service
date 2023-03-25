package models

import (
	"github.com/google/uuid"
)

type Event struct {
	EventID      uuid.UUID `bson:"event_id"`
	OrganizerID  uuid.UUID `bson:"organizer_id"`
	PvtBCChannel string    `bson:"pvt_bc_channel"`
	PubBCAddress string    `bson:"pub_bc_address"`
}
