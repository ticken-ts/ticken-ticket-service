package infra

import (
	"fmt"
	"github.com/gin-gonic/gin"
	pvtbc "github.com/ticken-ts/ticken-pvtbc-connector"
	"github.com/ticken-ts/ticken-pvtbc-connector/fabric/peerconnector"
	"ticken-ticket-service/config"
	"ticken-ticket-service/infra/db"
)

type Builder struct {
	tickenConfig *config.Config
}

var pc *peerconnector.PeerConnector = nil

func NewBuilder(tickenConfig *config.Config) (*Builder, error) {
	if tickenConfig == nil {
		return nil, fmt.Errorf("configuration is mandatory")
	}

	builder := new(Builder)
	builder.tickenConfig = tickenConfig

	return builder, nil
}

func (builder *Builder) BuildDb(connString string) Db {
	var tickenDb Db = nil

	switch builder.tickenConfig.Database.Driver {
	case config.MongoDriver:
		tickenDb = db.NewMongoDb()
	default:
		panic(fmt.Errorf("database driver %s not implemented", builder.tickenConfig.Database.Driver))
	}

	err := tickenDb.Connect(connString)
	if err != nil {
		panic(err)
	}

	return tickenDb
}

func (builder *Builder) BuildEngine() *gin.Engine {
	return gin.Default()
}

func (builder *Builder) BuildPvtbcCaller() *pvtbc.Caller {
	caller, err := pvtbc.NewCaller(buildPeerConnector(builder.tickenConfig.Pvtbc))
	if err != nil {
		panic(err)
	}
	return caller
}

func (builder *Builder) BuildPvtbcListener() *pvtbc.Listener {
	listener, err := pvtbc.NewListener(buildPeerConnector(builder.tickenConfig.Pvtbc))
	if err != nil {
		panic(err)
	}
	return listener
}

func buildPeerConnector(config config.PvtbcConfig) *peerconnector.PeerConnector {
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
