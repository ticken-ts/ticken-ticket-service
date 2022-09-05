package services

import "ticken-ticket-service/models"

type Provider interface {
	GetTicketIssuer() TicketIssuer
}

type TicketIssuer interface {
	IssueTicket(eventID string, section string, owner string) (*models.Ticket, error)
}
