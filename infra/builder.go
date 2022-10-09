package infra

// TODO
// * Handle more than one service type using config file
// * Log errors. This include passing a logger

import (
	"fmt"
	"github.com/gin-gonic/gin"
	pvtbc "github.com/ticken-ts/ticken-pvtbc-connector"
	"github.com/ticken-ts/ticken-pvtbc-connector/fabric/peerconnector"
	"ticken-ticket-service/infra/db"
	"ticken-ticket-service/utils"
)

type Builder struct {
	tickenConfig *utils.TickenConfig
}

var pc *peerconnector.PeerConnector = nil

func NewBuilder(tickenConfig *utils.TickenConfig) (*Builder, error) {
	if tickenConfig == nil {
		return nil, fmt.Errorf("configuration is mandatory")
	}

	builder := new(Builder)
	builder.tickenConfig = tickenConfig

	return builder, nil
}

func (builder *Builder) BuildDb() Db {
	switch builder.tickenConfig.Config.Database.Driver {
	case utils.MongoDriver:
		return buildMongoDb(builder.tickenConfig.Env.MongoUri)
	default:
		panic(fmt.Errorf("database driver %s not implemented", builder.tickenConfig.Config.Database.Driver))
	}
}

func (builder *Builder) BuildRouter() Router {
	return gin.Default()
}

func (builder *Builder) BuildPvtbcCaller() *pvtbc.Caller {
	caller, err := pvtbc.NewCaller(buildPeerConnector(builder.tickenConfig.Config.Pvtbc))
	if err != nil {
		panic(err)
	}
	return caller
}

func (builder *Builder) BuildPvtbcListener() *pvtbc.Listener {
	listener, err := pvtbc.NewListener(buildPeerConnector(builder.tickenConfig.Config.Pvtbc))
	if err != nil {
		panic(err)
	}
	return listener
}

func buildMongoDb(uri string) Db {
	mongoDb := db.NewMongoDb()
	err := mongoDb.Connect(uri)
	if err != nil {
		panic(err)
	}
	return mongoDb
}

func buildPeerConnector(config utils.PvtbcConfig) *peerconnector.PeerConnector {
	if pc != nil {
		return pc
	}

	pc := peerconnector.New(config.MspID, config.CertificatePath, config.PrivateKeyPath)

	err := pc.Connect(config.PeerEndpoint, config.GatewayPeer, config.TLSCertificatePath)
	if err != nil {
		panic(err)
	}

	return pc
}
