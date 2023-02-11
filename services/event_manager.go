package services

import (
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

func (eventManager *EventManager) AddEvent(eventID, organizerID, pvtBCChannel, pubBCAddress string) (*models.Event, error) {
	event := models.NewEvent(eventID, organizerID, pvtBCChannel, pubBCAddress)
	if err := eventManager.eventRepository.AddEvent(event); err != nil {
		return nil, err
	}
	return event, nil
}
