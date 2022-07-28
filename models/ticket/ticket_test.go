package ticket_test

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"ticken-ticket-service/models/ticket"
)

func Test_Ticket_AssignTo_CanAssignOwnerToNewTicket(t *testing.T) {
	eventID := uuid.NewString()
	section := "V.I.P"

	newTicket := ticket.New(eventID, section)
	assert.Empty(t, newTicket.Owner)

	owner := uuid.NewString()
	err := newTicket.AssignTo(owner)
	assert.NoError(t, err)

	assert.Equal(t, newTicket.Owner, owner)
}

func Test_Ticket_AssignTo_AssignOwnerToTicketThatHasOwnerReturnsError(t *testing.T) {
	eventID := uuid.NewString()
	section := "V.I.P"

	newTicket := ticket.New(eventID, section)

	owner := uuid.NewString()
	err := newTicket.AssignTo(owner)
	assert.NoError(t, err)

	otherOwner := uuid.NewString()
	err = newTicket.AssignTo(otherOwner)
	assert.Error(t, err)

	assert.Equal(t, newTicket.Owner, owner)
}

func Test_Ticket_HasOwner_ReturnsCorrectlyIfTicketHasOwner(t *testing.T) {
	eventID := uuid.NewString()
	section := "V.I.P"

	newTicket := ticket.New(eventID, section)
	assert.False(t, newTicket.HasOwner())

	owner := uuid.NewString()
	_ = newTicket.AssignTo(owner)
	assert.True(t, newTicket.HasOwner())
}
