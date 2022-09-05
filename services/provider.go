package services

import (
	"ticken-ticket-service/blockchain/tickenPVTBCConnector"
	"ticken-ticket-service/infra"
	"ticken-ticket-service/repositories"
	"ticken-ticket-service/utils"
)

type provider struct {
	ticketIssuer TicketIssuer
}

func NewProvider(db infra.Db, tickenConfig *utils.TickenConfig) (Provider, error) {
	provider := new(provider)

	repoProvider, err := repositories.NewProvider(db, tickenConfig)
	if err != nil {
		return nil, err
	}

	provider.ticketIssuer = NewTicketIssuer(
		repoProvider.GetEventRepository(),
		repoProvider.GetTicketRepository(),
		tickenPVTBCConnector.New(),
	)

	return provider, nil
}

func (provider *provider) GetTicketIssuer() TicketIssuer {
	return provider.ticketIssuer
}
