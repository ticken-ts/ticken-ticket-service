package models

import "github.com/google/uuid"

type User struct {
	ID uuid.UUID `json:"id"`
}

func NewUser(id uuid.UUID) *User {
	return &User{
		ID: id,
	}
}
