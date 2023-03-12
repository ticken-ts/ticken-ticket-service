package userController

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"ticken-ticket-service/api/mappers"
	"ticken-ticket-service/security/jwt"
	"ticken-ticket-service/utils"
)

type createAccountPayload struct {
	AddressPK string `json:"addressPK"`
}

func (controller *UserController) CreateAccount(c *gin.Context) {
	owner := c.MustGet("jwt").(*jwt.Token)
	token := owner.Subject
	email := owner.Email
	profile := owner.Profile

	payload := createAccountPayload{}
	err := c.BindJSON(&payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, utils.HttpResponse{Message: err.Error()})
		c.Abort()
		return
	}

	user, err := controller.serviceProvider.GetUserManager().CreateUser(token, payload.AddressPK)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.HttpResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, utils.HttpResponse{Message: "Account created successfully", Data: mappers.MapUserToDTO(user, email, &profile)})
}
