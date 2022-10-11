package repos

import (
	"fmt"
	"ticken-ticket-service/infra"
	"ticken-ticket-service/repos/mongoDBRepos"
	"ticken-ticket-service/utils"
)

type provider struct {
	reposFactory     Factory
	eventRepository  EventRepository
	ticketRepository TicketRepository
}

func NewProvider(db infra.Db, tickenConfig *utils.TickenConfig) (Provider, error) {
	provider := new(provider)

	switch tickenConfig.Config.Database.Driver {
	case utils.MongoDriver:
		provider.reposFactory = mongoDBRepos.NewFactory(db, &tickenConfig.Config.Database)
	default:
		return nil, fmt.Errorf("database driver %s not implemented", tickenConfig.Config.Database.Driver)
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
