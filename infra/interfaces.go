package infra

import (
	"context"
	"github.com/gin-gonic/gin"
	pvtbc "github.com/ticken-ts/ticken-pvtbc-connector"
	"ticken-ticket-service/infra/bus"
	"ticken-ticket-service/infra/public_blockchain"
	"ticken-ticket-service/security/jwt"
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

type HSM interface {
	Store(data []byte) (string, error)
	Retrieve(key string) ([]byte, error)
}

type BusSubscriber interface {
	Connect(connString string, exchangeName string) error
	IsConnected() bool
	Listen(handler func([]byte)) error
}

type BusPublisher interface {
	Connect(connString string, exchangeName string) error
	IsConnected() bool
	Publish(ctx context.Context, msg bus.Message) error
}

type IBuilder interface {
	BuildDb(connString string) Db
	BuildHSM(encryptionKey string) HSM
	BuildEngine() *gin.Engine
	BuildJWTVerifier() jwt.Verifier
	BuildPvtbcCaller() *pvtbc.Caller
	BuildPvtbcListener() *pvtbc.Listener
	BuildBusPublisher(connString string) BusPublisher
	BuildBusSubscriber(connString string) BusSubscriber
	BuildPublicBlockchain() public_blockchain.PublicBC

	// atomic buildings
	BuildAtomicPvtbcCaller(mspID, user, peerAddr string, userCert, userPriv, tlsCert []byte) (*pvtbc.Caller, error)
}
