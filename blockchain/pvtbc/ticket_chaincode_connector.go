package pvtbc

import (
	"encoding/json"
	"fmt"
)

const (
	TicketCCIssueFunction = "Issue"
	TicketCCSignFunction  = "Sign"
	TickenTicketChaincode = "ticken-ticket"
)

type TicketChaincode struct {
	TicketID string `json:"ticket_id"`
	EventID  string `json:"event_id"`
	Owner    string `json:"owner"`
	Status   string `json:"status"`
}

type ticketChaincodeConnector struct {
	baseCC ChaincodeConnector
}

func NewTicketChaincodeConnector(hfc HyperledgerFabricConnector, channelName string) (TicketChaincodeConnector, error) {
	cc, err := NewChaincodeConnector(hfc, channelName, TickenTicketChaincode)
	if err != nil {
		return nil, err
	}

	ticketCC := new(ticketChaincodeConnector)
	ticketCC.baseCC = cc

	return ticketCC, nil
}

func (ticketCC *ticketChaincodeConnector) IssueTicket(ticketID string, eventID string, section string, owner string) (*TicketChaincode, error) {
	data, err := ticketCC.baseCC.Submit(TicketCCIssueFunction, ticketID, eventID, section, owner)
	if err != nil {
		return nil, err
	}

	response := new(TicketChaincode)
	err = json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (ticketCC *ticketChaincodeConnector) SignTicket(ticketID string, eventID string, signer string, signature []byte) (*TicketChaincode, error) {
	hexSignature := fmt.Sprintf("%x", signature)

	data, err := ticketCC.baseCC.Submit(TicketCCSignFunction, ticketID, eventID, signer, hexSignature)
	if err != nil {
		return nil, err
	}

	response := new(TicketChaincode)
	err = json.Unmarshal(data, &response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
