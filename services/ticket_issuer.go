package services

import (
	"fmt"
	"github.com/google/uuid"
	pubbc "github.com/ticken-ts/ticken-pubbc-connector"
	pvtbc "github.com/ticken-ts/ticken-pvtbc-connector"
	"math/big"
	"ticken-ticket-service/infra"
	"ticken-ticket-service/models"
	"ticken-ticket-service/repos"
)

type ticketIssuer struct {
	eventRepository  repos.EventRepository
	ticketRepository repos.TicketRepository
	userRepository   repos.UserRepository
	hsm              infra.HSM
	pvtbcCaller      *pvtbc.Caller
	pubbcCaller      pubbc.Caller
}

func NewTicketIssuer(
	eventRepository repos.EventRepository,
	ticketRepository repos.TicketRepository,
	userRepository repos.UserRepository,
	hsm infra.HSM,
	pubbcCaller pubbc.Caller,
	pvtbcCaller *pvtbc.Caller,
) TicketIssuer {
	return &ticketIssuer{
		eventRepository:  eventRepository,
		ticketRepository: ticketRepository,
		userRepository:   userRepository,
		hsm:              hsm,
		pubbcCaller:      pubbcCaller,
		pvtbcCaller:      pvtbcCaller,
	}
}

func (s *ticketIssuer) IssueTicket(eventID uuid.UUID, section string, ownerID uuid.UUID) (*models.Ticket, error) {
	event := s.eventRepository.FindEvent(eventID)
	if event == nil {
		return nil, fmt.Errorf("could not determine organizer channel")
	}

	user := s.userRepository.FindUser(ownerID)
	if user == nil {
		return nil, fmt.Errorf("could not find user")
	}

	err := s.pvtbcCaller.SetChannel(event.PvtBCChannel)
	if err != nil {
		return nil, err
	}

	newTicket := &models.Ticket{
		TicketID: uuid.New(),
		EventID:  eventID,
		Section:  section,
		OwnerID:  ownerID,

		// this fields will be populated after each
		// blockchain transaction, indicating if the
		// ticket is synced with this blockchain
		PubbcTxID: "",
		PvtbcTxID: "",
	}

	tokenID, err := generateTokenID(newTicket.TicketID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token ID: %w", err)
	}
	newTicket.TokenID = tokenID

	if err := s.ticketRepository.AddTicket(newTicket); err != nil {
		return nil, err
	}

	// ******************* PVTBC Ticket Issuing ******************* //
	ticketResponse, pvtbcTxID, err := s.pvtbcCaller.IssueTicket(
		newTicket.TicketID,
		newTicket.EventID,
		newTicket.OwnerID,
		newTicket.Section,
		newTicket.TokenID,
	)
	if err != nil {
		return nil, err
	}
	newTicket.PvtbcTxID = pvtbcTxID
	newTicket.Status = ticketResponse.Status
	if err := s.ticketRepository.UpdateTicketStatus(newTicket); err != nil {
		return nil, err
	}
	// ************************************************************ //

	// ******************* PUBBC Ticket Issuing ******************* //
	pubbcTxID, err := s.pubbcCaller.MintTicket(
		event.PubBCAddress,
		user.WalletAddress,
		newTicket.Section,
		newTicket.TokenID,
	)
	if err != nil {
		return nil, fmt.Errorf("could not generate ticket on public blockchain")
	}
	newTicket.PubbcTxID = pubbcTxID
	if err := s.ticketRepository.UpdateTicketStatus(newTicket); err != nil {
		return nil, err
	}
	// ************************************************************ //

	return newTicket, nil
}

func (s *ticketIssuer) GetUserTickets(userID uuid.UUID) ([]*models.Ticket, error) {
	return s.ticketRepository.GetUserTickets(userID)
}

// GenerateTokenID generates a token ID in uint256 for the ticket by hashing the ticket ID.
func generateTokenID(ticketID uuid.UUID) (*big.Int, error) {
	bytes, err := ticketID.MarshalBinary()
	if err != nil {
		return nil, err
	}
	tokenID := big.NewInt(0).SetBytes(bytes)
	return tokenID, nil
}
