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

func New(serviceProvider services.IProvider) *TicketController {
	controller := new(TicketController)
	controller.validator = validator.New()
	controller.serviceProvider = serviceProvider
	return controller
}

func (controller *TicketController) Setup(router gin.IRouter) {
	router.POST("/events/:eventID/tickets/:ticketID/resells/:resellID", controller.BuyResell)
	router.PUT("/events/:eventID/tickets/:ticketID/resells", controller.ResellTicket)
	router.POST("/events/:eventID/tickets", controller.BuyTicket)
	router.GET("/events/tickets", controller.GetMyTickets)
	router.POST("/tickets/link", controller.LinkTicket)
}
