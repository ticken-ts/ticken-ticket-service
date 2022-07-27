package ticketController

import (
	"github.com/gin-gonic/gin"
)

type ticketController struct {
}

func RegisterRoutes(router *gin.Engine) {
	controller := new(ticketController)

	router.POST("/tickets", controller.BuyTicket)
}
