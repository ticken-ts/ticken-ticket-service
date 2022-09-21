package ticketController

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"ticken-ticket-service/utils"
)

func (controller *TicketController) SignTicket(c *gin.Context) {
	eventID, ticketID := c.Param("eventID"), c.Param("ticketID")

	ticketSigner := controller.serviceProvider.GetTicketSigner()

	signedTicket, err := ticketSigner.SignTicket(eventID, ticketID, owner)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.HttpResponse{Message: err.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, utils.HttpResponse{Data: signedTicket})
}
