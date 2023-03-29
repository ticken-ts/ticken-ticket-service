package userController

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"ticken-ticket-service/api/res"
	"ticken-ticket-service/security/jwt"
)

func (controller *UserController) GetUserPrivKey(c *gin.Context) {
	token := c.MustGet("jwt").(*jwt.Token)
	attendantID := token.Subject

	userPrivKey, err := controller.serviceProvider.GetUserManager().GetUserPrivKey(
		attendantID,
	)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, res.Success{
		Data: userPrivKey,
	})
}
