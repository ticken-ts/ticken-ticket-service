package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type mongoDb struct {
	connString string
	client     *mongo.Client
}

func NewMongoDb() *mongoDb {
	return new(mongoDb)
}

func (mongoDb *mongoDb) Connect(connString string) error {
	if mongoDb.IsConnected() {
		return fmt.Errorf("DB is already connected")
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(connString))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	//ping the database
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	mongoDb.client = client
	mongoDb.connString = connString

	fmt.Println("Connected to MongoDB")
	return nil
}

func (mongoDb *mongoDb) IsConnected() bool {
	return mongoDb.client != nil
}

func (mongoDb *mongoDb) GetClient() interface{} {
	return mongoDb.client
}
