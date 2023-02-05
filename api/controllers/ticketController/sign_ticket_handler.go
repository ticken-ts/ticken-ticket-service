package ticketController

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"ticken-ticket-service/api/mappers"
	"ticken-ticket-service/security/jwt"
	"ticken-ticket-service/utils"
)

func (controller *TicketController) SignTicket(c *gin.Context) {
	eventID, ticketID := c.Param("eventID"), c.Param("ticketID")
	owner := c.MustGet("jwt").(*jwt.Token).Subject

	ticketSigner := controller.serviceProvider.GetTicketSigner()

	signedTicket, err := ticketSigner.SignTicket(eventID, ticketID, owner.String())
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.HttpResponse{Message: err.Error()})
		c.Abort()
		return
	}

	ticketDTO := mappers.MapTicketToDTO(signedTicket)

	c.JSON(http.StatusOK, utils.HttpResponse{Data: ticketDTO})
}
