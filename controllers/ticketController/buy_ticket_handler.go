package ticketController

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"net/http"
	"ticken-ticket-service/models/ticket"
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

	var newTicket = ticket.New(payload.EventID, payload.Section)
	_ = newTicket.AssignTo(owner) // new ticket never has owner

	c.JSON(http.StatusCreated, utils.HttpResponse{Data: newTicket})
}
