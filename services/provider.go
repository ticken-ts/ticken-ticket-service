package services

import (
	pubbc "github.com/ticken-ts/ticken-pubbc-connector"
	pvtbc "github.com/ticken-ts/ticken-pvtbc-connector"
	"ticken-ticket-service/infra"
	"ticken-ticket-service/repos"
)

type Provider struct {
	ticketIssuer ITicketIssuer
	ticketLinker ITicketLinker
	eventManager IEventManager
	userManager  IUserManager
	ticketTrader ITicketTrader
}

func NewProvider(
	repoProvider repos.IProvider,
	pvtbcCaller *pvtbc.Caller,
	pubbcAdmin pubbc.Admin,
	pubbcCaller pubbc.Caller,
	hsm infra.HSM,
) (*Provider, error) {
	provider := new(Provider)

	provider.eventManager = NewEventManager(repoProvider)
	provider.ticketLinker = NewTicketLinker(repoProvider, pubbcCaller)
	provider.ticketTrader = NewTicketTrader(repoProvider, pubbcCaller)
	provider.userManager = NewUserManager(repoProvider, pubbcAdmin, hsm)
	provider.ticketIssuer = NewTicketIssuer(repoProvider, hsm, pubbcCaller, pvtbcCaller)

	return provider, nil
}

func (provider *Provider) GetTicketIssuer() ITicketIssuer {
	return provider.ticketIssuer
}

func (provider *Provider) GetEventManager() IEventManager {
	return provider.eventManager
}

func (provider *Provider) GetUserManager() IUserManager {
	return provider.userManager
}

func (provider *Provider) GetTicketLinker() ITicketLinker {
	return provider.ticketLinker
}

func (provider *Provider) GetTicketTrader() ITicketTrader {
	return provider.ticketTrader
}
