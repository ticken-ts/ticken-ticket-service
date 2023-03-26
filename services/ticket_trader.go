package services

import (
	"fmt"
	"github.com/google/uuid"
	pubbc "github.com/ticken-ts/ticken-pubbc-connector"
	"ticken-ticket-service/models"
	"ticken-ticket-service/repos"
	"ticken-ticket-service/utils/money"
)

type ticketTrader struct {
	eventRepository  repos.EventRepository
	ticketRepository repos.TicketRepository
	userRepository   repos.UserRepository
	pubbcCaller      pubbc.Caller
}

func NewTicketTrader(repoProvider repos.IProvider, pubbcCaller pubbc.Caller) TicketTrader {
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
		return nil, fmt.Errorf("user not found")
	}

	event := ticketTrader.eventRepository.FindEvent(eventID)
	if event == nil {
		return nil, fmt.Errorf("event not found")
	}

	ticket := ticketTrader.ticketRepository.FindTicket(eventID, ticketID)
	if ticket == nil {
		return nil, fmt.Errorf("ticket not found")
	}

	// todo -> handle ownership when the ticket were transferred without ht eapp
	if !ticket.IsOwnedBy(seller) {
		return nil, fmt.Errorf("%s is not the ticket owner", seller.UUID)
	}

	newResell, err := ticket.CreateResell(price)
	if err != nil {
		return nil, err
	}

	if newResell.IsOnBlockchain() {
		panic("still not supported")
	}

	if err := ticketTrader.ticketRepository.AddTicketResell(eventID, ticketID, newResell); err != nil {
		return nil, err
	}

	return ticket, nil
}

func (ticketTrader *ticketTrader) BuyResoldTicket(buyerID, eventID, ticketID, resellID uuid.UUID) (*models.Ticket, error) {
	buyer := ticketTrader.userRepository.FindUser(buyerID)
	if buyer == nil {
		return nil, fmt.Errorf("user not found")
	}

	event := ticketTrader.eventRepository.FindEvent(eventID)
	if event == nil {
		return nil, fmt.Errorf("event not found")
	}

	ticket := ticketTrader.ticketRepository.FindTicket(eventID, ticketID)
	if ticket == nil {
		return nil, fmt.Errorf("ticket not found")
	}

	seller := ticketTrader.userRepository.FindUser(ticket.OwnerID)
	if seller == nil {
		return nil, fmt.Errorf("ticket owner not found")
	}

	if err := ticket.SellTo(buyer, resellID); err != nil {
		return nil, err
	}

	_, err := ticketTrader.pubbcCaller.TransferTicket(
		event.PubBCAddress,
		ticket.TokenID,
		seller.WalletAddress,
		buyer.WalletAddress,
	)
	if err != nil {
		return nil, err
	}

	if err := ticketTrader.ticketRepository.UpdateResoldTicket(ticket); err != nil {
		return nil, err
	}

	return ticket, nil
}

func (ticketTrader *ticketTrader) GetResells(eventID uuid.UUID, section string) ([]*models.Ticket, error) {
	return ticketTrader.ticketRepository.GetTicketsInResell(eventID, section)
}
