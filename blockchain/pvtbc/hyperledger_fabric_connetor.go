package pvtbc

import (
	"crypto/x509"
	"fmt"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"time"
)

type hyperledgerFabricConnector struct {
	identity *identity.X509Identity
	sign     identity.Sign
	gateway  *client.Gateway
}

func NewHyperledgerFabricConnector(mspID string, certPath string, keyPath string) HyperledgerFabricConnector {
	return &hyperledgerFabricConnector{
		identity: newIdentity(certPath, mspID),
		sign:     newSign(keyPath),
		gateway:  nil,
	}
}

func (hfc *hyperledgerFabricConnector) IsConnected() bool {
	return hfc.gateway != nil
}

func (hfc *hyperledgerFabricConnector) Connect() error {
	if hfc.IsConnected() {
		return fmt.Errorf("gateway is already connected")
	}

	grpcConn, err := newGrpcConnection()
	if err != nil {
		return err
	}

	gateway, err := client.Connect(
		hfc.identity,
		client.WithSign(hfc.sign),
		client.WithClientConnection(grpcConn),

		// Default timeouts for different gRPC calls
		client.WithSubmitTimeout(5*time.Second),
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		return err
	}

	hfc.gateway = gateway
	// hfc.network = gateway.GetNetwork(channel)
	// hfc.contract = hfc.network.GetContract(chaincode)
	return nil
}

func (hfc *hyperledgerFabricConnector) GetChaincode(channelName string, chaincodeName string) (*client.Contract, error) {
	network := hfc.gateway.GetNetwork(channelName)
	if network == nil {
		return nil, fmt.Errorf("channel %s not exist", channelName)
	}

	chaincode := network.GetContract(chaincodeName)
	if chaincode == nil {
		return nil, fmt.Errorf("chaincode %s not exist", chaincodeName)
	}

	return chaincode, nil
}

// newIdentity creates a client identity for this
// Gateway connection using an X.509 certificate.
func newIdentity(certPath string, mspID string) *identity.X509Identity {
	certificatePEM, err := ioutil.ReadFile(certPath)
	if err != nil {
		panic(fmt.Errorf("failed to read certificate file: %w", err))
	}

	certificate, err := identity.CertificateFromPEM(certificatePEM)
	if err != nil {
		panic(err)
	}

	id, err := identity.NewX509Identity(mspID, certificate)
	if err != nil {
		panic(err)
	}

	return id
}

// newSign creates a function that generates a digital
// signature from a message digest using a private key.
func newSign(keyPath string) identity.Sign {
	privateKeyPEM, err := ioutil.ReadFile(keyPath)
	if err != nil {
		panic(fmt.Errorf("failed to read private key file: %w", err))
	}

	// TODO -> Undestand this
	privateKey, err := identity.PrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		panic(err)
	}

	sign, err := identity.NewPrivateKeySign(privateKey)
	if err != nil {
		panic(err)
	}

	return sign
}

func newGrpcConnection() (*grpc.ClientConn, error) {
	certificate, err := loadCertificate(tlsCertPath)
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	certPool.AddCert(certificate)
	transportCredentials := credentials.NewClientTLSFromCert(certPool, gatewayPeer)

	connection, err := grpc.Dial(peerEndpoint, grpc.WithTransportCredentials(transportCredentials))
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection: %w", err)
	}

	return connection, nil
}

func loadCertificate(filename string) (*x509.Certificate, error) {
	certificatePEM, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read certificate file: %w", err)
	}
	return identity.CertificateFromPEM(certificatePEM)
}
