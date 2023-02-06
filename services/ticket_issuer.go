package services

import (
	"fmt"
	"github.com/google/uuid"
	pvtbc "github.com/ticken-ts/ticken-pvtbc-connector"
	"ticken-ticket-service/infra"
	"ticken-ticket-service/infra/public_blockchain"
	"ticken-ticket-service/models"
	"ticken-ticket-service/repos"
)

type ticketIssuer struct {
	eventRepository  repos.EventRepository
	ticketRepository repos.TicketRepository
	userRepository   repos.UserRepository
	hsm              infra.HSM
	pvtbcConnector   *pvtbc.Caller
	pubBCConnector   public_blockchain.PublicBC
}

func NewTicketIssuer(
	eventRepository repos.EventRepository,
	ticketRepository repos.TicketRepository,
	userRepository repos.UserRepository,
	hsm infra.HSM,
	pvtbcConnector *pvtbc.Caller,
	blockchain public_blockchain.PublicBC,
) TicketIssuer {
	return &ticketIssuer{
		eventRepository:  eventRepository,
		ticketRepository: ticketRepository,
		userRepository:   userRepository,
		hsm:              hsm,
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

	// Get user address on public blockchain
	// ---------------------------------------------------------------------------
	user := s.userRepository.FindUser(ownerID)
	if user == nil {
		return nil, fmt.Errorf("could not find user")
	}

	res, err := s.hsm.Retrieve(user.AddressPKStoreKey)
	if err != nil {
		return nil, fmt.Errorf("could not retrieve user PK from HSM")
	}

	userPK := string(res)
	userAddress, err := s.pubBCConnector.GetAddressFromPK(userPK)
	if err != nil {
		return nil, fmt.Errorf("could not get user address from PK")
	}
	// ----------------------------------------------------------------------------

	contract, err := s.pubBCConnector.GetContract(event.PubBCAddress)
	if err != nil {
		return nil, fmt.Errorf("could not get contract")
	}

	_, err = contract.GenerateTicket(userAddress, newTicket)
	if err != nil {
		return nil, fmt.Errorf("could not generate ticket on public address")
	}

	newTicket.Status = ticketResponse.Status
	err = s.ticketRepository.AddTicket(newTicket)
	if err != nil {
		return nil, err
	}

	return newTicket, nil
}
