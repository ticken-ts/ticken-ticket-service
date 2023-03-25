package services

import (
	pubbc "github.com/ticken-ts/ticken-pubbc-connector"
	pvtbc "github.com/ticken-ts/ticken-pvtbc-connector"
	"ticken-ticket-service/infra"
	"ticken-ticket-service/repos"
)

type Provider struct {
	ticketIssuer TicketIssuer
	ticketLinker TicketLinker
	eventManager IEventManager
	userManager  UserManager
	ticketTrader TicketTrader
}

func NewProvider(repoProvider repos.IProvider, pvtbcCaller *pvtbc.Caller, pubbcAdmin pubbc.Admin, pubbcCaller pubbc.Caller, hsm infra.HSM) (*Provider, error) {
	provider := new(Provider)

	eventRepo := repoProvider.GetEventRepository()
	ticketRepo := repoProvider.GetTicketRepository()
	userRepo := repoProvider.GetUserRepository()

	provider.eventManager = NewEventManager(eventRepo, ticketRepo)
	provider.userManager = NewUserManager(eventRepo, ticketRepo, userRepo, pubbcAdmin, hsm)
	provider.ticketIssuer = NewTicketIssuer(repoProvider, hsm, pubbcCaller, pvtbcCaller)
	provider.ticketLinker = NewTicketLinker(repoProvider, pubbcCaller)
	provider.ticketTrader = NewTicketTrader(repoProvider, pubbcCaller)

	return provider, nil
}

func (provider *Provider) GetTicketIssuer() TicketIssuer {
	return provider.ticketIssuer
}

func (provider *Provider) GetEventManager() IEventManager {
	return provider.eventManager
}

func (provider *Provider) GetUserManager() UserManager {
	return provider.userManager
}

func (provider *Provider) GetTicketLinker() TicketLinker {
	return provider.ticketLinker
}

func (provider *Provider) GetTicketTrader() TicketTrader {
	return provider.ticketTrader
}
