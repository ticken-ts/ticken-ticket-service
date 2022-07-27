package repositories

import "go.mongodb.org/mongo-driver/mongo"

type ticketMongoDBRepository struct {
	database string
	db       *mongo.Client
}

func NewTicketMongoDBRepository(db *mongo.Client, database string) TicketRepository {
	return &ticketMongoDBRepository{
		db:       db,
		database: database}
}
