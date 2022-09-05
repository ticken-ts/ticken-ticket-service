package mongoDBRepositories

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"ticken-ticket-service/models"
	"time"
)

const TimeoutOperationsInSeconds = 10
const EventCollectionName = "events"

type eventMongoDBRepository struct {
	database string
	db       *mongo.Client
}

func NewEventRepository(db *mongo.Client, database string) *eventMongoDBRepository {
	return &eventMongoDBRepository{
		db:       db,
		database: database,
	}
}

func (r *eventMongoDBRepository) getCollection() *mongo.Collection {
	return r.db.Database(r.database).Collection(EventCollectionName)
}

func (r *eventMongoDBRepository) generateOpSubcontext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), TimeoutOperationsInSeconds*time.Second)
}

func (r *eventMongoDBRepository) AddEvent(event *models.Event) error {
	storeContext, cancel := r.generateOpSubcontext()
	defer cancel()

	events := r.getCollection()
	_, err := events.InsertOne(storeContext, event)
	if err != nil {
		return err
	}

	return nil
}

func (r *eventMongoDBRepository) FindEventByID(eventID string) *models.Event {
	findContext, cancel := r.generateOpSubcontext()
	defer cancel()

	events := r.getCollection()
	result := events.FindOne(findContext, bson.M{"_id": eventID})

	var foundEvent models.Event
	err := result.Decode(&foundEvent)

	if err != nil {
		return nil
	}

	return &foundEvent
}
