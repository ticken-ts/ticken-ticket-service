package services

import (
	"ticken-ticket-service/blockchain/pvtbc"
	"ticken-ticket-service/infra"
	"ticken-ticket-service/repos"
	"ticken-ticket-service/utils"
)

type provider struct {
	ticketIssuer TicketIssuer
	eventManager EventManager
}

func NewProvider(db infra.Db, tickenConfig *utils.TickenConfig) (Provider, error) {
	provider := new(provider)

	repoProvider, err := repos.NewProvider(db, tickenConfig)
	if err != nil {
		return nil, err
	}

	pvtbcTickenConnector, err := pvtbc.NewConnector()
	if err != nil {
		return nil, err
	}

	provider.ticketIssuer = NewTicketIssuer(
		repoProvider.GetEventRepository(),
		repoProvider.GetTicketRepository(),
		pvtbcTickenConnector,
	)

	provider.eventManager = NewEventManager(
		repoProvider.GetEventRepository(),
		repoProvider.GetTicketRepository(),
		pvtbcTickenConnector,
	)

	return provider, nil
}

// TODO -> see if it is better to do lazy
func (provider *provider) GetTicketIssuer() TicketIssuer {
	return provider.ticketIssuer
}

func (provider *provider) GetEventManager() EventManager {
	return provider.eventManager
}
