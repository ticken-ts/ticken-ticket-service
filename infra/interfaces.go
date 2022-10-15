package infra

import (
	"github.com/gin-gonic/gin"
	pvtbc "github.com/ticken-ts/ticken-pvtbc-connector"
)

type Db interface {
	Connect(connString string) error
	IsConnected() bool

	// GetClient is going to return the raw client.
	// The caller should cast the returned value
	// into the correct client depending on the
	// driver
	GetClient() interface{}
}

type IBuilder interface {
	BuildDb(connString string) Db
	BuildEngine() *gin.Engine
	BuildPvtbcCaller() *pvtbc.Caller
	BuildPvtbcListener() *pvtbc.Listener
}
