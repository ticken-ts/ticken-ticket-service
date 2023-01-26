package services

import (
	"fmt"
	"github.com/google/uuid"
	pvtbc "github.com/ticken-ts/ticken-pvtbc-connector"
	"ticken-ticket-service/infra/public_blockchain"
	"ticken-ticket-service/models"
	"ticken-ticket-service/repos"
)

type ticketIssuer struct {
	eventRepository  repos.EventRepository
	ticketRepository repos.TicketRepository
	pvtbcConnector   *pvtbc.Caller
	pubBCConnector   public_blockchain.PublicBC
}

func NewTicketIssuer(
	eventRepository repos.EventRepository,
	ticketRepository repos.TicketRepository,
	pvtbcConnector *pvtbc.Caller,
	blockchain public_blockchain.PublicBC,
) TicketIssuer {
	return &ticketIssuer{
		eventRepository:  eventRepository,
		ticketRepository: ticketRepository,
		pvtbcConnector:   pvtbcConnector,
		pubBCConnector:   blockchain,
	}
}

func (s *ticketIssuer) IssueTicket(eventID uuid.UUID, section string, ownerID uuid.UUID) (*models.Ticket, error) {
	event := s.eventRepository.FindEvent(eventID.String())
	if event == nil {
		return nil, fmt.Errorf("could not determine organizer channel")
	}

	err := s.pvtbcConnector.SetChannel(event.PvtBCChannel)
	if err != nil {
		return nil, err
	}

	newTicket := models.NewTicket(eventID, section, ownerID)

	ticketResponse, err := s.pvtbcConnector.IssueTicket(
		newTicket.TicketID, newTicket.EventID, newTicket.OwnerID, newTicket.Section,
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
