package api

import (
	"github.com/gin-gonic/gin"
)

type Controller interface {
	Setup(router gin.IRouter)
}

type Middleware interface {
	Setup(router gin.IRouter)
}
