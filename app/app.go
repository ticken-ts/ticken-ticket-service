package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"ticken-ticket-service/api"
	"ticken-ticket-service/api/controllers/ticketController"
	"ticken-ticket-service/api/middlewares"
	"ticken-ticket-service/config"
	"ticken-ticket-service/env"
	"ticken-ticket-service/infra"
	"ticken-ticket-service/services"
)

type TickenTicketApp struct {
	engine          *gin.Engine
	config          *config.Config
	serviceProvider services.Provider
}

func New(builder *infra.Builder, tickenConfig *config.Config) *TickenTicketApp {
	tickenTicketApp := new(TickenTicketApp)

	router := builder.BuildEngine()
	pvtbcCaller := builder.BuildPvtbcCaller()
	db := builder.BuildDb(env.TickenEnv.ConnString)

	// this provider is going to provide all services
	// needed by the controllers to execute it operations
	serviceProvider, err := services.NewProvider(db, pvtbcCaller, tickenConfig)
	if err != nil {
		panic(err)
	}

	tickenTicketApp.engine = router
	tickenTicketApp.serviceProvider = serviceProvider

	var appMiddlewares = []api.Middleware{
		middlewares.NewAuthMiddleware(serviceProvider, &tickenConfig.Server),
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
	url := ticketTicketApp.config.Server.GetServerURL()
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
		"sub":   "290c641a-55a1-40f5-acc3-d4ebe3626fdd",
		"email": "user@ticken.com",
	})

	signedJWT, err := fakeJWT.SigningString()
	if err != nil {
		panic(fmt.Errorf("error generation fake JWT: %s", err.Error()))
	}

	fmt.Printf("DEV JWT: %s \n", signedJWT)
}
