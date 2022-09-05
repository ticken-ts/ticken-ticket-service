package tickenPVTBCConnector

import (
	"crypto/x509"
	"fmt"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"ticken-ticket-service/blockchain/tickenPVTBCConnector/hyperledgerFabricConnectors"
	"ticken-ticket-service/models"
)

const globalPath = "/Users/facundotorraca/Documents/ticken/papers-and-books/repos/fabric-samples"
const cryptoPath = globalPath + "/test-network/organizations/peerOrganizations/org1.example.com"

const (
	tlsCertPath  = cryptoPath + "/peers/peer0.org1.example.com/tls/ca.crt"
	peerEndpoint = "localhost:7051"
	gatewayPeer  = "peer0.org1.example.com"
)

type TickenPVTBConnector interface {
	Connect(channel string) error
	IsConnected() bool
	IssueTicket(ticket *models.Ticket) error
}

type tickenPVTBCConnector struct {
	channel                  string
	grpcConn                 *grpc.ClientConn
	eventChaincodeConnector  hyperledgerFabricConnectors.EventChaincodeConnector
	ticketChaincodeConnector hyperledgerFabricConnectors.TicketChaincodeConnector
}

func New() TickenPVTBConnector {
	tickenPVTBCService := new(tickenPVTBCConnector)

	tickenPVTBCService.eventChaincodeConnector = hyperledgerFabricConnectors.NewEventChaincodeConnector()
	tickenPVTBCService.ticketChaincodeConnector = hyperledgerFabricConnectors.NewTicketChaincodeConnector()

	return tickenPVTBCService
}

func (s *tickenPVTBCConnector) Connect(channel string) error {
	if s.IsConnected() {
		if s.channel == channel {
			return fmt.Errorf("already connected to channel %s", channel)
		}
		s.channel = channel
	} else {
		newGrpcConn, err := newGrpcConnection()
		if err != nil {
			return err
		}
		s.grpcConn = newGrpcConn
	}

	err := s.eventChaincodeConnector.Connect(s.grpcConn, s.channel)
	if err != nil {
		return err
	}

	err = s.ticketChaincodeConnector.Connect(s.grpcConn, s.channel)
	if err != nil {
		return err
	}

	return nil
}

func (s *tickenPVTBCConnector) IsConnected() bool {
	return s.grpcConn != nil
}

func (s *tickenPVTBCConnector) IssueTicket(ticket *models.Ticket) error {
	if !s.IsConnected() {
		return fmt.Errorf("service is not connected")
	}

	err := s.ticketChaincodeConnector.IssueTicket(ticket)
	if err != nil {
		return err
	}
	return nil
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
