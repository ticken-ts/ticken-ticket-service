package userController

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"ticken-ticket-service/api/mappers"
	"ticken-ticket-service/security/jwt"
	"ticken-ticket-service/utils"
)

func (controller *UserController) GetUser(c *gin.Context) {
	owner := c.MustGet("jwt").(*jwt.Token)
	token := owner.Subject
	email := owner.Email
	profile := owner.Profile

	user, err := controller.serviceProvider.GetUserManager().GetUser(token)

	if err != nil {
		c.JSON(http.StatusNotFound, utils.HttpResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, utils.HttpResponse{Message: "User fetched successfully", Data: mappers.MapUserToDTO(user, email, &profile)})
}
