package config

import (
	"fmt"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"google.golang.org/grpc"
	"io/ioutil"
	"time"
)

type HyperledgerFabricServiceInterface interface {
	Connect(grpcConn *grpc.ClientConn, channel string, chaincode string)
	Query(function string, args ...string) []byte
	Submit(function string, args ...string) []byte
	SubmitAsync(function string, args ...string) ([]byte, *client.Commit)
}

type hyperledgerFabricService struct {
	identity *identity.X509Identity
	sign     identity.Sign
	gateway  *client.Gateway
	network  *client.Network
	contract *client.Contract
}

func New(mspID string, certPath string, keyPath string) HyperledgerFabricServiceInterface {
	return &hyperledgerFabricService{
		identity: newIdentity(certPath, mspID),
		sign:     newSign(keyPath),
		gateway:  nil,
		network:  nil,
		contract: nil,
	}
}

func (hfs *hyperledgerFabricService) Connect(grpcConn *grpc.ClientConn, channel string, chaincode string) {
	if hfs.gateway != nil {
		return
	}

	gateway, err := client.Connect(
		hfs.identity,
		client.WithSign(hfs.sign),
		client.WithClientConnection(grpcConn),

		// Default timeouts for different gRPC calls
		client.WithSubmitTimeout(5*time.Second),
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)

	if err != nil {
		panic(err)
	}

	hfs.gateway = gateway
	hfs.network = gateway.GetNetwork(channel)
	hfs.contract = hfs.network.GetContract(chaincode)
}

func (hfs *hyperledgerFabricService) Query(function string, args ...string) []byte {
	evaluateResult, err := hfs.contract.EvaluateTransaction(function, args...)

	if err != nil {
		panic(fmt.Errorf("failed to evaluate transaction: %w", err))
	}

	return evaluateResult
}

func (hfs *hyperledgerFabricService) Submit(function string, args ...string) []byte {
	evaluateResult, err := hfs.contract.SubmitTransaction(function, args...)

	if err != nil {
		panic(fmt.Errorf("failed to evaluate transaction: %w", err))
	}

	return evaluateResult
}

func (hfs *hyperledgerFabricService) SubmitAsync(function string, args ...string) ([]byte, *client.Commit) {
	submitResult, commit, err := hfs.contract.SubmitAsync(function, client.WithArguments(args...))
	if err != nil {
		panic(fmt.Errorf("failed to submit transaction asynchronously: %w", err))
	}

	return submitResult, commit
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
	privateKey, err := ioutil.ReadFile(keyPath)
	if err != nil {
		panic(fmt.Errorf("failed to read private key file: %w", err))
	}

	sign, err := identity.NewPrivateKeySign(privateKey)
	if err != nil {
		panic(err)
	}

	return sign
}
