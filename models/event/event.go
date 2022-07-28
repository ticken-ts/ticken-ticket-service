package event

type Event struct {
	EventID      string `json:"event_id" bson:"_id"`
	OrganizerID  string `json:"organizer_id" bson:"organizer_id"`
	PvtBCChannel string `json:"pvt_bc_channel" bson:"pvt_bc_channel"`
}
