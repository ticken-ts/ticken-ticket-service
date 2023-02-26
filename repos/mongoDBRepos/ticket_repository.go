package mongoDBRepos

import (
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"ticken-ticket-service/models"
)

const TicketCollectionName = "tickets"

type TicketMongoDBRepository struct {
	baseRepository
}

func NewTicketRepository(db *mongo.Client, database string) *TicketMongoDBRepository {
	return &TicketMongoDBRepository{
		baseRepository{
			dbClient:       db,
			dbName:         database,
			collectionName: TicketCollectionName,
		},
	}
}

func (r *TicketMongoDBRepository) AddTicket(ticket *models.Ticket) error {
	storeContext, cancel := r.generateOpSubcontext()
	defer cancel()

	tickets := r.getCollection()
	_, err := tickets.InsertOne(storeContext, ticket)
	if err != nil {
		return err
	}

	return nil
}

func (r *TicketMongoDBRepository) FindTicket(eventID uuid.UUID, ticketID uuid.UUID) *models.Ticket {
	findContext, cancel := r.generateOpSubcontext()
	defer cancel()

	tickets := r.getCollection()
	result := tickets.FindOne(findContext, bson.M{"event_id": eventID, "ticket_id": ticketID})

	var foundTicket models.Ticket
	err := result.Decode(&foundTicket)

	if err != nil {
		return nil
	}

	return &foundTicket
}

func (r *TicketMongoDBRepository) UpdateTicketStatus(ticket *models.Ticket) error {
	updateContext, cancel := r.generateOpSubcontext()
	defer cancel()

	tickets := r.getCollection()

	filter := bson.M{"event_id": ticket.EventID, "ticket_id": ticket.TicketID}
	update := bson.M{"status": ticket.Status}

	_, err := tickets.UpdateOne(updateContext, filter, update)
	if err != nil {
		return err
	}

	return nil
}

// GetUserTickets Get tickets for a user
func (r *TicketMongoDBRepository) GetUserTickets(userID uuid.UUID) ([]*models.Ticket, error) {
	findContext, cancel := r.generateOpSubcontext()
	defer cancel()

	tickets := r.getCollection()
	findOptions := options.Find()
	findOptions.SetLimit(100)

	filter := bson.M{"owner": userID}
	cursor, err := tickets.Find(findContext, filter, findOptions)
	if err != nil {
		return nil, err
	}

	var foundTickets []*models.Ticket
	if err = cursor.All(findContext, &foundTickets); err != nil {
		return nil, err
	}

	return foundTickets, nil
}

// UpdateTicket Replace ticket with new ticket
func (r *TicketMongoDBRepository) UpdateTicketBlockchainData(ticket *models.Ticket) error {
	updateContext, cancel := r.generateOpSubcontext()
	defer cancel()

	tickets := r.getCollection()

	filter := bson.M{"event_id": ticket.EventID, "ticket_id": ticket.TicketID}
	update := bson.M{
		"$set": bson.M{
			"status":      ticket.Status,
			"pubbc_tx_id": ticket.PubbcTxID,
			"pvtbc_tx_id": ticket.PvtbcTxID,
		},
	}

	_, err := tickets.UpdateOne(updateContext, filter, update)
	if err != nil {
		return err
	}

	return nil
}
