package repos

import (
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"ticken-ticket-service/infra"
	"ticken-ticket-service/repos/mongoDBRepositories"
	"ticken-ticket-service/utils"
)

type provider struct {
	eventRepository  EventRepository
	ticketRepository TicketRepository
}

func NewProvider(db infra.Db, tickenConfig *utils.TickenConfig) (*provider, error) {
	provider := new(provider)

	switch tickenConfig.Config.Database.Driver {
	case utils.MongoDriver:
		mongoDbClient := db.GetClient().(*mongo.Client)

		provider.eventRepository = mongoDBRepositories.NewEventRepository(
			mongoDbClient,
			tickenConfig.Config.Database.Name)

		provider.ticketRepository = mongoDBRepositories.NewTicketRepository(
			mongoDbClient,
			tickenConfig.Config.Database.Name)

	default:
		return nil, fmt.Errorf("database driver %s not implemented", tickenConfig.Config.Database.Driver)
	}

	return provider, nil
}

func (provider *provider) GetEventRepository() EventRepository {
	return provider.eventRepository
}

func (provider *provider) GetTicketRepository() TicketRepository {
	return provider.ticketRepository
}
