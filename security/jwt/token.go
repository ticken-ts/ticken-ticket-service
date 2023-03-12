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
	Profile  Profile
}

type Profile struct {
	FirstName     string
	LastName      string
	EmailVerified bool
}

type Claims struct {
	NotBefore         int64          `json:"nbf,omitempty"`
	ExpiresAt         int64          `json:"exp"`
	IssuedAt          int64          `json:"iat"`
	AuthTime          int64          `json:"auth_time"`
	Id                string         `json:"jti"`
	Issuer            string         `json:"iss"`
	Audience          []string       `json:"aud"`
	Subject           string         `json:"sub"`
	Typ               string         `json:"typ"`
	Azp               string         `json:"azp"`
	SessionState      string         `json:"session_state"`
	AllowedOrigins    []string       `json:"allowed-origins"`
	RealmAccess       RealmAccess    `json:"realm_access"`
	ResourceAccess    ResourceAccess `json:"resource_access"`
	Scope             string         `json:"scope"`
	Sid               string         `json:"sid"`
	EmailVerified     bool           `json:"email_verified"`
	Name              string         `json:"name"`
	PreferredUsername string         `json:"preferred_username"`
	GivenName         string         `json:"given_name"`
	FamilyName        string         `json:"family_name"`
	Email             string         `json:"email"`
}

type RealmAccess struct {
	Roles []string `json:"roles"`
}

type ResourceAccess struct {
	Account RealmAccess `json:"account"`
}

func (*Claims) Valid() error {
	return nil
}
