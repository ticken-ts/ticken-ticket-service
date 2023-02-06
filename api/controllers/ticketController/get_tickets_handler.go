package ticketController

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"ticken-ticket-service/api/dto"
	"ticken-ticket-service/api/mappers"
	"ticken-ticket-service/security/jwt"
	"ticken-ticket-service/utils"
)

func (controller *TicketController) GetMyTickets(c *gin.Context) {
	ownerID := c.MustGet("jwt").(*jwt.Token).Subject

	tickets, err := controller.serviceProvider.GetTicketIssuer().GetUserTickets(ownerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.HttpResponse{Message: err.Error()})
		c.Abort()
		return
	}

	// Create array of dto.Ticket
	ticketsDTO := make([]*dto.Ticket, len(tickets))

	// Map each ticket to dto.Ticket
	for i, ticket := range tickets {
		ticketsDTO[i] = mappers.MapTicketToDTO(ticket)
	}

	c.JSON(http.StatusOK, utils.HttpResponse{Data: ticketsDTO})
}
