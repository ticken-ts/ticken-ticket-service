package sync

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"net/url"
	"ticken-ticket-service/security/auth"
	"ticken-ticket-service/utils"
)

type KeycloakUser struct {
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Enabled   bool      `json:"enabled"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	ID        uuid.UUID `json:"id"`
}

type KeycloakUserCredential struct {
	Temporary bool   `json:"temporary"`
	Type      string `json:"type"`
	Value     string `json:"value"`
}

type KeycloakHTTPClient struct {
	adminURL   string
	userType   auth.UserType
	authIssuer *auth.Issuer
}

func NewKeycloakHTTPClient(issuerURL string, userType auth.UserType, authIssuer *auth.Issuer) *KeycloakHTTPClient {
	adminURL, _ := url.JoinPath(issuerURL, "admin", "realms", getUserRealm(userType))
	return &KeycloakHTTPClient{adminURL: adminURL, authIssuer: authIssuer, userType: userType}
}

func (client *KeycloakHTTPClient) RegisterUser(firstname, lastname, password string, email string) (*KeycloakUser, error) {
	user := &KeycloakUser{
		FirstName: firstname,
		LastName:  lastname,
		Email:     email,
		Enabled:   true,
	}

	jsonUser, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	usersURL, _ := url.JoinPath(client.adminURL, "/users")
	req, err := http.NewRequest(http.MethodPost, usersURL, bytes.NewReader(jsonUser))
	if err != nil {
		return nil, err
	}

	rawJWT, err := client.authIssuer.IssueInService(getUserService(client.userType))
	if err != nil || len(rawJWT.Token) == 0 {
		return nil, fmt.Errorf("could not issue access token: %s", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", rawJWT.Token))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusCreated {
		return nil, errors.New(utils.ReadHTTPResponseBody(res))
	}

	refreshedUser, err := client.GetUserByEmail(email)
	if err != nil {
		return nil, fmt.Errorf("could not refresh user: %s", email)
	}

	if err := client.UpdateUserPassword(refreshedUser.ID, password); err != nil {
		return nil, err
	}

	return refreshedUser, nil
}

func (client *KeycloakHTTPClient) GetUserByEmail(email string) (*KeycloakUser, error) {
	params := url.Values{}
	params.Set("exact", "true")
	params.Set("username", email) // realm uses username as email

	usersURL, _ := url.JoinPath(client.adminURL, "/users")
	req, err := http.NewRequest(http.MethodGet, usersURL+"?"+params.Encode(), http.NoBody)
	if err != nil {
		return nil, err
	}

	rawJWT, err := client.authIssuer.IssueInService(getUserService(client.userType))
	if err != nil {
		return nil, fmt.Errorf("could not issue access token: %s", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", rawJWT.Token))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, errors.New(utils.ReadHTTPResponseBody(res))
	}

	var users []KeycloakUser
	err = json.Unmarshal([]byte(utils.ReadHTTPResponseBody(res)), &users)
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, fmt.Errorf("user with email %s not found", email)
	}

	return &users[0], nil
}

func (client *KeycloakHTTPClient) UpdateUserPassword(userID uuid.UUID, newPassword string) error {
	credentials := &KeycloakUserCredential{
		Temporary: false,
		Type:      "password",
		Value:     newPassword,
	}

	jsonCredentials, err := json.Marshal(credentials)
	if err != nil {
		return err
	}

	usersURL, _ := url.JoinPath(client.adminURL, "/users")
	resetPassURL, _ := url.JoinPath(usersURL, userID.String(), "reset-password")
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPut, resetPassURL, bytes.NewReader(jsonCredentials))
	if err != nil {
		return err
	}

	rawJWT, err := client.authIssuer.IssueInService(getUserService(client.userType))
	if err != nil {
		return fmt.Errorf("could not issue access token: %s", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", rawJWT.Token))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusNoContent {
		return errors.New(utils.ReadHTTPResponseBody(res))
	}

	return nil
}

func getUserRealm(forUser auth.UserType) string {
	var realm = ""
	switch forUser {
	case auth.Validator:
		realm = "validators"
		break
	case auth.Attendant:
		realm = "attendants"
		break
	case auth.Organizer:
		realm = "organizers"
		break
	}
	return realm
}

func getUserService(forUser auth.UserType) auth.ServiceType {
	switch forUser {
	case auth.Validator:
		return auth.TickenValidatorService
	case auth.Attendant:
		return auth.TickenTicketService
	case auth.Organizer:
		return auth.TickenEventService
	default:
		return 0
	}
}
