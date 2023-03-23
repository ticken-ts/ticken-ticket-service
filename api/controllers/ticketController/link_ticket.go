package ticketController

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"ticken-ticket-service/api/dto"
	"ticken-ticket-service/api/mappers"
	"ticken-ticket-service/security/jwt"
	"ticken-ticket-service/utils"
)

type linkTicketPayload struct {
	EventContractAddr string `json:"event_contract_addr"`
}

func (controller *TicketController) LinkTicket(c *gin.Context) {
	var payload linkTicketPayload

	ownerID := c.MustGet("jwt").(*jwt.Token).Subject

	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, utils.HttpResponse{Message: err.Error()})
		c.Abort()
		return
	}

	linkedTickets, err := controller.serviceProvider.GetTicketLinker().LinkTickets(ownerID, payload.EventContractAddr)

	if err != nil {
		c.JSON(http.StatusBadRequest, utils.HttpResponse{Message: err.Error()})
		c.Abort()
		return
	}

	// Create array of dto.Ticket
	ticketsDTO := make([]*dto.Ticket, len(linkedTickets))

	// Map each ticket to dto.Ticket
	for i, ticket := range linkedTickets {
		ticketsDTO[i] = mappers.MapTicketToDTO(ticket)
	}

	c.JSON(http.StatusOK, utils.HttpResponse{
		Message: fmt.Sprintf("%d tickets linked", len(linkedTickets)),
		Data:    ticketsDTO,
	})
}
