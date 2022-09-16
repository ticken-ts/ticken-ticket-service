package pvtbc

import (
	"fmt"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"google.golang.org/grpc"
	"io/ioutil"
	"time"
)

type baseConnector struct {
	identity *identity.X509Identity
	sign     identity.Sign
	gateway  *client.Gateway
	network  *client.Network
	contract *client.Contract
}

func FabricConnector(mspID string, certPath string, keyPath string) BaseConnector {
	return &baseConnector{
		identity: newIdentity(certPath, mspID),
		sign:     newSign(keyPath),
		gateway:  nil,
		network:  nil,
		contract: nil,
	}
}

func (bc *baseConnector) Connect(grpcConn *grpc.ClientConn, channel string, chaincode string) error {
	if bc.gateway != nil {
		return fmt.Errorf("gateway is already connected")
	}

	gateway, err := client.Connect(
		bc.identity,
		client.WithSign(bc.sign),
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

	bc.gateway = gateway
	bc.network = gateway.GetNetwork(channel)
	bc.contract = bc.network.GetContract(chaincode)
	return nil
}

func (bc *baseConnector) Query(function string, args ...string) ([]byte, error) {
	evaluateResult, err := bc.contract.EvaluateTransaction(function, args...)

	if err != nil {
		return nil, fmt.Errorf("failed to evaluate transaction: %w", err)
	}

	return evaluateResult, nil
}

func (bc *baseConnector) Submit(function string, args ...string) ([]byte, error) {
	evaluateResult, err := bc.contract.SubmitTransaction(function, args...)

	if err != nil {
		return nil, fmt.Errorf("failed to evaluate transaction: %w", err)
	}

	return evaluateResult, nil
}

func (bc *baseConnector) SubmitAsync(function string, args ...string) ([]byte, *client.Commit) {
	submitResult, commit, err := bc.contract.SubmitAsync(function, client.WithArguments(args...))
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
