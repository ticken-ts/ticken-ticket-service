package services

import (
	"ticken-ticket-service/models"
	"ticken-ticket-service/repos"
)

type eventManager struct {
	eventRepository repos.EventRepository
}

func NewEventManager(eventRepository repos.EventRepository) EventManager {
	return &eventManager{
		eventRepository: eventRepository,
	}
}

func (eventManager *eventManager) AddEvent(EventID string, OrganizerID string, PvtBCChannel string) (*models.Event, error) {
	event := models.NewEvent(EventID, OrganizerID, PvtBCChannel)
	err := eventManager.eventRepository.AddEvent(event)
	if err != nil {
		return nil, err
	}
	return event, err
}
