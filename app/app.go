package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	gojwt "github.com/golang-jwt/jwt"
	"ticken-ticket-service/api"
	"ticken-ticket-service/api/controllers/healthController"
	"ticken-ticket-service/api/controllers/ticketController"
	"ticken-ticket-service/api/controllers/userController"
	"ticken-ticket-service/api/middlewares"
	"ticken-ticket-service/app/fakes"
	"ticken-ticket-service/async"
	"ticken-ticket-service/config"
	"ticken-ticket-service/env"
	"ticken-ticket-service/infra"
	"ticken-ticket-service/log"
	"ticken-ticket-service/repos"
	"ticken-ticket-service/security/jwt"
	"ticken-ticket-service/services"
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

func New(infraBuilder infra.IBuilder, tickenConfig *config.Config) *TickenTicketApp {
	log.TickenLogger.Info().Msg("initializing " + ServiceName)

	tickenTicketApp := new(TickenTicketApp)

	engine := infraBuilder.BuildEngine()
	pvtbcCaller := infraBuilder.BuildPvtbcCaller()
	jwtVerifier := infraBuilder.BuildJWTVerifier()
	db := infraBuilder.BuildDb(env.TickenEnv.DbConnString)
	hsm := infraBuilder.BuildHSM(env.TickenEnv.HSMEncryptionKey)
	pubbcAdmin := infraBuilder.BuildPubbcAdmin(env.TickenEnv.TickenWalletKey)
	pubbcCaller := infraBuilder.BuildPubbcCaller(env.TickenEnv.TickenWalletKey)
	busSubscriber := infraBuilder.BuildBusSubscriber(env.TickenEnv.BusConnString)

	// this provider is going to provider all repositories
	// to the services
	repoProvider, err := repos.NewProvider(db, &tickenConfig.Database)
	if err != nil {
		panic(err)
	}

	// this provider is going to provide all services
	// needed by the controllers to execute it operations
	serviceProvider, err := services.NewProvider(repoProvider, pvtbcCaller, pubbcAdmin, pubbcCaller, hsm)
	if err != nil {
		panic(err)
	}

	subscriber, err := async.NewSubscriber(busSubscriber, serviceProvider)
	if err != nil {
		panic(err)
	}

	tickenTicketApp.engine = engine
	tickenTicketApp.config = tickenConfig
	tickenTicketApp.subscriber = subscriber
	tickenTicketApp.repoProvider = repoProvider
	tickenTicketApp.serviceProvider = serviceProvider

	var appMiddlewares = []api.Middleware{
		middlewares.NewErrorMiddleware(),
		middlewares.NewAuthMiddleware(serviceProvider, jwtVerifier, tickenConfig.Server.APIPrefix),
	}

	var appControllers = []api.Controller{
		healthController.New(serviceProvider),
		ticketController.New(serviceProvider),
		userController.New(serviceProvider),
	}

	apiRouter := engine.Group(tickenConfig.Server.APIPrefix)

	for _, middleware := range appMiddlewares {
		middleware.Setup(apiRouter)
	}

	for _, controller := range appControllers {
		controller.Setup(apiRouter)
	}

	var appPopulators = []Populator{
		//&fakes.FakeEventsPopulator{EventRepo: repoProvider.GetEventRepository(), Pubbc: pubbcAdmin},
		&fakes.FakeUsersPopulator{Repo: repoProvider.GetUserRepository(), Config: tickenConfig.Dev.User, HSM: hsm},
	}
	tickenTicketApp.populators = appPopulators

	return tickenTicketApp
}

func (ticketTicketApp *TickenTicketApp) Start() {
	url := ticketTicketApp.config.Server.GetServerURL()

	if err := ticketTicketApp.subscriber.Start(); err != nil {
		panic(err)
	}

	if err := ticketTicketApp.engine.Run(url); err != nil {
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

	fakeJWT := gojwt.NewWithClaims(gojwt.SigningMethodRS256, &jwt.Claims{
		Subject:           ticketTicketApp.config.Dev.User.UserID,
		Email:             ticketTicketApp.config.Dev.User.Email,
		PreferredUsername: ticketTicketApp.config.Dev.User.Username,
	})

	signedJWT, err := fakeJWT.SignedString(rsaPrivKey)

	if err != nil {
		panic(fmt.Errorf("error generation fake Token: %s", err.Error()))
	}

	fmt.Printf("DEV Token: %s \n", signedJWT)
}
