package app

import (
	"ticken-ticket-service/api"
	"ticken-ticket-service/api/controllers/ticketController"
	"ticken-ticket-service/infra"
	"ticken-ticket-service/services"
	"ticken-ticket-service/utils"
)

type tickenTicketApp struct {
	router          infra.Router
	serviceProvider services.Provider
}

func New(router infra.Router, db infra.Db, tickenConfig *utils.TickenConfig) *tickenTicketApp {
	tickenTicketApp := new(tickenTicketApp)

	// this provider is going to provide all services
	// needed by the controllers to execute it operations
	serviceProvider, _ := services.NewProvider(db, tickenConfig)

	tickenTicketApp.router = router
	tickenTicketApp.serviceProvider = serviceProvider

	var controllers = []api.Controller{
		ticketController.NewTicketController(serviceProvider),
	}

	for _, controller := range controllers {
		controller.Setup(router)
	}

	return tickenTicketApp
}

func (tickenTicketApp *tickenTicketApp) Start() {
	err := tickenTicketApp.router.Run("localhost:8080")
	if err != nil {
		panic(err)
	}
}
