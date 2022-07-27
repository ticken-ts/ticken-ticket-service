package repositories

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"ticken-ticket-service/models/event"
	"time"
)

const TimeoutOperationsInSeconds = 10
const EventCollectionName = "events"

type eventMongoDBRepository struct {
	database string
	db       *mongo.Client
}

func NewEventMongoDBRepository(db *mongo.Client, database string) EventRepository {
	return &eventMongoDBRepository{
		db:       db,
		database: database}
}

func (r *eventMongoDBRepository) getCollection() *mongo.Collection {
	return r.db.Database(r.database).Collection(EventCollectionName)
}

func (r *eventMongoDBRepository) generateOpSubcontext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), TimeoutOperationsInSeconds*time.Second)
}

func (r *eventMongoDBRepository) AddEvent(event *event.Event) error {
	storeContext, cancel := r.generateOpSubcontext()
	defer cancel()

	events := r.getCollection()
	_, err := events.InsertOne(storeContext, event)
	if err != nil {
		return err
	}

	return nil
}

func (r *eventMongoDBRepository) FindEventByID(eventID string) *event.Event {
	findContext, cancel := r.generateOpSubcontext()
	defer cancel()

	events := r.getCollection()
	result := events.FindOne(findContext, bson.M{"eventID": eventID})

	var foundEvent event.Event
	err := result.Decode(&foundEvent)

	if err != nil {
		return nil
	}

	return &foundEvent
}
