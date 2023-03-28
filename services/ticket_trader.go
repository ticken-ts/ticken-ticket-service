package services

import (
	"fmt"
	"github.com/google/uuid"
	pubbc "github.com/ticken-ts/ticken-pubbc-connector"
	"ticken-ticket-service/models"
	"ticken-ticket-service/repos"
	"ticken-ticket-service/tickenerr"
	"ticken-ticket-service/tickenerr/commonerr"
	"ticken-ticket-service/tickenerr/eventerr"
	"ticken-ticket-service/tickenerr/ticketerr"
	"ticken-ticket-service/tickenerr/usererr"
	"ticken-ticket-service/utils/money"
)

type ticketTrader struct {
	eventRepository  repos.EventRepository
	ticketRepository repos.TicketRepository
	userRepository   repos.UserRepository
	pubbcCaller      pubbc.Caller
}

func NewTicketTrader(repoProvider repos.IProvider, pubbcCaller pubbc.Caller) ITicketTrader {
	return &ticketTrader{
		eventRepository:  repoProvider.GetEventRepository(),
		ticketRepository: repoProvider.GetTicketRepository(),
		userRepository:   repoProvider.GetUserRepository(),
		pubbcCaller:      pubbcCaller,
	}
}

func (ticketTrader *ticketTrader) SellTicket(ownerID, eventID, ticketID uuid.UUID, price *money.Money) (*models.Ticket, error) {
	seller := ticketTrader.userRepository.FindUser(ownerID)
	if seller == nil {
		return nil, tickenerr.New(usererr.UserNotFoundErrorCode)
	}

	event := ticketTrader.eventRepository.FindEvent(eventID)
	if event == nil {
		return nil, tickenerr.New(eventerr.EventNotFoundErrorCode)
	}

	ticket := ticketTrader.ticketRepository.FindTicket(eventID, ticketID)
	if ticket == nil {
		return nil, tickenerr.New(ticketerr.TicketNotFoundErrorCode)
	}

	newResell, err := ticket.CreateResell(price)
	if err != nil {
		return nil, tickenerr.FromError(ticketerr.CreateResellErrorCode, err)
	}
	if newResell.IsOnBlockchain() {
		return nil, tickenerr.NewWithMessage(
			ticketerr.ResellCurrencyNotSupportedErrorCode,
			"you can the resell directly from the contract",
		)
	}

	if err := ticketTrader.ticketRepository.AddTicketResell(eventID, ticketID, newResell); err != nil {
		return nil, tickenerr.FromErrorWithMessage(
			commonerr.FailedToUpdateElement, err,
			fmt.Sprintf("could not add resell to ticket"))
	}

	return ticket, nil
}

func (ticketTrader *ticketTrader) BuyResoldTicket(buyerID, eventID, ticketID, resellID uuid.UUID) (*models.Ticket, error) {
	buyer := ticketTrader.userRepository.FindUser(buyerID)
	if buyer == nil {
		return nil, tickenerr.NewWithMessage(
			usererr.UserNotFoundErrorCode,
			"seller user not present in the database")
	}

	event := ticketTrader.eventRepository.FindEvent(eventID)
	if event == nil {
		return nil, tickenerr.New(eventerr.EventNotFoundErrorCode)
	}

	ticket := ticketTrader.ticketRepository.FindTicket(eventID, ticketID)
	if ticket == nil {
		return nil, tickenerr.New(ticketerr.TicketNotFoundErrorCode)
	}

	seller := ticketTrader.userRepository.FindUser(ticket.OwnerID)
	if seller == nil { // this never should happen
		return nil, tickenerr.NewWithMessage(
			usererr.UserNotFoundErrorCode,
			"ticket owner not present in the database")
	}

	if err := ticket.SellTo(buyer, resellID); err != nil {
		return nil, tickenerr.FromError(ticketerr.BuyResellErrorCode, err)
	}

	_, err := ticketTrader.pubbcCaller.TransferTicket(
		event.PubBCAddress,
		ticket.TokenID,
		seller.WalletAddress,
		buyer.WalletAddress,
	)
	if err != nil {
		return nil, tickenerr.FromError(ticketerr.TransferTicketInPUBBCErrorCode, err)
	}

	if err := ticketTrader.ticketRepository.UpdateResoldTicket(ticket); err != nil {
		return nil, tickenerr.FromErrorWithMessage(
			commonerr.FailedToUpdateElement, err,
			"could not update resold ticket information")
	}

	return ticket, nil
}

func (ticketTrader *ticketTrader) GetTicketsInResells(eventID uuid.UUID, section string) ([]*models.Ticket, error) {
	return ticketTrader.ticketRepository.GetTicketsInResell(eventID, section)
}
