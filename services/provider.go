package services

import (
	pvtbc "github.com/ticken-ts/ticken-pvtbc-connector"
	"ticken-ticket-service/infra"
	"ticken-ticket-service/repos"
	"ticken-ticket-service/utils"
)

type provider struct {
	ticketIssuer TicketIssuer
	ticketSigner TicketSigner
	eventManager EventManager
}

func NewProvider(db infra.Db, pvtbcCaller *pvtbc.Caller, tickenConfig *utils.TickenConfig) (Provider, error) {
	provider := new(provider)

	repoProvider, err := repos.NewProvider(db, tickenConfig)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	eventRepo := repoProvider.GetEventRepository()
	ticketRepo := repoProvider.GetTicketRepository()
	userManager := NewUserManager()

	provider.eventManager = NewEventManager(repoProvider.GetEventRepository())
	provider.ticketIssuer = NewTicketIssuer(eventRepo, ticketRepo, pvtbcCaller)
	provider.ticketSigner = NewTicketSigner(eventRepo, ticketRepo, pvtbcCaller, userManager)

	return provider, nil
}

// TODO -> see if it is better to do lazy
func (provider *provider) GetTicketIssuer() TicketIssuer {
	return provider.ticketIssuer
}

func (provider *provider) GetEventManager() EventManager {
	return provider.eventManager
}

func (provider *provider) GetTicketSigner() TicketSigner {
	return provider.ticketSigner
}
