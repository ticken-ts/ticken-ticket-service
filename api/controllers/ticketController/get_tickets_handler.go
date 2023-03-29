package ticketController

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"ticken-ticket-service/api/dto"
	"ticken-ticket-service/api/mappers"
	"ticken-ticket-service/api/res"
	"ticken-ticket-service/security/jwt"
)

func (controller *TicketController) GetMyTickets(c *gin.Context) {
	attendantID := c.MustGet("jwt").(*jwt.Token).Subject

	tickets, err := controller.serviceProvider.GetTicketIssuer().GetUserTickets(
		attendantID,
	)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	// Create array of dto.Ticket
	ticketsDTO := make([]*dto.Ticket, len(tickets))

	// Map each ticket to dto.Ticket
	for i, ticket := range tickets {
		ticketsDTO[i] = mappers.MapTicketToDTO(ticket)
	}

	c.JSON(http.StatusOK, res.Success{
		Message: fmt.Sprintf("%d ticket's found", len(tickets)),
		Data:    ticketsDTO,
	})
}
