package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"ticken-ticket-service/controllers/ticketController"
	"ticken-ticket-service/db"
	"ticken-ticket-service/utils"
)

func main() {
	router := gin.Default()

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	mongoDB := db.ConnectMongoDB(utils.GetEnvOrPanic("MONGO_URI"))

	ticketController.RegisterRoutes(router, mongoDB)

	err = router.Run("localhost:8080")
	if err != nil {
		panic(err)
	}
}
