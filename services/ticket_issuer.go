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
	event := s.eventRepository.FindEvent(eventID.String())
	if event == nil {
		return nil, fmt.Errorf("could not determine organizer channel")
	}

	err := s.pvtbcCaller.SetChannel(event.PvtBCChannel)
	if err != nil {
		return nil, err
	}

	newTicket := models.NewTicket(eventID, section, ownerID)

	tokenID, err := generateTokenID(newTicket.TicketID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token ID: %w", err)
	}

	newTicket.TokenID = tokenID

	ticketResponse, err := s.pvtbcCaller.IssueTicket(
		newTicket.TicketID, newTicket.EventID, newTicket.OwnerID, newTicket.Section, newTicket.TokenID,
	)
	if err != nil {
		return nil, err
	}

	user := s.userRepository.FindUser(ownerID)
	if user == nil {
		return nil, fmt.Errorf("could not find user")
	}

	txHash, err := s.pubbcCaller.MintTicket(event.PubBCAddress, user.WalletAddress, newTicket.Section)
	if err != nil {
		return nil, fmt.Errorf("could not generate ticket on public address")
	}

	newTicket.Status = ticketResponse.Status
	newTicket.TxHash = txHash
	err = s.ticketRepository.AddTicket(newTicket)
	if err != nil {
		return nil, err
	}

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
