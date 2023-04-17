package infra

import (
	"context"
	"github.com/gin-gonic/gin"
	pubbc "github.com/ticken-ts/ticken-pubbc-connector"
	pvtbc "github.com/ticken-ts/ticken-pvtbc-connector"
	"ticken-ticket-service/infra/bus"
	"ticken-ticket-service/security/auth"
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

	BuildPubbcAdmin(privateKey string) pubbc.Admin
	BuildPubbcCaller(privateKey string) pubbc.Caller
	BuildAuthIssuer(clientSecret string) *auth.Issuer
	BuildBusPublisher(connString string) BusPublisher
	BuildBusSubscriber(connString string) BusSubscriber

	// atomic buildings
	BuildAtomicPvtbcCaller(mspID, user, peerAddr string, userCert, userPriv, tlsCert []byte) (*pvtbc.Caller, error)
}
