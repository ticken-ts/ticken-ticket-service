package services

// TODO
// * Check lazy build

import (
	pvtbc "github.com/ticken-ts/ticken-pvtbc-connector"
	"ticken-ticket-service/infra/public_blockchain"
	"ticken-ticket-service/repos"
)

type Provider struct {
	ticketIssuer TicketIssuer
	ticketSigner TicketSigner
	eventManager EventManager
	userManager  UserManager
}

func NewProvider(repoProvider repos.IProvider, pvtbcCaller *pvtbc.Caller, publicBlockchain public_blockchain.PublicBC) (*Provider, error) {
	provider := new(Provider)

	eventRepo := repoProvider.GetEventRepository()
	ticketRepo := repoProvider.GetTicketRepository()
	userRepo := repoProvider.GetUserRepository()

	provider.eventManager = NewEventManager(eventRepo, ticketRepo, publicBlockchain)
	provider.ticketIssuer = NewTicketIssuer(eventRepo, ticketRepo, pvtbcCaller, publicBlockchain)
	provider.ticketSigner = NewTicketSigner(eventRepo, ticketRepo, pvtbcCaller, publicBlockchain)
	provider.userManager = NewUserManager(eventRepo, ticketRepo, userRepo, publicBlockchain)

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

func (provider *Provider) GetUserManager() UserManager {
	return provider.userManager
}
