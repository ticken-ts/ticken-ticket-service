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
	event := s.eventRepository.FindEvent(eventID)
	if event == nil {
		return nil, fmt.Errorf("could not determine organizer channel")
	}

	err := s.pvtbcConnector.Connect(event.PvtBCChannel)
	if err != nil {
		return nil, err
	}

	newTicket := models.NewTicket(eventID, section, owner)

	ticketResponse, err := s.pvtbcConnector.IssueTicket(
		newTicket.TicketID,
		newTicket.EventID,
		newTicket.Section,
		newTicket.Owner,
	)

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
