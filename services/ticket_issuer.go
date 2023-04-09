package services

import (
	"fmt"
	"github.com/google/uuid"
	pubbc "github.com/ticken-ts/ticken-pubbc-connector"
	pvtbc "github.com/ticken-ts/ticken-pvtbc-connector"
	"math/big"
	"ticken-ticket-service/infra"
	"ticken-ticket-service/log"
	"ticken-ticket-service/models"
	"ticken-ticket-service/repos"
	"ticken-ticket-service/tickenerr"
	"ticken-ticket-service/tickenerr/commonerr"
	"ticken-ticket-service/tickenerr/eventerr"
	"ticken-ticket-service/tickenerr/ticketerr"
	"ticken-ticket-service/tickenerr/usererr"
)

type TicketIssuer struct {
	eventRepository  repos.EventRepository
	ticketRepository repos.TicketRepository
	userRepository   repos.UserRepository
	hsm              infra.HSM
	pvtbcCaller      *pvtbc.Caller
	pubbcCaller      pubbc.Caller
}

func NewTicketIssuer(
	repoProvider repos.IProvider,
	hsm infra.HSM,
	pubbcCaller pubbc.Caller,
	pvtbcCaller *pvtbc.Caller,
) ITicketIssuer {
	return &TicketIssuer{
		ticketRepository: repoProvider.GetTicketRepository(),
		eventRepository:  repoProvider.GetEventRepository(),
		userRepository:   repoProvider.GetUserRepository(),
		hsm:              hsm,
		pubbcCaller:      pubbcCaller,
		pvtbcCaller:      pvtbcCaller,
	}
}

func (s *TicketIssuer) IssueTicket(eventID uuid.UUID, section string, attendantID uuid.UUID) (*models.Ticket, error) {
	event := s.eventRepository.FindEvent(eventID)
	if event == nil {
		return nil, tickenerr.New(eventerr.EventNotFoundErrorCode)
	}

	attendant := s.userRepository.FindUser(attendantID)
	if attendant == nil {
		return nil, tickenerr.New(usererr.UserNotFoundErrorCode)
	}

	err := s.pvtbcCaller.SetChannel(event.PvtBCChannel)
	if err != nil {
		return nil, tickenerr.FromError(
			commonerr.FailedToEstablishConnectionWithPVTBCErrorCode, err)
	}

	tickenID := uuid.New()
	tokenID, err := generateTokenID(tickenID)
	if err != nil {
		return nil, tickenerr.FromError(
			ticketerr.FailedToGenerateTokenIDErrorCode,
			err)
	}

	newTicket := &models.Ticket{
		/*************** ticket key ****************/
		TicketID: tickenID,
		EventID:  eventID,
		/*******************************************/

		/****************** owner ******************/
		OwnerID: attendantID,
		/*******************************************/

		/****************** info *******************/
		Section: section,
		Status:  models.TicketStatusIssued,
		/*******************************************/

		/************** blockchain *****************/
		TokenID: tokenID,

		// this fields will be populated after each
		// blockchain transaction, indicating if the
		// ticket is synced with this blockchain
		PubbcTxID: "",
		PvtbcTxID: "",
		/*******************************************/
	}

	if err := s.ticketRepository.AddOne(newTicket); err != nil {
		return nil, err
	}

	// ******************* PVTBC Ticket Issuing ******************* //
	_, pvtbcTxID, err := s.pvtbcCaller.IssueTicket(
		newTicket.TicketID,
		newTicket.EventID,
		newTicket.OwnerID,
		newTicket.Section,
		newTicket.TokenID,
	)
	if err != nil {
		return nil, tickenerr.FromError(ticketerr.FailedToMintTicketInPVTBC, err)
	}
	newTicket.PvtbcTxID = pvtbcTxID
	if err := s.ticketRepository.UpdateTicketBlockchainData(newTicket); err != nil {
		return nil, err
	}
	// ************************************************************ //

	// ******************* PUBBC Ticket Issuing ******************* //
	pubbcTxID, err := s.pubbcCaller.MintTicket(
		event.PubBCAddress,
		attendant.WalletAddress,
		newTicket.Section,
		newTicket.TokenID,
	)
	if err != nil {
		return nil, tickenerr.FromError(ticketerr.FailedToMintTicketInPUBBC, err)
	}
	newTicket.PubbcTxID = pubbcTxID
	if err := s.ticketRepository.UpdateTicketBlockchainData(newTicket); err != nil {
		return nil, err
	}
	// ************************************************************ //

	return newTicket, nil
}

func (s *TicketIssuer) GetUserTickets(attendantID uuid.UUID) ([]*models.Ticket, error) {
	attendant := s.userRepository.FindUser(attendantID)
	if attendant == nil {
		return nil, tickenerr.New(usererr.UserNotFoundErrorCode)
	}

	userTickets, err := s.ticketRepository.GetUserTickets(attendantID)
	if err != nil {
		return nil, tickenerr.FromError(ticketerr.ReadUserTicketsFromDatabaseErrorCode, err)
	}

	// represents the tickets that the user still has
	// after syncing it with the public blockchain
	filteredTickets := make([]*models.Ticket, 0)

	for _, ticket := range userTickets {
		event := s.eventRepository.FindEvent(ticket.EventID)
		if event == nil {
			log.TickenLogger.Warn().Msg(fmt.Sprintf(
				"Couldn't find event for ticket %s", ticket.TicketID.String(),
			))
			continue
		}

		pubbcTicket, err := s.pubbcCaller.GetTicket(event.PubBCAddress, ticket.TokenID)
		if err != nil {
			log.TickenLogger.Warn().Msg(fmt.Sprintf(
				"Failed to get ticket %s from contract with addr %s: %s", ticket.TicketID.String(), event.PubBCAddress,
				err.Error(),
			))
			continue
		}

		if pubbcTicket.TokenID.Text(16) != ticket.TokenID.Text(16) {
			log.TickenLogger.Panic().Msg("token ID differs") // should never happen
		}

		if pubbcTicket.OwnerAddr != attendant.WalletAddress {
			ticket.ToBatman()
			_ = s.ticketRepository.UpdateTicketOwner(ticket)
		} else {
			filteredTickets = append(filteredTickets, ticket)
		}
	}
	return filteredTickets, nil
}

// GenerateTokenID generates a token ID in uint256
// for the ticket by hashing the ticket ID.
func generateTokenID(ticketID uuid.UUID) (*big.Int, error) {
	bytes, err := ticketID.MarshalBinary()
	if err != nil {
		return nil, err
	}
	tokenID := big.NewInt(0).SetBytes(bytes)
	return tokenID, nil
}
