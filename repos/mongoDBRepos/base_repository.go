package mongoDBRepos

import (
	"context"
	"github.com/google/uuid"
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
	primKeyField   string
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

func (r *baseRepository) Count() int64 {
	findContext, cancel := r.generateOpSubcontext()
	defer cancel()

	collection := r.getCollection()
	totalDocs, err := collection.CountDocuments(findContext, bson.M{})
	if err != nil {
		return 0
	}
	return totalDocs
}

func (r *baseRepository) AddOne(model any) error {
	storeContext, cancel := r.generateOpSubcontext()
	defer cancel()

	collection := r.getCollection()
	_, err := collection.InsertOne(storeContext, model)
	if err != nil {
		return err
	}

	return nil
}

func (r *baseRepository) AnyWithID(id uuid.UUID) bool {
	findContext, cancel := r.generateOpSubcontext()
	defer cancel()

	collection := r.getCollection()
	totalDocs, err := collection.CountDocuments(findContext, bson.M{r.primKeyField: id})
	if err != nil {
		return false
	}

	return totalDocs > 0
}
