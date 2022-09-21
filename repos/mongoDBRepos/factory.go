package mongoDBRepos

import (
	"go.mongodb.org/mongo-driver/mongo"
	"ticken-ticket-service/infra"
	"ticken-ticket-service/utils"
)

type Factory struct {
	dbClient *mongo.Client
	dbName   string
}

func NewFactory(db infra.Db, tf *utils.TickenConfig) *Factory {
	return &Factory{
		dbClient: db.GetClient().(*mongo.Client),
		dbName:   tf.Config.Database.Name,
	}
}

func (factory *Factory) BuildEventRepository() any {
	return NewEventRepository(factory.dbClient, factory.dbName)
}

func (factory *Factory) BuildTicketRepository() any {
	return NewTicketRepository(factory.dbClient, factory.dbName)
}
