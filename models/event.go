package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Event struct {
	ID       primitive.ObjectID `json:"id"`
	IsActive bool               `json:"is_active" validate:"required"`
}
