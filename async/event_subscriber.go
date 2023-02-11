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
	PubBCAddress string    `json:"pub_bc_address"`
}

type EventSubscriber struct {
	eventManager services.IEventManager
}

func NewEventSubscriber(eventManager services.IEventManager) *EventSubscriber {
	return &EventSubscriber{eventManager: eventManager}
}

func (s *EventSubscriber) NewEventHandler(rawEvent []byte) error {
	dto := new(eventDTO)

	err := json.Unmarshal(rawEvent, dto)
	if err != nil {
		return err
	}

	_, err = s.eventManager.AddEvent(
		dto.EventID.String(), dto.OrganizerID.String(), dto.PvtBCChannel, dto.PubBCAddress,
	)
	if err != nil {
		return err
	}

	return nil
}
