package services

import (
	"fmt"
	"ticken-ticket-service/blockchain/pvtbc"
	"ticken-ticket-service/helpers"
	"ticken-ticket-service/models"
	"ticken-ticket-service/repos"
)

type baseTicket struct {
	TicketID string `json:"ticket_id"`
	EventID  string `json:"event_id"`
}

type ticketSigner struct {
	eventRepository  repos.EventRepository
	ticketRepository repos.TicketRepository
	pvtbcConnector   pvtbc.TickenConnector
	userManager      *UserManager
}

func NewTicketSigner(
	eventRepository repos.EventRepository,
	ticketRepository repos.TicketRepository,
	pvtbcConnector pvtbc.TickenConnector,
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

	err := s.pvtbcConnector.Connect(event.PvtBCChannel)
	if err != nil {
		return nil, err
	}

	privateKey, err := s.userManager.GetUserPrivateKey(signer)
	if err != nil {
		return nil, err
	}

	signerHelper, err := helpers.NewSigner(privateKey)
	if err != nil {
		return nil, err
	}

	// this ticket is only used to sign
	baseTicket := &baseTicket{TicketID: ticketID, EventID: eventID}
	signature, err := signerHelper.Sign(baseTicket)
	if err != nil {
		return nil, err
	}

	ticketResponse, err := s.pvtbcConnector.SignTicket(ticketID, eventID, signer, signature)
	if err != nil {
		return nil, err
	}

	err = s.ticketRepository.UpdateTicketStatus(ticketID, eventID, ticketResponse.Status)
	if err != nil {
		return nil, err
	}

	ticket := s.ticketRepository.FindTicket(eventID, ticketID)

	return ticket, nil
}
