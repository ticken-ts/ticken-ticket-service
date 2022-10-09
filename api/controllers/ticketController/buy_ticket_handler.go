package ticketController

import (
	"github.com/coreos/go-oidc"
	"github.com/gin-gonic/gin"
	"net/http"
	"ticken-ticket-service/api/mappers"
	"ticken-ticket-service/utils"
)

type buyTicketPayload struct {
	Section string `json:"section"`
}

func (controller *TicketController) BuyTicket(c *gin.Context) {
	var payload buyTicketPayload
	eventID := c.Param("eventID")
	owner := c.MustGet("jwt").(*oidc.IDToken).Subject

	err := c.BindJSON(&payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.HttpResponse{Message: err.Error()})
		c.Abort()
		return
	}

	err = controller.validator.Struct(&payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.HttpResponse{Message: err.Error()})
		c.Abort()
		return
	}

	ticketIssuer := controller.serviceProvider.GetTicketIssuer()

	newTicket, err := ticketIssuer.IssueTicket(eventID, payload.Section, owner)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.HttpResponse{Message: err.Error()})
		c.Abort()
		return
	}

	ticketDTO := mappers.MapTicketToDTO(newTicket)

	c.JSON(http.StatusCreated, utils.HttpResponse{Data: ticketDTO})
}
