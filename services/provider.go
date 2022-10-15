package services

// TODO
// * Check lazy build

import (
	pvtbc "github.com/ticken-ts/ticken-pvtbc-connector"
	"ticken-ticket-service/config"
	"ticken-ticket-service/infra"
	"ticken-ticket-service/repos"
)

type provider struct {
	ticketIssuer TicketIssuer
	ticketSigner TicketSigner
	eventManager EventManager
}

func NewProvider(db infra.Db, pvtbcCaller *pvtbc.Caller, tickenConfig *config.Config) (Provider, error) {
	provider := new(provider)

	repoProvider, err := repos.NewProvider(db, &tickenConfig.Database)
	if err != nil {
		return nil, err
	}

	userManager := NewUserManager()
	eventRepo := repoProvider.GetEventRepository()
	ticketRepo := repoProvider.GetTicketRepository()

	provider.eventManager = NewEventManager(eventRepo)
	provider.ticketIssuer = NewTicketIssuer(eventRepo, ticketRepo, pvtbcCaller)
	provider.ticketSigner = NewTicketSigner(eventRepo, ticketRepo, pvtbcCaller, userManager)

	return provider, nil
}

func (provider *provider) GetTicketIssuer() TicketIssuer {
	return provider.ticketIssuer
}

func (provider *provider) GetEventManager() EventManager {
	return provider.eventManager
}

func (provider *provider) GetTicketSigner() TicketSigner {
	return provider.ticketSigner
}
