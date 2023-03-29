package ticketController

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"ticken-ticket-service/api/dto"
	"ticken-ticket-service/api/mappers"
	"ticken-ticket-service/api/res"
)

func (controller *TicketController) GetTicketsInResell(c *gin.Context) {
	eventID, err := uuid.Parse(c.Param("eventID"))
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	section := c.Query("section")

	tickets, err := controller.serviceProvider.GetTicketTrader().GetTicketsInResells(
		eventID,
		section,
	)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	ticketsDTO := make([]*dto.Ticket, len(tickets))

	// Map each ticket to dto.Ticket
	for i, ticket := range tickets {
		ticketsDTO[i] = mappers.MapTicketToDTO(ticket)
	}

	c.JSON(http.StatusOK, res.Success{
		Message: fmt.Sprintf("%d tickets in resell found", len(tickets)),
		Data:    ticketsDTO,
	})
}
