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

type BusSubscriber interface {
	Connect(connString string, exchangeName string) error
	IsConnected() bool
	Listen(handler func([]byte)) error
}

type BusPublisher interface {
	Connect(connString string, qName string) error
	IsConnected() bool
	Listen(handler func([]byte)) error
}

type IBuilder interface {
	BuildEngine() *gin.Engine
	BuildPvtbcCaller() *pvtbc.Caller
	BuildPvtbcListener() *pvtbc.Listener
	BuildDb(connString string) Db
	BuildBusPublisher(connString string) BusPublisher
	BuildBusSubscriber(connString string) BusSubscriber
}
