package main

import (
	"github.com/gin-gonic/gin"
	"ticken-ticket-service/controllers/ticketController"
)

func main() {
	router := gin.Default()

	ticketController.RegisterRoutes(router)

	err := router.Run("localhost:8080")
	if err != nil {
		panic(err)
	}
}
