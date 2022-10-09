package services

import (
	"fmt"
	pvtbc "github.com/ticken-ts/ticken-pvtbc-connector"
	"ticken-ticket-service/models"
	"ticken-ticket-service/repos"
)

type ticketSigner struct {
	eventRepository  repos.EventRepository
	ticketRepository repos.TicketRepository
	pvtbcConnector   *pvtbc.Caller
	userManager      *UserManager
}

func NewTicketSigner(
	eventRepository repos.EventRepository,
	ticketRepository repos.TicketRepository,
	pvtbcConnector *pvtbc.Caller,
	userManager *UserManager,
) TicketSigner {
	return &ticketSigner{
		eventRepository:  eventRepository,
		ticketRepository: ticketRepository,
		pvtbcConnector:   pvtbcConnector,
		userManager:      userManager,
	}
}

func (s *ticketSigner) SignTicket(eventID string, ticketID string, signer string) (*models.Ticket, error) {
	event := s.eventRepository.FindEvent(eventID)
	if event == nil {
		return nil, fmt.Errorf("could not determine organizer channel")
	}

	err := s.pvtbcConnector.SetChannel(event.PvtBCChannel)
	if err != nil {
		return nil, err
	}

	privateKey, err := s.userManager.GetUserPrivateKey(signer)
	if err != nil {
		return nil, err
	}

	ticket := s.ticketRepository.FindTicket(eventID, ticketID)
	if ticket == nil {
		// TODO - handle this situation with a sync with
		// pvt blockchain to ensure that the ticket exist or not
		return nil, fmt.Errorf("ticket %s not found in event %s", ticketID, eventID)
	}

	signature, err := ticket.Sign(privateKey)
	if err != nil {
		return nil, err
	}

	ticketResponse, err := s.pvtbcConnector.SignTicket(ticketID, eventID, signer, signature)
	if err != nil {
		return nil, err
	}

	ticket.Status = ticketResponse.Status
	err = s.ticketRepository.UpdateTicketStatus(ticket)
	if err != nil {
		return nil, err
	}

	return ticket, nil
}