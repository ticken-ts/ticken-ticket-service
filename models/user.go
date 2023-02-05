package models

import (
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	mongoID primitive.ObjectID `bson:"_id"`
	UUID    uuid.UUID          `json:"uuid" bson:"uuid"`
}

func NewUser(id uuid.UUID) *User {
	return &User{
		UUID: id,
	}
}
