package ticketController

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"net/http"
	"ticken-ticket-service/utils"
)

type buyTicketPayload struct {
	EventID string `json:"event_id"`
	Section string `json:"section"`
}

var validate = validator.New()
var owner = uuid.New().String()

func (controller *ticketController) BuyTicket(c *gin.Context) {
	var payload buyTicketPayload

	err := c.BindJSON(&payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.HttpResponse{Message: err.Error()})
		c.Abort()
		return
	}

	err = validate.Struct(&payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.HttpResponse{Message: err.Error()})
		c.Abort()
		return
	}

	ticketIssuer := controller.serviceProvider.GetTicketIssuer()

	newTicket, err := ticketIssuer.IssueTicket(payload.EventID, payload.Section, owner)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.HttpResponse{Message: err.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, utils.HttpResponse{Data: newTicket})
}
