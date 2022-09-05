package ticketController

import (
	"ticken-ticket-service/infra"
	"ticken-ticket-service/services"
)

type ticketController struct {
	serviceProvider services.Provider
}

func NewTicketController(serviceProvider services.Provider) *ticketController {
	controller := new(ticketController)
	controller.serviceProvider = serviceProvider
	return controller
}

func (controller *ticketController) Setup(router infra.Router) {
	router.POST("/tickets", controller.BuyTicket)
}
