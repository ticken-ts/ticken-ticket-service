package services

import (
	"ticken-ticket-service/infra/public_blockchain"
	"ticken-ticket-service/models"
	"ticken-ticket-service/repos"
)

type eventManager struct {
	eventRepository  repos.EventRepository
	ticketRepository repos.TicketRepository
	blockchain       public_blockchain.PublicBC
}

func NewEventManager(
	eventRepository repos.EventRepository,
	ticketRepository repos.TicketRepository,
	blockchain public_blockchain.PublicBC,
) EventManager {
	return &eventManager{
		ticketRepository: ticketRepository,
		eventRepository:  eventRepository,
		blockchain:       blockchain,
	}
}

func (eventManager *eventManager) AddEvent(EventID string, OrganizerID string, PvtBCChannel string) (*models.Event, error) {
	addr, err := eventManager.blockchain.DeployContract()
	if err != nil {
		return nil, err
	}

	event := models.NewEvent(EventID, OrganizerID, PvtBCChannel, addr)
	err = eventManager.eventRepository.AddEvent(event)
	if err != nil {
		return nil, err
	}
	return event, err
}
