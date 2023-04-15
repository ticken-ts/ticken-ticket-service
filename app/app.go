package app

import (
	"fmt"
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

	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	gojwt "github.com/golang-jwt/jwt"
)

type TickenTicketApp struct {
	engine          *gin.Engine
	config          *config.Config
	repoProvider    repos.IProvider
	serviceProvider services.IProvider
	jwtVerifier     jwt.Verifier
	subscriber      *async.Subscriber

	// populators are intended to populate
	// useful data. It can be testdata or
	// data that should be present on the db
	// before the service is available
	populators []Populator
}

func New(infraBuilder infra.IBuilder, tickenConfig *config.Config) *TickenTicketApp {
	log.TickenLogger.Info().Msg(color.BlueString("initializing " + tickenConfig.Server.ClientID))

	tickenTicketApp := new(TickenTicketApp)

	/******************************** infra builds ********************************/
	engine := infraBuilder.BuildEngine()
	pvtbcCaller := infraBuilder.BuildPvtbcCaller()
	jwtVerifier := infraBuilder.BuildJWTVerifier()
	db := infraBuilder.BuildDb(env.TickenEnv.DbConnString)
	hsm := infraBuilder.BuildHSM(env.TickenEnv.HSMEncryptionKey)
	pubbcAdmin := infraBuilder.BuildPubbcAdmin(env.TickenEnv.TickenWalletKey)
	pubbcCaller := infraBuilder.BuildPubbcCaller(env.TickenEnv.TickenWalletKey)
	busSubscriber := infraBuilder.BuildBusSubscriber(env.TickenEnv.BusConnString)
	/**************************++***************************************************/

	/********************************** providers **********************************/
	repoProvider, err := repos.NewProvider(
		db,
		&tickenConfig.Database,
	)
	if err != nil {
		log.TickenLogger.Panic().Msg(err.Error())
	}

	serviceProvider, err := services.NewProvider(
		repoProvider,
		pvtbcCaller,
		pubbcAdmin,
		pubbcCaller,
		hsm,
	)
	if err != nil {
		log.TickenLogger.Panic().Msg(err.Error())
	}
	/**************************++***************************************************/

	/********************************* subscriber **********************************/
	subscriber, err := async.NewSubscriber(busSubscriber, serviceProvider)
	if err != nil {
		log.TickenLogger.Panic().Msg(err.Error())
	}
	/**************************++***************************************************/

	tickenTicketApp.engine = engine
	tickenTicketApp.config = tickenConfig
	tickenTicketApp.subscriber = subscriber
	tickenTicketApp.jwtVerifier = jwtVerifier
	tickenTicketApp.repoProvider = repoProvider
	tickenTicketApp.serviceProvider = serviceProvider

	tickenTicketApp.loadMiddlewares(engine)
	tickenTicketApp.loadControllers(engine)

	/********************************* populators **********************************/
	tickenTicketApp.populators = []Populator{
		&fakes.FakeUsersPopulator{Repo: repoProvider.GetUserRepository(), Config: tickenConfig.Dev.User, HSM: hsm},
	}
	/**************************++***************************************************/

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

func (ticketTicketApp *TickenTicketApp) loadControllers(apiRouter gin.IRouter) {
	apiRouterGroup := apiRouter.Group(ticketTicketApp.config.Server.APIPrefix)
	var appControllers = []api.Controller{
		userController.New(ticketTicketApp.serviceProvider),
		healthController.New(ticketTicketApp.serviceProvider),
		ticketController.New(ticketTicketApp.serviceProvider),
	}

	for _, controller := range appControllers {
		controller.Setup(apiRouterGroup)
	}
}

func (ticketTicketApp *TickenTicketApp) loadMiddlewares(apiRouter gin.IRouter) {
	var appMiddlewares = []api.Middleware{
		middlewares.NewErrorMiddleware(),
		middlewares.NewAuthMiddleware(ticketTicketApp.jwtVerifier, ticketTicketApp.config.Server.APIPrefix),
	}

	for _, middleware := range appMiddlewares {
		middleware.Setup(apiRouter)
	}
}
