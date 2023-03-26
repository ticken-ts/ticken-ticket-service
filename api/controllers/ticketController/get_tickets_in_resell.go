package ticketController

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"ticken-ticket-service/api/dto"
	"ticken-ticket-service/api/mappers"
	"ticken-ticket-service/utils"
)

func (controller *TicketController) GetTicketsInResell(c *gin.Context) {
	eventID, err := uuid.Parse(c.Param("eventID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.HttpResponse{Message: err.Error()})
		c.Abort()
		return
	}

	section := c.Query("section")

	tickets, err := controller.serviceProvider.GetTicketTrader().GetResells(eventID, section)
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
