package ticketController

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"ticken-ticket-service/services"
)

type TicketController struct {
	validator       *validator.Validate
	serviceProvider services.IProvider
}

func NewTicketController(serviceProvider services.IProvider) *TicketController {
	controller := new(TicketController)
	controller.validator = validator.New()
	controller.serviceProvider = serviceProvider
	return controller
}

func (controller *TicketController) Setup(router gin.IRouter) {
	router.POST("/events/:eventID/tickets", controller.BuyTicket)
	router.PUT("/events/:eventID/tickets/:ticketID/sign", controller.SignTicket)
}
