package services

// TODO
// * Check lazy build

import (
	pvtbc "github.com/ticken-ts/ticken-pvtbc-connector"
	"ticken-ticket-service/repos"
)

type Provider struct {
	ticketIssuer TicketIssuer
	ticketSigner TicketSigner
	eventManager EventManager
}

func NewProvider(repoProvider repos.IProvider, pvtbcCaller *pvtbc.Caller) (*Provider, error) {
	provider := new(Provider)

	userManager := NewUserManager()

	eventRepo := repoProvider.GetEventRepository()
	ticketRepo := repoProvider.GetTicketRepository()

	provider.eventManager = NewEventManager(eventRepo)
	provider.ticketIssuer = NewTicketIssuer(eventRepo, ticketRepo, pvtbcCaller)
	provider.ticketSigner = NewTicketSigner(eventRepo, ticketRepo, pvtbcCaller, userManager)

	return provider, nil
}

func (provider *Provider) GetTicketIssuer() TicketIssuer {
	return provider.ticketIssuer
}

func (provider *Provider) GetEventManager() EventManager {
	return provider.eventManager
}

func (provider *Provider) GetTicketSigner() TicketSigner {
	return provider.ticketSigner
}
