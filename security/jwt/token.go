package jwt

import "github.com/google/uuid"

// This struct is an abstraction over the different
// libraries used for offline (dev/test) and online
// jwt validation.
// Missing properties can be added on demand

type Token struct {
	Email    string
	Subject  uuid.UUID
	Username string
}

type Claims struct {
	Audience          []string `json:"aud,omitempty"`
	ExpiresAt         int64    `json:"exp,omitempty"`
	Id                string   `json:"jti,omitempty"`
	IssuedAt          int64    `json:"iat,omitempty"`
	Issuer            string   `json:"iss,omitempty"`
	NotBefore         int64    `json:"nbf,omitempty"`
	Subject           string   `json:"sub,omitempty"`
	PreferredUsername string   `json:"preferred_username,omitempty"`
	Email             string   `json:"email,omitempty"`
}

func (*Claims) Valid() error {
	return nil
}
