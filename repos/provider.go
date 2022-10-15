package repos

import (
	"fmt"
	"ticken-ticket-service/config"
	"ticken-ticket-service/infra"
	"ticken-ticket-service/repos/mongoDBRepos"
)

type provider struct {
	reposFactory     Factory
	eventRepository  EventRepository
	ticketRepository TicketRepository
}

func NewProvider(db infra.Db, dbConfig *config.DatabaseConfig) (Provider, error) {
	provider := new(provider)

	switch dbConfig.Driver {
	case config.MongoDriver:
		provider.reposFactory = mongoDBRepos.NewFactory(db, dbConfig)
	default:
		return nil, fmt.Errorf("database driver %s not implemented", dbConfig.Driver)
	}

	return provider, nil
}

func (provider *provider) GetEventRepository() EventRepository {
	if provider.eventRepository == nil {
		provider.eventRepository = provider.reposFactory.BuildEventRepository().(EventRepository)
	}
	return provider.eventRepository
}

func (provider *provider) GetTicketRepository() TicketRepository {
	if provider.ticketRepository == nil {
		provider.ticketRepository = provider.reposFactory.BuildTicketRepository().(TicketRepository)
	}
	return provider.ticketRepository
}
