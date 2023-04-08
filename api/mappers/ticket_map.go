package mappers

import (
	"fmt"
	"ticken-ticket-service/api/dto"
	"ticken-ticket-service/models"
)

func MapTicketToDTO(ticket *models.Ticket) *dto.Ticket {
	ticketDTO := &dto.Ticket{
		TicketID: ticket.TicketID.String(),
		EventID:  ticket.EventID.String(),
		Status:   ticket.Status,
		Section:  ticket.Section,

		PubbcTxID: ticket.PubbcTxID,
		PvtbcTxID: ticket.PvtbcTxID,
		TokenID:   ticket.TokenID.Text(16),

		Resells: make([]*dto.Resells, 0),
	}

	if ticket.Resells != nil {
		for _, resell := range ticket.Resells {
			if !resell.Active {
				continue
			}
			ticketDTO.Resells = append(
				ticketDTO.Resells,
				MapResellToDTO(resell),
			)
		}
	}

	return ticketDTO
}

func MapResellToDTO(resell *models.Resell) *dto.Resells {
	return &dto.Resells{
		Price:    fmt.Sprintf("%.2f", resell.Price.Amount),
		Currency: resell.Price.Currency.Symbol,
		ResellID: resell.ResellID.String(),
	}
}
