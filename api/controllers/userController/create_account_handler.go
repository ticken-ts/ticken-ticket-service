package userController

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"ticken-ticket-service/security/jwt"
	"ticken-ticket-service/utils"
)

func (controller *UserController) CreateAccount(c *gin.Context) {
	owner := c.MustGet("jwt").(*jwt.Token).Subject
	user, err := controller.serviceProvider.GetUserManager().CreateUser(owner)
	if err != nil {
		c.JSON(http.StatusInternalServerError, utils.HttpResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, utils.HttpResponse{Message: "Account created successfully", Data: user})
}
