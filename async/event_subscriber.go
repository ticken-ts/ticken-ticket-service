package async

import (
	"encoding/json"
	"github.com/google/uuid"
	"ticken-ticket-service/services"
)

const (
	NewEventMessageType = "new_event"
)

type eventDTO struct {
	EventID      uuid.UUID `json:"event_id"`
	OrganizerID  uuid.UUID `json:"organizer_id"`
	PvtBCChannel string    `json:"pvt_bc_channel"`
}

type EventSubscriber struct {
	eventManager services.EventManager
}

func NewEventSubscriber(eventManager services.EventManager) *EventSubscriber {
	return &EventSubscriber{eventManager: eventManager}
}

func (s *EventSubscriber) NewEventHandler(rawEvent []byte) error {
	dto := new(eventDTO)

	err := json.Unmarshal(rawEvent, dto)
	if err != nil {
		return err
	}

	_, err = s.eventManager.AddEvent(dto.EventID.String(), dto.OrganizerID.String(), dto.PvtBCChannel)
	if err != nil {
		return err
	}

	println("llegue")

	return nil
}
