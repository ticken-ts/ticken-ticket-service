package repos

import (
	"fmt"
	"ticken-ticket-service/config"
	"ticken-ticket-service/infra"
	"ticken-ticket-service/repos/mongoDBRepos"
)

type Provider struct {
	reposFactory     IFactory
	eventRepository  EventRepository
	ticketRepository TicketRepository
	userRepository   UserRepository
}

func NewProvider(db infra.Db, dbConfig *config.DatabaseConfig) (*Provider, error) {
	provider := new(Provider)

	switch dbConfig.Driver {
	case config.MongoDriver:
		provider.reposFactory = mongoDBRepos.NewFactory(db, dbConfig)
	default:
		return nil, fmt.Errorf("database driver %s not implemented", dbConfig.Driver)
	}

	return provider, nil
}

func (provider *Provider) GetEventRepository() EventRepository {
	if provider.eventRepository == nil {
		provider.eventRepository = provider.reposFactory.BuildEventRepository().(EventRepository)
	}
	return provider.eventRepository
}

func (provider *Provider) GetTicketRepository() TicketRepository {
	if provider.ticketRepository == nil {
		provider.ticketRepository = provider.reposFactory.BuildTicketRepository().(TicketRepository)
	}
	return provider.ticketRepository
}

func (provider *Provider) GetUserRepository() UserRepository {
	if provider.userRepository == nil {
		provider.userRepository = provider.reposFactory.BuildUserRepository().(UserRepository)
	}
	return provider.userRepository
}
