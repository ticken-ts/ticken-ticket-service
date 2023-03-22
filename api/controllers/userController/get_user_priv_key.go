package userController

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"ticken-ticket-service/security/jwt"
	"ticken-ticket-service/utils"
)

func (controller *UserController) GetUserPrivKey(c *gin.Context) {
	token := c.MustGet("jwt").(*jwt.Token)
	attendantID := token.Subject

	userPrivKey, err := controller.serviceProvider.GetUserManager().GetUserPrivKey(attendantID)
	if err != nil {
		c.JSON(http.StatusNotFound, utils.HttpResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, utils.HttpResponse{Data: userPrivKey})
}
