package services

import (
	"fmt"
	"github.com/google/uuid"
	pubbc "github.com/ticken-ts/ticken-pubbc-connector"
	"ticken-ticket-service/env"
	"ticken-ticket-service/log"
	"ticken-ticket-service/models"
	"ticken-ticket-service/repos"
	"ticken-ticket-service/tickenerr"
	"ticken-ticket-service/tickenerr/commonerr"
	"ticken-ticket-service/tickenerr/eventerr"
	"ticken-ticket-service/tickenerr/ticketerr"
	"ticken-ticket-service/tickenerr/usererr"
)

type ticketLinker struct {
	eventRepository  repos.EventRepository
	ticketRepository repos.TicketRepository
	userRepository   repos.UserRepository
	pubbcCaller      pubbc.Caller
}

func NewTicketLinker(repoProvider repos.IProvider, pubbcCaller pubbc.Caller) ITicketLinker {
	return &ticketLinker{
		eventRepository:  repoProvider.GetEventRepository(),
		ticketRepository: repoProvider.GetTicketRepository(),
		userRepository:   repoProvider.GetUserRepository(),
		pubbcCaller:      pubbcCaller,
	}
}

func (t ticketLinker) LinkTickets(attendantID uuid.UUID, eventContractAddress string) ([]*models.Ticket, error) {
	attendant := t.userRepository.FindUser(attendantID)
	if attendant == nil {
		return nil, tickenerr.New(usererr.UserNotFoundErrorCode)
	}

	event := t.eventRepository.FindEventByContractAddress(eventContractAddress)
	if event == nil {
		return nil, tickenerr.NewWithMessage(
			eventerr.EventNotFoundErrorCode,
			fmt.Sprintf("contract addr: %s", eventContractAddress))
	}

	pubbcTickets, err := t.pubbcCaller.GetTickets(eventContractAddress, attendant.WalletAddress)
	if err != nil {
		return nil, tickenerr.FromError(ticketerr.FailedToRetrievePUBBCTicketsErrorCode, err)
	}

	newTickets := make([]*models.Ticket, 0)
	for _, pubbcTicket := range pubbcTickets {
		ticket := t.ticketRepository.FindTicketByPUBBCToken(event.EventID, pubbcTicket.TokenID)

		// this case is when ticket was never in the database.
		// this should never happen, because the user cant mint
		// the tickets by themselves.
		if ticket == nil {
			msg := fmt.Sprintf("ticket with token ID %s present in PUBBC not found in local database: %v",
				ticket.TokenID.Text(16), err)

			if env.TickenEnv.IsDev() {
				panic(msg)
			}
			log.TickenLogger.Warn().Msg(msg)
		}

		// if the user didn't change, we continue
		if ticket.IsOwnedBy(attendant) {
			continue
		}

		if err := ticket.TransferTo(attendant); err == nil {
			if err := t.ticketRepository.UpdateTicketOwner(ticket); err != nil {
				return nil, tickenerr.FromErrorWithMessage(
					commonerr.FailedToUpdateElement, err,
					fmt.Sprintf("owner of ticket with token ID %s", ticket.TokenID.Text(16)))
			}
		} else { // this should never happen
			msg := fmt.Sprintf("error transferring ticket locally from %s to %s: %v",
				ticket.OwnerID.String(), attendant.UUID.String(), err)

			if env.TickenEnv.IsDev() {
				panic(msg)
			}
			log.TickenLogger.Warn().Msg(msg)
		}

		newTickets = append(newTickets, ticket)
	}

	return newTickets, nil
}
