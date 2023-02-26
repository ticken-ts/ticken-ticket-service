package ticketController

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"ticken-ticket-service/api/mappers"
	"ticken-ticket-service/security/jwt"
	"ticken-ticket-service/utils"
)

type buyTicketPayload struct {
	Section string `json:"section"`
}

func (controller *TicketController) BuyTicket(c *gin.Context) {
	var payload buyTicketPayload

	ownerID := c.MustGet("jwt").(*jwt.Token).Subject

	eventID, err := uuid.Parse(c.Param("eventID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.HttpResponse{Message: err.Error()})
		c.Abort()
		return
	}

	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, utils.HttpResponse{Message: err.Error()})
		c.Abort()
		return
	}

	if err := controller.validator.Struct(&payload); err != nil {
		c.JSON(http.StatusBadRequest, utils.HttpResponse{Message: err.Error()})
		c.Abort()
		return
	}

	ticketIssuer := controller.serviceProvider.GetTicketIssuer()

	newTicket, err := ticketIssuer.IssueTicket(eventID, payload.Section, ownerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.HttpResponse{Message: err.Error()})
		c.Abort()
		return
	}

	ticketDTO := mappers.MapTicketToDTO(newTicket)

	c.JSON(http.StatusCreated, utils.HttpResponse{Data: ticketDTO})
}
