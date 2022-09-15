package services

import (
	"ticken-ticket-service/blockchain/tickenPVTBCConnector"
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

	provider.ticketIssuer = NewTicketIssuer(
		repoProvider.GetEventRepository(),
		repoProvider.GetTicketRepository(),
		tickenPVTBCConnector.New(),
	)

	provider.eventManager = NewEventManager(
		repoProvider.GetEventRepository(),
		repoProvider.GetTicketRepository(),
		tickenPVTBCConnector.New(),
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
