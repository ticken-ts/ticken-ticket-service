package ticket

import (
	"fmt"
	"github.com/google/uuid"
)

type Ticket struct {
	TicketID string `json:"id"`
	Owner    string `json:"owner"`
	EventID  string `json:"event_id"`
	Section  string `json:"section"`
}

func New(eventID string, section string) *Ticket {
	return &Ticket{
		TicketID: uuid.New().String(),
		EventID:  eventID,
		Section:  section,
	}
}

func (t *Ticket) AssignTo(owner string) error {
	if t.HasOwner() {
		return fmt.Errorf("ticket already has owner")
	}

	t.Owner = owner
	return nil
}

func (t *Ticket) HasOwner() bool {
	return len(t.Owner) != 0
}
