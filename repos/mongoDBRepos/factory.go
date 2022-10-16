package mongoDBRepos

import (
	"go.mongodb.org/mongo-driver/mongo"
	"ticken-ticket-service/config"
	"ticken-ticket-service/infra"
)

type MongoRepoFactory struct {
	dbClient *mongo.Client
	dbName   string
}

func NewFactory(db infra.Db, dbConfig *config.DatabaseConfig) *MongoRepoFactory {
	return &MongoRepoFactory{
		dbClient: db.GetClient().(*mongo.Client),
		dbName:   dbConfig.Name,
	}
}

func (factory *MongoRepoFactory) BuildEventRepository() any {
	return NewEventRepository(factory.dbClient, factory.dbName)
}

func (factory *MongoRepoFactory) BuildTicketRepository() any {
	return NewTicketRepository(factory.dbClient, factory.dbName)
}
