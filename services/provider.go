package services

import (
	pubbc "github.com/ticken-ts/ticken-pubbc-connector"
	pvtbc "github.com/ticken-ts/ticken-pvtbc-connector"
	"ticken-ticket-service/async"
	"ticken-ticket-service/config"
	"ticken-ticket-service/env"
	"ticken-ticket-service/infra"
	"ticken-ticket-service/repos"
	"ticken-ticket-service/security/auth"
	"ticken-ticket-service/sync"
)

type Provider struct {
	ticketIssuer ITicketIssuer
	ticketLinker ITicketLinker
	eventManager IEventManager
	userManager  IUserManager
	ticketTrader ITicketTrader
}

func NewProvider(
	repoProvider repos.IProvider,
	pvtbcCaller *pvtbc.Caller,
	pubbcAdmin pubbc.Admin,
	pubbcCaller pubbc.Caller,
	asyncPublisher async.IAsyncPublisher,
	hsm infra.HSM,
	authIssuer *auth.Issuer,
	tickenConfig *config.Config,
) (*Provider, error) {
	provider := new(Provider)

	var attendantsKeycloakClient *sync.KeycloakHTTPClient
	if !env.TickenEnv.IsDev() || tickenConfig.Dev.Mock.DisableAuthMock {
		attendantsKeycloakClient = sync.NewKeycloakHTTPClient(tickenConfig.Services.Keycloak, auth.Attendant, authIssuer)
	}

	provider.eventManager = NewEventManager(repoProvider)
	provider.ticketLinker = NewTicketLinker(repoProvider, pubbcCaller)
	provider.ticketTrader = NewTicketTrader(repoProvider, pubbcCaller)
	provider.ticketIssuer = NewTicketIssuer(repoProvider, hsm, pubbcCaller, pvtbcCaller)
	provider.userManager = NewUserManager(repoProvider, pubbcAdmin, hsm, asyncPublisher, attendantsKeycloakClient)

	return provider, nil
}

func (provider *Provider) GetTicketIssuer() ITicketIssuer {
	return provider.ticketIssuer
}

func (provider *Provider) GetEventManager() IEventManager {
	return provider.eventManager
}

func (provider *Provider) GetUserManager() IUserManager {
	return provider.userManager
}

func (provider *Provider) GetTicketLinker() ITicketLinker {
	return provider.ticketLinker
}

func (provider *Provider) GetTicketTrader() ITicketTrader {
	return provider.ticketTrader
}
