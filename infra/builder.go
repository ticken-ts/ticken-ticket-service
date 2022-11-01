package infra

import (
	"fmt"
	"github.com/gin-gonic/gin"
	pvtbc "github.com/ticken-ts/ticken-pvtbc-connector"
	"github.com/ticken-ts/ticken-pvtbc-connector/fabric/peerconnector"
	"ticken-ticket-service/config"
	"ticken-ticket-service/infra/bus"
	"ticken-ticket-service/infra/db"
	"ticken-ticket-service/log"
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
		log.TickenLogger.Info().Msg("using db: " + config.MongoDriver)
		tickenDb = db.NewMongoDb()
	default:
		err := fmt.Errorf("database driver %s not implemented", builder.tickenConfig.Database.Driver)
		log.TickenLogger.Panic().Err(err)
	}

	err := tickenDb.Connect(connString)
	if err != nil {
		log.TickenLogger.Panic().Err(err)
	}
	log.TickenLogger.Info().Msg("db connection established")

	return tickenDb
}

func (builder *Builder) BuildBusSubscriber(connString string) BusSubscriber {
	var tickenBus BusSubscriber = nil

	switch builder.tickenConfig.Bus.Driver {
	case config.RabbitMQDriver:
		log.TickenLogger.Info().Msg("using bus subscriber: " + config.RabbitMQDriver)
		tickenBus = bus.NewRabbitMQSubscriber()
	default:
		err := fmt.Errorf("bus driver %s not implemented", builder.tickenConfig.Bus.Driver)
		log.TickenLogger.Panic().Err(err)
	}

	err := tickenBus.Connect(connString, builder.tickenConfig.Bus.Exchange)
	if err != nil {
		log.TickenLogger.Panic().Err(err)
	}
	log.TickenLogger.Info().Msg("bus subscriber connection established")

	return tickenBus
}

func (builder *Builder) BuildBusPublisher(connString string) BusPublisher {
	var tickenBus BusPublisher = nil

	switch builder.tickenConfig.Bus.Driver {
	case config.RabbitMQDriver:
		log.TickenLogger.Info().Msg("using bus publisher: " + config.RabbitMQDriver)
		tickenBus = bus.NewRabbitMQPublisher()
	default:
		err := fmt.Errorf("bus driver %s not implemented", builder.tickenConfig.Bus.Driver)
		log.TickenLogger.Panic().Err(err)
	}

	err := tickenBus.Connect(connString, builder.tickenConfig.Bus.Exchange)
	if err != nil {
		log.TickenLogger.Panic().Err(err)
	}
	log.TickenLogger.Info().Msg("bus publisher connection established")

	return tickenBus
}

func (builder *Builder) BuildEngine() *gin.Engine {
	return gin.Default()
}

func (builder *Builder) BuildPvtbcCaller() *pvtbc.Caller {
	caller, err := pvtbc.NewCaller(buildPeerConnector(builder.tickenConfig.Pvtbc))
	if err != nil {
		log.TickenLogger.Panic().Err(err)
	}
	log.TickenLogger.Info().Msg("pvtbc caller created successfully")
	return caller
}

func (builder *Builder) BuildPvtbcListener() *pvtbc.Listener {
	listener, err := pvtbc.NewListener(buildPeerConnector(builder.tickenConfig.Pvtbc))
	if err != nil {
		log.TickenLogger.Panic().Err(err)
	}
	log.TickenLogger.Info().Msg("pvtbc listener created successfully")
	return listener
}

func buildPeerConnector(config config.PvtbcConfig) *peerconnector.PeerConnector {
	if pc != nil {
		return pc
	}

	pc := peerconnector.New(config.MspID, config.CertificatePath, config.PrivateKeyPath)
	log.TickenLogger.Info().Msg("pvtbc peer connector created for org " + config.MspID)

	err := pc.Connect(config.PeerEndpoint, config.GatewayPeer, config.TLSCertificatePath)
	if err != nil {
		log.TickenLogger.Panic().Err(err)
	}

	log.TickenLogger.Info().Msg("pvtbc peer connection established on " + config.PeerEndpoint)
	return pc
}
