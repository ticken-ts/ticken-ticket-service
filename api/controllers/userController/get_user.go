package userController

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"ticken-ticket-service/api/mappers"
	"ticken-ticket-service/api/res"
	"ticken-ticket-service/security/jwt"
)

func (controller *UserController) GetUser(c *gin.Context) {
	token := c.MustGet("jwt").(*jwt.Token)

	attendantID := token.Subject
	email := token.Email
	profile := token.Profile

	user, err := controller.serviceProvider.GetUserManager().GetUser(
		attendantID,
	)

	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, res.Success{
		Message: "User fetched successfully",
		Data:    mappers.MapUserToDTO(user, email, &profile),
	})
}
