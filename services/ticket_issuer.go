package services

import (
	"fmt"
	pvtbc "github.com/ticken-ts/ticken-pvtbc-connector"
	"ticken-ticket-service/infra/public_blockchain"
	"ticken-ticket-service/models"
	"ticken-ticket-service/repos"
)

type ticketIssuer struct {
	eventRepository  repos.EventRepository
	ticketRepository repos.TicketRepository
	pvtbcConnector   *pvtbc.Caller
	pubBCConnector   *public_blockchain.PublicBlockchain
}

func NewTicketIssuer(
	eventRepository repos.EventRepository,
	ticketRepository repos.TicketRepository,
	pvtbcConnector *pvtbc.Caller,
	blockchain *public_blockchain.PublicBlockchain,
) TicketIssuer {
	return &ticketIssuer{
		eventRepository:  eventRepository,
		ticketRepository: ticketRepository,
		pvtbcConnector:   pvtbcConnector,
		pubBCConnector:   blockchain,
	}
}

func (s *ticketIssuer) IssueTicket(eventID string, section string, owner string) (*models.Ticket, error) {
	event := s.eventRepository.FindEvent(eventID)
	if event == nil {
		return nil, fmt.Errorf("could not determine organizer channel")
	}

	err := s.pvtbcConnector.SetChannel(event.PvtBCChannel)
	if err != nil {
		return nil, err
	}

	newTicket := models.NewTicket(eventID, section, owner)

	ticketResponse, err := s.pvtbcConnector.IssueTicket(newTicket.TicketID, newTicket.EventID, newTicket.Section, newTicket.Owner)
	if err != nil {
		return nil, err
	}

	newTicket.Status = ticketResponse.Status

	err = s.ticketRepository.AddTicket(newTicket)
	if err != nil {
		return nil, err
	}

	return newTicket, nil
}
