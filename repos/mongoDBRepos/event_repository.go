package mongoDBRepos

import (
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"ticken-ticket-service/models"
)

const EventCollectionName = "events"

type EventMongoDBRepository struct {
	baseRepository
}

func NewEventRepository(dbClient *mongo.Client, dbName string) *EventMongoDBRepository {
	return &EventMongoDBRepository{
		baseRepository{
			dbClient:       dbClient,
			dbName:         dbName,
			collectionName: EventCollectionName,
			primKeyField:   "event_id",
		},
	}
}

func (r *EventMongoDBRepository) FindEvent(eventID uuid.UUID) *models.Event {
	findContext, cancel := r.generateOpSubcontext()
	defer cancel()

	events := r.getCollection()
	result := events.FindOne(findContext, bson.M{"event_id": eventID})

	var foundEvent models.Event
	err := result.Decode(&foundEvent)

	if err != nil {
		return nil
	}

	return &foundEvent
}

func (r *EventMongoDBRepository) FindEventByContractAddress(contractAddr string) *models.Event {
	findContext, cancel := r.generateOpSubcontext()
	defer cancel()

	events := r.getCollection()
	result := events.FindOne(findContext, bson.M{"pub_bc_address": contractAddr})

	var foundEvent models.Event
	err := result.Decode(&foundEvent)

	if err != nil {
		return nil
	}

	return &foundEvent
}

func (r *EventMongoDBRepository) GetActiveEvents() ([]*models.Event, error) {
	findContext, cancel := r.generateOpSubcontext()
	defer cancel()

	events := r.getCollection()
	cursor, err := events.Find(findContext, bson.M{})
	if err != nil {
		return nil, err
	}

	var foundEvents []*models.Event
	err = cursor.All(findContext, &foundEvents)
	if err != nil {
		return nil, err
	}

	return foundEvents, nil
}
