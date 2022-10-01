package dto

type Ticket struct {
	TicketID string `json:"ticket_id"`
	EventID  string `json:"event_id"`
	Status   string `json:"status"`
	Section  string `json:"section"`
}

// --- blockchain
