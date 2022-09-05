package infra

import "github.com/gin-gonic/gin"

type Db interface {
	Connect(connString string) error
	IsConnected() bool

	// GetClient is going to return the raw client.
	// The caller should cast the returned value
	// into the correct client depending on the
	// driver
	GetClient() interface{}
}

type Router interface {
	gin.IRouter
	Run(addresses ...string) (err error)
}
