package models

import (
	"container/list"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Event struct {
	EventID  primitive.ObjectID `json:"event_id"`
	Name     string             `json:"name"`
	Date     time.Time          `json:"date"`
	Sections list.List          `json:"sections"`
}
