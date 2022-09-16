package pvtbc

import (
	"encoding/json"
	"ticken-ticket-service/models"
)

const (
	TicketCCIssueFunction = "Issue"
	TickenTicketChaincode = "ticken-ticket"
)

type issueTicketResponse struct {
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

func (ticketCC *ticketChaincodeConnector) IssueTicket(ticket *models.Ticket) error {
	data, err := ticketCC.baseCC.Submit(
		TicketCCIssueFunction,
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
