package controllers

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"ticken-ticket-service/config"
	"ticken-ticket-service/models"
	"ticken-ticket-service/responses"
	"time"
)

var ticketCollection = config.GetCollection(config.DB, "tickets")
var validate = validator.New()

var owner = uuid.New()

func BuyTicket() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var ticket models.Ticket
		defer cancel()

		//validate the request body
		if err := c.BindJSON(&ticket); err != nil {
			c.JSON(http.StatusBadRequest,
				responses.GenericResponse{
					Status:  http.StatusBadRequest,
					Message: "error-1",
					Data:    map[string]interface{}{"data": err.Error()},
				})
			return
		}

		//use the validator library to validate required fields
		if validationErr := validate.Struct(&ticket); validationErr != nil {
			c.JSON(http.StatusBadRequest,
				responses.GenericResponse{
					Status:  http.StatusBadRequest,
					Message: "error",
					Data:    map[string]interface{}{"data": validationErr.Error()},
				})
			return
		}

		newTicket := models.Ticket{
			ID:      primitive.NewObjectID(),
			Owner:   owner.String(),
			EventID: ticket.EventID,
			Section: ticket.Section,
		}

		result, err := ticketCollection.InsertOne(ctx, newTicket)
		if err != nil {
			c.JSON(http.StatusInternalServerError,
				responses.GenericResponse{
					Status:  http.StatusInternalServerError,
					Message: "error",
					Data:    map[string]interface{}{"data": err.Error()},
				})
			return
		}

		c.JSON(http.StatusCreated,
			responses.GenericResponse{
				Status:  http.StatusCreated,
				Message: "success",
				Data:    map[string]interface{}{"data": result},
			})
	}
}

func GetTicket() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		id := c.Param("id")
		var event models.Event
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(id)
		println("ID: ", objId.String())

		err := ticketCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&event)
		if err != nil {
			c.JSON(http.StatusInternalServerError,
				responses.GenericResponse{
					Status:  http.StatusInternalServerError,
					Message: "error",
					Data:    map[string]interface{}{"data": err.Error()},
				})
			return
		}

		c.JSON(http.StatusOK,
			responses.GenericResponse{
				Status:  http.StatusOK,
				Message: "success",
				Data:    map[string]interface{}{"data": event},
			})
	}
}
