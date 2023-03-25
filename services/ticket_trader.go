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
	attendant := ticketTrader.userRepository.FindUser(ownerID)
	if attendant == nil {
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
	if !ticket.IsOwnedBy(attendant) {
		return nil, fmt.Errorf("%s is not the ticket owner", attendant.UUID)
	}

	newSaleAnnouncement, err := ticket.CreateSaleAnnouncement(price)
	if err != nil {
		return nil, err
	}

	err = ticketTrader.ticketRepository.AddTicketSaleAnnouncement(eventID, ticketID, newSaleAnnouncement)
	if err != nil {
		return nil, err
	}

	if newSaleAnnouncement.IsOnBlockchain() {
		panic("still not supported")
	}

	return ticket, nil
}
