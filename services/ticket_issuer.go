package services

import (
	"fmt"
	"ticken-ticket-service/blockchain/tickenPVTBCConnector"
	"ticken-ticket-service/db/repositories"
	"ticken-ticket-service/models/ticket"
)

type ticketIssuer struct {
	eventRepository  repositories.EventRepository
	ticketRepository repositories.TicketRepository
	pvtbcConnector   tickenPVTBCConnector.TickenPVTBConnector
}

type TicketIssuer interface {
	IssueTicket(eventID string, section string, owner string) (*ticket.Ticket, error)
}

func NewTicketIssuer(
	eventRepository repositories.EventRepository,
	ticketRepository repositories.TicketRepository,
	pvtbcConnector tickenPVTBCConnector.TickenPVTBConnector,
) TicketIssuer {
	return &ticketIssuer{
		eventRepository:  eventRepository,
		ticketRepository: ticketRepository,
		pvtbcConnector:   pvtbcConnector,
	}
}

func (s *ticketIssuer) IssueTicket(eventID string, section string, owner string) (*ticket.Ticket, error) {
	event := s.eventRepository.FindEventByID(eventID)
	if event == nil {
		return nil, fmt.Errorf("event not found")
	}

	err := s.pvtbcConnector.Connect(event.PvtBCChannel)
	if err != nil {
		return nil, err
	}

	newTicket := ticket.New(eventID, section)
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
