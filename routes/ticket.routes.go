package routes

import (
	"github.com/gin-gonic/gin"
	"ticken-ticket-service/controllers"
)

func TicketRoute(router *gin.Engine) {
	router.POST("/tickets", controllers.BuyTicket())
	router.GET("/tickets/:id", controllers.GetTicket())
}
