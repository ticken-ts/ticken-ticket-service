package async

import (
	"encoding/json"
	"ticken-ticket-service/models"
	"ticken-ticket-service/repos"
)

const (
	CreateEventMessageType = "create_event"
)

type eventDTO struct {
	EventID      string `json:"event_id"`
	OrganizerID  string `json:"organizer_id"`
	PvtBCChannel string `json:"pvt_bc_channel"`
}

type EventProcessor struct {
	eventRepo repos.EventRepository
}

func NewEventProcessor(eventRepo repos.EventRepository) *EventProcessor {
	return &EventProcessor{eventRepo: eventRepo}
}

func (processor *EventProcessor) CreateEvent(rawEvent []byte) error {
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
