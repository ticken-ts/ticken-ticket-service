package pvtbc

import (
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"ticken-ticket-service/models"
)

type TickenConnector interface {
	EventChaincodeConnector
	TicketChaincodeConnector
	Connect(channel string) error
	IsConnected() bool
}

type HyperledgerFabricConnector interface {
	IsConnected() bool
	Connect() error
	GetChaincode(channelName string, chaincodeName string) (*client.Contract, error)
}

type ChaincodeConnector interface {
	Query(function string, args ...string) ([]byte, error)
	Submit(function string, args ...string) ([]byte, error)
	SubmitAsync(function string, args ...string) ([]byte, *client.Commit)
}

type TicketChaincodeConnector interface {
	IssueTicket(ticket *models.Ticket) error
}

type EventChaincodeConnector interface {
}
