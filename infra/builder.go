package infra

import (
	"fmt"
	"github.com/gin-gonic/gin"
	pvtbc "github.com/ticken-ts/ticken-pvtbc-connector"
	"github.com/ticken-ts/ticken-pvtbc-connector/fabric/peerconnector"
	"ticken-ticket-service/config"
	"ticken-ticket-service/env"
	"ticken-ticket-service/infra/bus"
	"ticken-ticket-service/infra/db"
	"ticken-ticket-service/infra/hsm"
	"ticken-ticket-service/infra/public_blockchain"
	"ticken-ticket-service/log"
	"ticken-ticket-service/security/jwt"
	"ticken-ticket-service/utils"
)

type Builder struct {
	tickenConfig *config.Config
}

var pc peerconnector.PeerConnector = nil

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

func (builder *Builder) BuildHSM(encryptingKey string) HSM {
	rootPath, err := utils.GetServiceRootPath()
	if err != nil {
		panic(err)
	}

	localFileSystemHSM, err := hsm.NewLocalFileSystemHSM(encryptingKey, rootPath)
	if err != nil {
		panic(err)
	}
	return localFileSystemHSM
}

func (builder *Builder) BuildEngine() *gin.Engine {
	return gin.Default()
}

func (builder *Builder) BuildJWTVerifier() jwt.Verifier {
	var jwtVerifier jwt.Verifier

	if false {
		//if env.TickenEnv.IsDev() || env.TickenEnv.IsTest() {
		jwtPublicKey := builder.tickenConfig.Dev.JWTPublicKey
		jwtPrivateKey := builder.tickenConfig.Dev.JWTPrivateKey

		rsaPrivKey, err := utils.LoadRSA(jwtPrivateKey, jwtPublicKey)
		if err != nil {
			log.TickenLogger.Panic().Err(err)
		}

		jwtVerifier = jwt.NewOfflineVerifier(rsaPrivKey)
	} else {
		appClientID := builder.tickenConfig.Server.ClientID
		identityIssuer := builder.tickenConfig.Server.IdentityIssuer
		jwtVerifier = jwt.NewOnlineVerifier(identityIssuer, appClientID)
	}

	return jwtVerifier
}

func (builder *Builder) BuildPvtbcCaller() *pvtbc.Caller {
	caller, err := pvtbc.NewCaller(buildPeerConnector(builder.tickenConfig.Pvtbc, builder.tickenConfig.Dev))
	if err != nil {
		panic(err)
	}
	return caller
}

func (builder *Builder) BuildPvtbcListener() *pvtbc.Listener {
	listener, err := pvtbc.NewListener(buildPeerConnector(builder.tickenConfig.Pvtbc, builder.tickenConfig.Dev))
	if err != nil {
		panic(err)
	}
	return listener
}

func (builder *Builder) BuildBusPublisher(connString string) BusPublisher {
	var tickenBus BusPublisher = nil

	driverToUse := builder.tickenConfig.Bus.Driver
	if env.TickenEnv.IsDev() && !builder.tickenConfig.Dev.Mock.DisableBusMock {
		driverToUse = config.DevBusDriver
	}

	switch driverToUse {
	case config.DevBusDriver:
		log.TickenLogger.Info().Msg("using bus publisher: " + config.DevBusDriver)
		tickenBus = bus.NewTickenDevBusPublisher()
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

func (builder *Builder) BuildBusSubscriber(connString string) BusSubscriber {
	var tickenBus BusSubscriber = nil

	driverToUse := builder.tickenConfig.Bus.Driver
	if env.TickenEnv.IsDev() && !builder.tickenConfig.Dev.Mock.DisableBusMock {
		driverToUse = config.DevBusDriver
	}

	switch driverToUse {
	case config.DevBusDriver:
		log.TickenLogger.Info().Msg("using bus publisher: " + config.DevBusDriver)
		tickenBus = bus.NewTickenDevBusSubscriber()
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

func (builder *Builder) BuildPublicBlockchain() public_blockchain.PublicBC {
	pbConfig := builder.tickenConfig.PublicBlockchain

	//If is dev, generate a dev blockchain connector
	if env.TickenEnv.IsDev() {
		log.TickenLogger.Info().Msg("dev public blockchain connection established")
		return public_blockchain.NewDevPublicBlockchain()
	}

	pbbc := public_blockchain.NewPublicBlockchain(pbConfig.ChainURL, pbConfig.ChainID, pbConfig.AddressPK)
	//err := pbbc.Connect()
	//if err != nil {
	//	log.TickenLogger.Panic().Err(err)
	//}
	log.TickenLogger.Info().Msg("public blockchain connection established")
	return pbbc
}

func (builder *Builder) BuildAtomicPvtbcCaller(mspID, user, peerAddr string, userCert, userPriv, tlsCert []byte) (*pvtbc.Caller, error) {
	var pc peerconnector.PeerConnector
	if env.TickenEnv.IsDev() && !builder.tickenConfig.Dev.Mock.DisablePVTBCMock {
		pc = peerconnector.NewDev(mspID, user)
	} else {
		pc = peerconnector.NewWithRawCredentials(mspID, userCert, userPriv)
	}

	err := pc.ConnectWithRawTlsCert(peerAddr, peerAddr, tlsCert)
	if err != nil {
		return nil, err
	}

	caller, err := pvtbc.NewCaller(buildPeerConnector(builder.tickenConfig.Pvtbc, builder.tickenConfig.Dev))
	if err != nil {
		return nil, err
	}

	return caller, nil
}

func buildPeerConnector(config config.PvtbcConfig, devConfig config.DevConfig) peerconnector.PeerConnector {
	if pc != nil {
		return pc
	}

	var pc peerconnector.PeerConnector
	if env.TickenEnv.IsDev() && !devConfig.Mock.DisablePVTBCMock {
		pc = peerconnector.NewDev(config.MspID, "admin")
	} else {
		pc = peerconnector.New(config.MspID, config.CertificatePath, config.PrivateKeyPath)
	}

	err := pc.Connect(config.PeerEndpoint, config.GatewayPeer, config.TLSCertificatePath)
	if err != nil {
		panic(err)
	}

	return pc
}
