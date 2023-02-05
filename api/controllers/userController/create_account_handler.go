package userController

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"ticken-ticket-service/security/jwt"
	"ticken-ticket-service/utils"
)

func (controller *UserController) CreateAccount(c *gin.Context) {
	owner := c.MustGet("jwt").(*jwt.Token).Subject

	c.JSON(http.StatusOK, utils.HttpResponse{Message: "Account created successfully", Data: owner})
}
