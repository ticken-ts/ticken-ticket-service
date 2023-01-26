package mappers

import (
	"ticken-ticket-service/api/dto"
	"ticken-ticket-service/models"
)

func MapTicketToDTO(ticket *models.Ticket) *dto.Ticket {
	return &dto.Ticket{
		TicketID: ticket.TicketID.String(),
		EventID:  ticket.EventID.String(),
		Status:   ticket.Status,
		Section:  ticket.Section,
	}
}
