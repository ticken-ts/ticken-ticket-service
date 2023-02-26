package services

import (
	"github.com/google/uuid"
	"ticken-ticket-service/models"
	"ticken-ticket-service/repos"
)

type EventManager struct {
	eventRepository  repos.EventRepository
	ticketRepository repos.TicketRepository
}

func NewEventManager(eventRepository repos.EventRepository, ticketRepository repos.TicketRepository) *EventManager {
	return &EventManager{ticketRepository: ticketRepository, eventRepository: eventRepository}
}

func (eventManager *EventManager) AddEvent(eventID, organizerID uuid.UUID, pvtBCChannel, pubBCAddress string) (*models.Event, error) {
	event := &models.Event{
		EventID:      eventID,
		OrganizerID:  organizerID,
		PvtBCChannel: pvtBCChannel,
		PubBCAddress: pubBCAddress,
	}

	if err := eventManager.eventRepository.AddEvent(event); err != nil {
		return nil, err
	}
	return event, nil
}
