package services

import (
	"fmt"
	"github.com/google/uuid"
	pubbc "github.com/ticken-ts/ticken-pubbc-connector"
	"ticken-ticket-service/models"
	"ticken-ticket-service/repos"
)

type ticketLinker struct {
	eventRepository  repos.EventRepository
	ticketRepository repos.TicketRepository
	userRepository   repos.UserRepository
	pubbcCaller      pubbc.Caller
}

func NewTicketLinker(repoProvider repos.IProvider, pubbcCaller pubbc.Caller) TicketLinker {
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
		return nil, fmt.Errorf("user not found")
	}

	event := t.eventRepository.FindEventByContractAddress(eventContractAddress)
	if event == nil {
		return nil, fmt.Errorf("event with contract address %s not found", eventContractAddress)
	}

	pubbcTickets, err := t.pubbcCaller.GetTickets(eventContractAddress, attendant.WalletAddress)
	if err != nil {
		return nil, fmt.Errorf("could not obtain tickets from the blockchain: %s", err.Error())
	}

	newTickets := make([]*models.Ticket, 0)
	for _, pubbcTicket := range pubbcTickets {
		ticket := t.ticketRepository.FindTicketByPUBBCToken(event.EventID, pubbcTicket.TokenID)

		// this case is when ticket was never in the database.
		// this should never happen, because the user cant mint
		// the tickets by themselves.
		if ticket == nil {
			panic("found a ticket in blockchain that is not in db") // todo just for dev
		}

		// if the user didn't change, we continue
		if ticket.IsOwnedBy(attendant) {
			continue
		}

		// todo -> handle better these failures
		if err := ticket.TransferTo(attendant); err != nil {
			return nil, fmt.Errorf("failed to tranfer ticket %s: %s", ticket.TokenID.Text(16), err.Error())
		}

		if err := t.ticketRepository.UpdateTicketOwner(ticket); err != nil {
			return nil, fmt.Errorf("fialed to update ticket %s owner: %s", ticket.TokenID.Text(16), err.Error())
		}

		newTickets = append(newTickets, ticket)
	}

	return newTickets, nil
}
