package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Ticket struct {
	ID      primitive.ObjectID `json:"id"`
	Owner   string             `json:"owner"`
	EventID string             `json:"event_id"`
	Section string             `json:"section"`
}
