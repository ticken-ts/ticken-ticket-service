package ticketController

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"ticken-ticket-service/api/mappers"
	"ticken-ticket-service/security/jwt"
	"ticken-ticket-service/utils"
	"ticken-ticket-service/utils/money"
)

type sellTicketPayload struct {
	Price    float64 `json:"price"`
	Currency string  `json:"currency"`
}

func (controller *TicketController) ResellTicket(c *gin.Context) {
	var payload sellTicketPayload

	ownerID := c.MustGet("jwt").(*jwt.Token).Subject

	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, utils.HttpResponse{Message: err.Error()})
		c.Abort()
		return
	}

	eventID, err := uuid.Parse(c.Param("eventID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.HttpResponse{Message: err.Error()})
		c.Abort()
		return
	}

	ticketID, err := uuid.Parse(c.Param("ticketID"))
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.HttpResponse{Message: err.Error()})
		c.Abort()
		return
	}

	price, err := money.BuildFrom(payload.Price, payload.Currency)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.HttpResponse{Message: err.Error()})
		c.Abort()
		return
	}

	ticket, err := controller.serviceProvider.GetTicketTrader().SellTicket(
		ownerID,
		eventID,
		ticketID,
		price,
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.HttpResponse{Message: err.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, utils.HttpResponse{
		Message: "Ticket published for sale successfully",
		Data:    mappers.MapTicketToDTO(ticket),
	})
}
