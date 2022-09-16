package pvtbcConnector

import (
	"encoding/json"
	"google.golang.org/grpc"
	"ticken-ticket-service/models"
)

const globalPath = "/Users/facundotorraca/Documents/ticken/ticken-dev"
const cryptoPath = globalPath + "/test-network/organizations/peerOrganizations/org1.example.com"

const (
	mspID    = "Org1MSP"
	certPath = cryptoPath + "/users/User1@org1.example.com/msp/signcerts/User1@org1.example.com-cert.pem"
	keyPath  = cryptoPath + "/users/User1@org1.example.com/msp/keystore/priv_sk"
)

const (
	TICKET_CC_ISSUE_FUNCTION = "Issue"
	TickenTickenChaincode    = "ticken-ticket"
)

type issueTicketResponse struct {
	TicketID string `json:"ticket_id"`
	EventID  string `json:"event_id"`
	Owner    string `json:"owner"`
	Status   string `json:"status"`
}

type ticketChaincodeConnector struct {
	hyperledgerFabricBaseConnector BaseConnector
}

type TicketChaincodeConnector interface {
	Connect(grpcConn *grpc.ClientConn, channel string) error
	IssueTicket(ticket *models.Ticket) error
}

func NewTicketChaincodeConnector() TicketChaincodeConnector {
	tcc := new(ticketChaincodeConnector)
	return tcc
}

func (c *ticketChaincodeConnector) Connect(grpcConn *grpc.ClientConn, channel string) error {
	c.hyperledgerFabricBaseConnector = NewBaseConnector(mspID, certPath, keyPath)
	err := c.hyperledgerFabricBaseConnector.Connect(grpcConn, channel, TickenTickenChaincode)
	if err != nil {
		return err
	}
	return nil
}

func (c *ticketChaincodeConnector) IssueTicket(ticket *models.Ticket) error {
	data, err := c.hyperledgerFabricBaseConnector.Submit(
		TICKET_CC_ISSUE_FUNCTION,
		ticket.TicketID,
		ticket.EventID,
		ticket.Section,
		ticket.Owner,
	)

	if err != nil {
		return err
	}

	response := new(issueTicketResponse)
	err = json.Unmarshal(data, &response)
	if err != nil {
		return err
	}

	return nil
}
