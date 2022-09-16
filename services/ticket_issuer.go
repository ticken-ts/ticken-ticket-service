package services

import (
	"fmt"
	"ticken-ticket-service/blockchain/pvtbc"
	"ticken-ticket-service/models"
	"ticken-ticket-service/repos"
)

type ticketIssuer struct {
	eventRepository  repos.EventRepository
	ticketRepository repos.TicketRepository
	pvtbcConnector   pvtbc.TickenConnector
}

func NewTicketIssuer(
	eventRepository repos.EventRepository,
	ticketRepository repos.TicketRepository,
	pvtbcConnector pvtbc.TickenConnector,
) TicketIssuer {
	return &ticketIssuer{
		eventRepository:  eventRepository,
		ticketRepository: ticketRepository,
		pvtbcConnector:   pvtbcConnector,
	}
}

func (s *ticketIssuer) IssueTicket(eventID string, section string, owner string) (*models.Ticket, error) {
	event := s.eventRepository.FindEventByID(eventID)
	if event == nil {
		return nil, fmt.Errorf("could not determine organizer channel")
	}

	fmt.Printf("channel found %s", event.PvtBCChannel)

	err := s.pvtbcConnector.Connect(event.PvtBCChannel)
	if err != nil {
		return nil, err
	}

	newTicket := models.NewTicket(eventID, section)
	err = newTicket.AssignTo(owner)
	if err != nil {
		return nil, err
	}

	err = s.pvtbcConnector.IssueTicket(newTicket)
	if err != nil {
		return nil, err
	}

	return newTicket, nil
}
