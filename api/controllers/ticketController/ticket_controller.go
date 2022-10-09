package ticketController

import (
	"github.com/go-playground/validator/v10"
	"ticken-ticket-service/infra"
	"ticken-ticket-service/services"
)

type TicketController struct {
	validator       *validator.Validate
	serviceProvider services.Provider
}

func NewTicketController(serviceProvider services.Provider) *TicketController {
	controller := new(TicketController)
	controller.validator = validator.New()
	controller.serviceProvider = serviceProvider
	return controller
}

func (controller *TicketController) Setup(router infra.Router) {
	router.POST("/events/:eventID/tickets", controller.BuyTicket)
	router.PUT("/events/:eventID/tickets/:ticketID/sign", controller.SignTicket) // <- Es REST LCTM
}
