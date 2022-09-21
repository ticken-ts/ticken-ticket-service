package mongoDBRepos

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

const TimeoutOperationsInSeconds = 10

type baseRepository struct {
	dbName         string
	dbClient       *mongo.Client
	collectionName string
}

func (r *baseRepository) getCollection() *mongo.Collection {
	return r.dbClient.Database(r.dbName).Collection(r.collectionName)
}

func (r *baseRepository) generateOpSubcontext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), TimeoutOperationsInSeconds*time.Second)
}

func (r *baseRepository) findByMongoID(id primitive.ObjectID) *mongo.SingleResult {
	findContext, cancel := r.generateOpSubcontext()
	defer cancel()

	events := r.getCollection()
	result := events.FindOne(findContext, bson.M{"_id": id})

	return result
}
