package healthController

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"ticken-ticket-service/api/res"
	"ticken-ticket-service/services"
)

const HealthMessage = "Everything is fine"

type HealthController struct {
	serviceProvider services.IProvider
}

func New(serviceProvider services.IProvider) *HealthController {
	controller := new(HealthController)
	controller.serviceProvider = serviceProvider
	return controller
}

func (controller *HealthController) Setup(router gin.IRouter) {
	router.GET("/healthz", controller.Healthz)
}

func (controller *HealthController) Healthz(c *gin.Context) {
	c.JSON(http.StatusOK, res.Success{Message: HealthMessage})
}
