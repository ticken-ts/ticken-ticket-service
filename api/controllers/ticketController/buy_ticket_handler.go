package ticketController

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"ticken-ticket-service/api/mappers"
	"ticken-ticket-service/api/res"
	"ticken-ticket-service/security/jwt"
)

type buyTicketPayload struct {
	Section string `json:"section"`
}

func (controller *TicketController) BuyTicket(c *gin.Context) {
	var payload buyTicketPayload

	ownerID := c.MustGet("jwt").(*jwt.Token).Subject

	eventID, err := uuid.Parse(c.Param("eventID"))
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	if err := c.BindJSON(&payload); err != nil {
		c.Abort()
		return
	}

	newTicket, err := controller.serviceProvider.GetTicketIssuer().IssueTicket(
		eventID,
		payload.Section,
		ownerID,
	)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, res.Success{
		Message: "ticket bought successfully",
		Data:    mappers.MapTicketToDTO(newTicket),
	})
}
