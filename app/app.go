package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"ticken-ticket-service/api"
	"ticken-ticket-service/api/controllers/healthController"
	"ticken-ticket-service/api/controllers/ticketController"
	"ticken-ticket-service/api/middlewares"
	"ticken-ticket-service/api/security"
	"ticken-ticket-service/async"
	"ticken-ticket-service/config"
	"ticken-ticket-service/env"
	"ticken-ticket-service/infra"
	"ticken-ticket-service/log"
	"ticken-ticket-service/repos"
	"ticken-ticket-service/services"
	"ticken-ticket-service/sync"
	"ticken-ticket-service/utils"
)

const ServiceName = "ticken-ticket-service"

type TickenTicketApp struct {
	engine          *gin.Engine
	config          *config.Config
	repoProvider    repos.IProvider
	serviceProvider services.IProvider
	subscriber      *async.Subscriber

	// populators are intended to populate
	// useful data. It can be testdata or
	// data that should be present on the db
	// before the service is available
	populators []Populator
}

func New(builder infra.IBuilder, tickenConfig *config.Config) *TickenTicketApp {
	log.TickenLogger.Info().Msg("initializing " + ServiceName)

	tickenTicketApp := new(TickenTicketApp)

	db := builder.BuildDb(env.TickenEnv.DbConnString)
	engine := builder.BuildEngine()
	pvtbcCaller := builder.BuildPvtbcCaller()
	publicBlockchain := builder.BuildPublicBlockchain()
	busSubscriber := builder.BuildBusSubscriber(env.TickenEnv.BusConnString)

	// this provider is going to provider all repositories
	// to the services
	repoProvider, err := repos.NewProvider(db, &tickenConfig.Database)
	if err != nil {
		panic(err)
	}

	// this provider is going to provide all services
	// needed by the controllers to execute it operations
	serviceProvider, err := services.NewProvider(repoProvider, pvtbcCaller, sync.NewUserServiceClient(), publicBlockchain)
	if err != nil {
		panic(err)
	}

	err = serviceProvider.GetEventManager().ListenBlockchainEvents()
	if err != nil {
		panic(err)
	}

	subscriber, err := async.NewSubscriber(busSubscriber, serviceProvider)
	if err != nil {
		panic(err)
	}

	err = subscriber.Start()
	if err != nil {
		panic(err)
	}

	tickenTicketApp.engine = engine
	tickenTicketApp.config = tickenConfig
	tickenTicketApp.subscriber = subscriber
	tickenTicketApp.repoProvider = repoProvider
	tickenTicketApp.serviceProvider = serviceProvider

	var appMiddlewares = []api.Middleware{
		middlewares.NewAuthMiddleware(serviceProvider, &tickenConfig.Server, &tickenConfig.Dev),
	}

	for _, middleware := range appMiddlewares {
		middleware.Setup(engine)
	}

	var appControllers = []api.Controller{
		healthController.New(serviceProvider),
		ticketController.New(serviceProvider),
	}

	for _, controller := range appControllers {
		controller.Setup(engine)
	}

	return tickenTicketApp
}

func (ticketTicketApp *TickenTicketApp) Start() {
	url := ticketTicketApp.config.Server.GetServerURL()
	err := ticketTicketApp.engine.Run(url)
	if err != nil {
		panic(err)
	}
}

func (ticketTicketApp *TickenTicketApp) Populate() {
	for _, populator := range ticketTicketApp.populators {
		err := populator.Populate()
		if err != nil {
			panic(err)
		}
	}
}

func (ticketTicketApp *TickenTicketApp) EmitFakeJWT() {
	rsaPrivKey, err := utils.LoadRSA(ticketTicketApp.config.Dev.JWTPrivateKey, ticketTicketApp.config.Dev.JWTPublicKey)
	if err != nil {
		panic(err)
	}

	fakeJWT := jwt.NewWithClaims(jwt.SigningMethodRS256, &security.Claims{
		Subject:           ticketTicketApp.config.Dev.User.UserID,
		Email:             ticketTicketApp.config.Dev.User.Email,
		PreferredUsername: ticketTicketApp.config.Dev.User.Username,
	})

	signedJWT, err := fakeJWT.SignedString(rsaPrivKey)

	if err != nil {
		panic(fmt.Errorf("error generation fake JWT: %s", err.Error()))
	}

	fmt.Printf("DEV JWT: %s \n", signedJWT)
}
