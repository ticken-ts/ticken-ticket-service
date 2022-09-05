package services

import (
	"fmt"
	"ticken-ticket-service/blockchain/tickenPVTBCConnector"
	"ticken-ticket-service/models"
	"ticken-ticket-service/repositories"
)

type ticketIssuer struct {
	eventRepository  repositories.EventRepository
	ticketRepository repositories.TicketRepository
	pvtbcConnector   tickenPVTBCConnector.TickenPVTBConnector
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

func (s *ticketIssuer) IssueTicket(eventID string, section string, owner string) (*models.Ticket, error) {
	event := s.eventRepository.FindEventByID(eventID)
	if event == nil {
		return nil, fmt.Errorf("could not determine organizer channel")
	}

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
