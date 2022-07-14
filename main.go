package main

import (
	"github.com/gin-gonic/gin"
	"ticken-ticket-service/config"
	"ticken-ticket-service/routes"
)

func main() {
	router := gin.Default()

	config.ConnectDB()

	//routes
	routes.TicketRoute(router)

	router.Run("localhost:8080")
}
