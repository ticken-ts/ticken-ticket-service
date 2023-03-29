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

type linkTicketPayload struct {
	EventContractAddr string `json:"event_contract_addr"`
}

func (controller *TicketController) LinkTicket(c *gin.Context) {
	var payload linkTicketPayload

	attendantID := c.MustGet("jwt").(*jwt.Token).Subject

	if err := c.BindJSON(&payload); err != nil {
		c.Abort()
		return
	}

	linkedTickets, err := controller.serviceProvider.GetTicketLinker().LinkTickets(
		attendantID,
		payload.EventContractAddr,
	)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	// Create array of dto.Ticket
	ticketsDTO := make([]*dto.Ticket, len(linkedTickets))

	// Map each ticket to dto.Ticket
	for i, ticket := range linkedTickets {
		ticketsDTO[i] = mappers.MapTicketToDTO(ticket)
	}

	c.JSON(http.StatusOK, res.Success{
		Message: fmt.Sprintf("%d tickets linked", len(linkedTickets)),
		Data:    ticketsDTO,
	})
}
