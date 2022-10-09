package app

import (
	"ticken-ticket-service/api"
	"ticken-ticket-service/api/controllers/ticketController"
	"ticken-ticket-service/api/middlewares"
	"ticken-ticket-service/infra"
	"ticken-ticket-service/services"
	"ticken-ticket-service/utils"
)

type TickenTicketApp struct {
	router          infra.Router
	serviceProvider services.Provider
}

func New(builder *infra.Builder, tickenConfig *utils.TickenConfig) *TickenTicketApp {
	tickenTicketApp := new(TickenTicketApp)

	db := builder.BuildDb()
	router := builder.BuildRouter()
	pvtbcCaller := builder.BuildPvtbcCaller()

	// this provider is going to provide all services
	// needed by the controllers to execute it operations
	serviceProvider, err := services.NewProvider(db, pvtbcCaller, tickenConfig)
	if err != nil {
		panic(err)
	}

	tickenTicketApp.router = router
	tickenTicketApp.serviceProvider = serviceProvider

	var appMiddlewares = []api.Middleware{
		middlewares.NewAuthMiddleware(serviceProvider),
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

func (tickenTicketApp *TickenTicketApp) Start() {
	err := tickenTicketApp.router.Run("localhost:9000")
	if err != nil {
		panic(err)
	}
}

func (tickenTicketApp *TickenTicketApp) Populate() {
	eventManager := tickenTicketApp.serviceProvider.GetEventManager()
	_, err := eventManager.AddEvent("test-event-id", "organizer", "ticken-test-channel")
	if err != nil {
		return // HANDLER DUPLICATES
	}
}
