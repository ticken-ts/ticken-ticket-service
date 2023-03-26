package userController

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"ticken-ticket-service/api/mappers"
	"ticken-ticket-service/api/res"
	"ticken-ticket-service/security/jwt"
)

type createAccountPayload struct {
	AddressPK string `json:"addressPK"`
}

func (controller *UserController) CreateAccount(c *gin.Context) {
	attendant := c.MustGet("jwt").(*jwt.Token)

	token := attendant.Subject
	email := attendant.Email
	profile := attendant.Profile

	payload := createAccountPayload{}
	if err := c.BindJSON(&payload); err != nil {
		c.Abort()
		return
	}

	user, err := controller.serviceProvider.GetUserManager().CreateUser(token, payload.AddressPK)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, res.Success{
		Message: "Account created successfully",
		Data:    mappers.MapUserToDTO(user, email, &profile),
	})
}
