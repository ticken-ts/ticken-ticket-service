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
		TokenID:   ticket.TokenID.String(),

		SaleAnnouncements: make([]*dto.SaleAnnouncement, 0),
	}

	if ticket.SaleAnnouncements != nil {
		for _, saleAnnouncement := range ticket.SaleAnnouncements {
			ticketDTO.SaleAnnouncements = append(
				ticketDTO.SaleAnnouncements,
				MapSaleAnnouncementToDTO(ticket, saleAnnouncement),
			)
		}
	}

	return ticketDTO
}

func MapSaleAnnouncementToDTO(ticket *models.Ticket, saleAnnoucement *models.SaleAnnouncement) *dto.SaleAnnouncement {
	return &dto.SaleAnnouncement{
		Price:    fmt.Sprintf("%.2f", saleAnnoucement.Price.Amount),
		Currency: saleAnnoucement.Price.Currency.Symbol,
		TicketID: ticket.TicketID.String(),
		EventID:  ticket.EventID.String(),
	}
}
