package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type UserType int8
type ServiceType int8

const (
	TickenValidatorService ServiceType = 1
	TickenTicketService                = 2
	TickenEventService                 = 3
)

const (
	Validator UserType = 1
	Attendant          = 2
	Organizer          = 3
)

type RawJWT struct {
	Token string `json:"access_token"`
}

type Issuer struct {
	issuerURL    string
	clientID     string
	clientSecret string
	selfService  ServiceType
}

// NewAuthIssuer creates a new JWT issues based on the
// authority provider set up. This constructor receives
// three parameters:
// @ selfService -> is the service which is running
// @ issuerURL   -> base url where the identity issues is running
// @ clientID    -> client ID of the service specified in selfService
func NewAuthIssuer(selfService ServiceType, issuerURL, clientID, clientSecret string) (*Issuer, error) {
	return &Issuer{
		issuerURL:    issuerURL,
		clientID:     clientID,
		clientSecret: clientSecret,
		selfService:  selfService,
	}, nil
}

func (issuer *Issuer) IssueForSelf() (*RawJWT, error) {
	data := url.Values{}
	data.Set("client_id", issuer.clientID)
	data.Set("client_secret", issuer.clientSecret)
	data.Set("grant_type", "client_credentials")

	realm := getServiceRealm(issuer.selfService)
	if len(realm) == 0 {
		return nil, fmt.Errorf("could not determine service realm")
	}

	tokenURL, _ := url.JoinPath(issuer.issuerURL, "realms", realm, "/protocol/openid-connect/token")

	response, err := http.Post(
		tokenURL,
		"application/x-www-form-urlencoded",
		bytes.NewReader([]byte(data.Encode())),
	)
	if err != nil {
		return nil, err
	}

	bodyData, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	rawJWT := new(RawJWT)
	if err := json.Unmarshal(bodyData, rawJWT); err != nil {
		return nil, err
	}

	return rawJWT, nil
}

func (issuer *Issuer) IssueInService(serviceType ServiceType) (*RawJWT, error) {
	data := url.Values{}
	data.Set("client_id", issuer.clientID)
	data.Set("client_secret", issuer.clientSecret)
	data.Set("grant_type", "client_credentials")

	realm := getServiceRealm(serviceType)
	if len(realm) == 0 {
		return nil, fmt.Errorf("could not determine service realm")
	}

	tokenURL, _ := url.JoinPath(issuer.issuerURL, "realms", realm, "/protocol/openid-connect/token")

	response, err := http.Post(
		tokenURL,
		"application/x-www-form-urlencoded",
		bytes.NewReader([]byte(data.Encode())),
	)
	if err != nil {
		return nil, err
	}

	bodyData, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	rawJWT := new(RawJWT)
	if err := json.Unmarshal(bodyData, rawJWT); err != nil {
		return nil, err
	}

	return rawJWT, nil
}

func (issuer *Issuer) IssueForUser(userType UserType, email, password string) (*RawJWT, error) {
	data := url.Values{}
	data.Set("username", email)
	data.Set("password", password)
	data.Set("client_id", issuer.clientID)
	data.Set("client_secret", issuer.clientSecret)
	data.Set("grant_type", "password")

	realm := getUserRealm(userType)
	if len(realm) == 0 {
		return nil, fmt.Errorf("could not determine user realm")
	}

	tokenURL, _ := url.JoinPath(issuer.issuerURL, "realms", realm, "/protocol/openid-connect/token")

	response, err := http.Post(
		tokenURL,
		"application/x-www-form-urlencoded",
		bytes.NewReader([]byte(data.Encode())),
	)
	if err != nil {
		return nil, err
	}

	bodyData, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	rawJWT := new(RawJWT)
	if err := json.Unmarshal(bodyData, rawJWT); err != nil {
		return nil, err
	}

	return rawJWT, nil
}

func getServiceRealm(forService ServiceType) string {
	var realm = ""
	switch forService {
	case TickenValidatorService:
		realm = "validators"
		break
	case TickenTicketService:
		realm = "attendants"
		break
	case TickenEventService:
		realm = "organizers"
		break
	}
	return realm
}

func getUserRealm(forUser UserType) string {
	var realm = ""
	switch forUser {
	case Validator:
		realm = "validators"
		break
	case Attendant:
		realm = "attendants"
		break
	case Organizer:
		realm = "organizers"
		break
	}
	return realm
}
