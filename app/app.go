package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	pvtbc "github.com/ticken-ts/ticken-pvtbc-connector"
	"ticken-ticket-service/api"
	"ticken-ticket-service/api/controllers/ticketController"
	"ticken-ticket-service/api/middlewares"
	"ticken-ticket-service/infra"
	"ticken-ticket-service/services"
	"ticken-ticket-service/utils"
)

type TickenTicketApp struct {
	engine          *gin.Engine
	serviceProvider services.Provider
	config          *utils.TickenConfig
}

func New(builder *infra.Builder, tickenConfig *utils.TickenConfig) *TickenTicketApp {
	tickenTicketApp := new(TickenTicketApp)

	db := builder.BuildDb()
	router := builder.BuildEngine()
	pvtbcCaller := new(pvtbc.Caller)

	// this provider is going to provide all services
	// needed by the controllers to execute it operations
	serviceProvider, err := services.NewProvider(db, pvtbcCaller, tickenConfig)
	if err != nil {
		panic(err)
	}

	tickenTicketApp.engine = router
	tickenTicketApp.serviceProvider = serviceProvider

	var appMiddlewares = []api.Middleware{
		middlewares.NewAuthMiddleware(serviceProvider, tickenConfig),
	}

	for _, middleware := range appMiddlewares {
		middleware.Setup(router)
	}

	var appControllers = []api.Controller{
		ticketController.NewTicketController(serviceProvider),
	}

	for _, controller := range appControllers {
		controller.Setup(router)
	}

	return tickenTicketApp
}

func (ticketTicketApp *TickenTicketApp) Start() {
	url := getServerURL(&ticketTicketApp.config.Config.Server)
	err := ticketTicketApp.engine.Run(url)
	if err != nil {
		panic(err)
	}
}

func (ticketTicketApp *TickenTicketApp) Populate() {
	eventManager := ticketTicketApp.serviceProvider.GetEventManager()
	_, err := eventManager.AddEvent("test-event-id", "organizer", "ticken-test-channel")
	if err != nil {
		return // HANDLER DUPLICATES
	}
}

func (ticketTicketApp *TickenTicketApp) EmitFakeJWT() {
	fakeJWT := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"sub": "290c641a-55a1-40f5-acc3-d4ebe3626fdd",
	})

	signedJWT, err := fakeJWT.SigningString()
	if err != nil {
		panic(fmt.Errorf("error generation fake JWT: %s", err.Error()))
	}

	fmt.Printf("DEV JWT: %s \n", signedJWT)
}

func getServerURL(serverConfig *utils.ServerConfig) string {
	return serverConfig.Host + ":" + serverConfig.Port
}
