package services

import (
	"ticken-ticket-service/infra/public_blockchain"
	"ticken-ticket-service/models"
	"ticken-ticket-service/repos"
)

type eventManager struct {
	eventRepository  repos.EventRepository
	ticketRepository repos.TicketRepository
	blockchain       *public_blockchain.PublicBlockchain
}

func NewEventManager(
	eventRepository repos.EventRepository,
	ticketRepository repos.TicketRepository,
	blockchain *public_blockchain.PublicBlockchain,
) EventManager {
	return &eventManager{
		ticketRepository: ticketRepository,
		eventRepository:  eventRepository,
		blockchain:       blockchain,
	}
}

// ListenBlockchainEvents Listen event ticket created from all events on database
func (eventManager *eventManager) ListenBlockchainEvents() error {
	events, err := eventManager.eventRepository.GetActiveEvents()
	if err != nil {
		return err
	}

	for _, event := range events {
		instance, err := eventManager.blockchain.GetContract(event.PubBcAddress)
		if err != nil {
			return err
		}

		err = instance.WatchTicketCreatedEvent(eventManager.getNewTicketHandler(event.EventID))
		if err != nil {
			return err
		}
	}

	return nil
}

func (eventManager *eventManager) AddEvent(EventID string, OrganizerID string, PvtBCChannel string) (*models.Event, error) {
	addr, err := eventManager.blockchain.DeployContract()
	if err != nil {
		return nil, err
	}

	instance, err := eventManager.blockchain.GetContract(addr)
	if err != nil {
		return nil, err
	}

	err = instance.WatchTicketCreatedEvent(eventManager.getNewTicketHandler(EventID))
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

func (eventManager *eventManager) getNewTicketHandler(eventID string) func(ticket *public_blockchain.CreatedTicket) {
	return func(ticket *public_blockchain.CreatedTicket) {
		// TODO: create ticket in our database
	}
}
