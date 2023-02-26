package dto

type Ticket struct {
	TicketID string `json:"ticket_id"`
	EventID  string `json:"event_id"`
	Status   string `json:"status"`
	Section  string `json:"section"`

	PubbcTxID string `json:"pubbc_tx_id"`
	PvtbcTxID string `json:"pvtbc_tx_id"`
	TokenID   string `json:"token_id"`
}
