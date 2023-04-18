package userController

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"ticken-ticket-service/api/res"
	"ticken-ticket-service/security/jwt"
)

func (controller *UserController) GetUserPrivKey(c *gin.Context) {
	token := c.MustGet("jwt").(*jwt.Token)
	attendantID := token.Subject

	format := c.DefaultQuery("format", "pem")

	userPrivKey, err := controller.serviceProvider.GetUserManager().GetUserPrivKey(
		attendantID,
		format,
	)
	if err != nil {
		c.Error(err)
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, res.Success{
		Message: fmt.Sprintf("private key obtainer successfully in format %s", format),
		Data:    userPrivKey,
	})
}
