package async

import (
	"encoding/json"
	"ticken-ticket-service/models"
	"ticken-ticket-service/repos"
)

const (
	NewEventMessageType = "new_event"
)

type eventDTO struct {
	EventID      string `json:"event_id"`
	OrganizerID  string `json:"organizer_id"`
	PvtBCChannel string `json:"pvt_bc_channel"`
}

type EventReceiver struct {
	eventRepo repos.EventRepository
}

func NewEventReceiver(eventRepo repos.EventRepository) *EventReceiver {
	return &EventReceiver{eventRepo: eventRepo}
}

func (processor *EventReceiver) NewEventHandler(rawEvent []byte) error {
	dto := new(eventDTO)

	err := json.Unmarshal(rawEvent, dto)
	if err != nil {
		return err
	}

	event := models.NewEvent(dto.EventID, dto.OrganizerID, dto.PvtBCChannel)

	err = processor.eventRepo.AddEvent(event)
	if err != nil {
		return err
	}

	return nil
}
