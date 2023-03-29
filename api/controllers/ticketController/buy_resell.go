package ticketController

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"ticken-ticket-service/api/mappers"
	"ticken-ticket-service/api/res"
	"ticken-ticket-service/security/jwt"
)

func (controller *TicketController) BuyResell(c *gin.Context) {
	var payload sellTicketPayload

	buyerID := c.MustGet("jwt").(*jwt.Token).Subject

	if err := c.BindJSON(&payload); err != nil {
		c.Abort()
		return
	}

	eventID, err := uuid.Parse(c.Param("eventID"))
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	ticketID, err := uuid.Parse(c.Param("ticketID"))
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	resellID, err := uuid.Parse(c.Param("resellID"))
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	ticket, err := controller.serviceProvider.GetTicketTrader().BuyResoldTicket(
		buyerID,
		eventID,
		ticketID,
		resellID,
	)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, res.Success{
		Message: "Ticket bought successfully",
		Data:    mappers.MapTicketToDTO(ticket),
	})
}
