package mongoDBRepositories

import (
	"go.mongodb.org/mongo-driver/mongo"
)

type ticketMongoDBRepository struct {
	database string
	db       *mongo.Client
}

func NewTicketRepository(db *mongo.Client, database string) *ticketMongoDBRepository {
	return &ticketMongoDBRepository{
		db:       db,
		database: database}
}
