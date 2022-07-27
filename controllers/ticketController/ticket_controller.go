package ticketController

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"ticken-ticket-service/blockchain/tickenPVTBCConnector"
	"ticken-ticket-service/db/repositories"
	"ticken-ticket-service/services"
)

type ticketController struct {
	ticketIssuer services.TicketIssuer
}

func RegisterRoutes(router *gin.Engine, db *mongo.Client) {
	controller := new(ticketController)

	controller.ticketIssuer = services.NewTicketIssuer(
		repositories.NewEventMongoDBRepository(db, "ticken-db"),
		repositories.NewTicketMongoDBRepository(db, "ticken-db"),
		tickenPVTBCConnector.New(),
	)

	router.POST("/tickets", controller.BuyTicket)
}
