package userController

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"ticken-ticket-service/services"
)

type UserController struct {
	validator       *validator.Validate
	serviceProvider services.IProvider
}

func New(serviceProvider services.IProvider) *UserController {
	controller := new(UserController)
	controller.validator = validator.New()
	controller.serviceProvider = serviceProvider
	return controller
}

func (controller *UserController) Setup(router gin.IRouter) {
	router.POST("/users", controller.CreateAccount)
	router.GET("/users/myUser", controller.GetUser)
	router.GET("/users/myUser/privKey", controller.GetUserPrivKey)
}
