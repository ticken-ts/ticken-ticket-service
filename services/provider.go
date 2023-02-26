package services

// TODO
// * Check lazy build

import (
	pubbc "github.com/ticken-ts/ticken-pubbc-connector"
	pvtbc "github.com/ticken-ts/ticken-pvtbc-connector"
	"ticken-ticket-service/infra"
	"ticken-ticket-service/repos"
)

type Provider struct {
	ticketIssuer TicketIssuer
	ticketSigner TicketSigner
	eventManager IEventManager
	userManager  UserManager
}

func NewProvider(repoProvider repos.IProvider, pvtbcCaller *pvtbc.Caller, pubbcAdmin pubbc.Admin, pubbcCaller pubbc.Caller, hsm infra.HSM) (*Provider, error) {
	provider := new(Provider)

	eventRepo := repoProvider.GetEventRepository()
	ticketRepo := repoProvider.GetTicketRepository()
	userRepo := repoProvider.GetUserRepository()

	provider.eventManager = NewEventManager(eventRepo, ticketRepo)
	provider.userManager = NewUserManager(eventRepo, ticketRepo, userRepo, pubbcAdmin, hsm)
	provider.ticketIssuer = NewTicketIssuer(eventRepo, ticketRepo, userRepo, hsm, pubbcCaller, pvtbcCaller)

	return provider, nil
}

func (provider *Provider) GetTicketIssuer() TicketIssuer {
	return provider.ticketIssuer
}

func (provider *Provider) GetEventManager() IEventManager {
	return provider.eventManager
}

func (provider *Provider) GetTicketSigner() TicketSigner {
	return provider.ticketSigner
}

func (provider *Provider) GetUserManager() UserManager {
	return provider.userManager
}
