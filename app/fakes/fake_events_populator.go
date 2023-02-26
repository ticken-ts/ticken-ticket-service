package fakes

import (
	"github.com/google/uuid"
	ticken_pubbc_connector "github.com/ticken-ts/ticken-pubbc-connector"
	"ticken-ticket-service/env"
	"ticken-ticket-service/models"
	"ticken-ticket-service/repos"
)

type FakeEventsPopulator struct {
	EventRepo repos.EventRepository
	Pubbc     ticken_pubbc_connector.Admin
}

func (populator *FakeEventsPopulator) Populate() error {
	if !env.TickenEnv.IsDev() {
		return nil
	}

	eventID := uuid.MustParse("8709adbb-0504-4707-9cb2-867126c8172f")
	event := populator.EventRepo.FindEvent(eventID)
	if event != nil {
		return nil
	}
	
	fakeEvent := &models.Event{
		EventID:      eventID,
		OrganizerID:  uuid.New(),
		PvtBCChannel: "ticken-event-name",
		PubBCAddress: "0xfafafa",
	}

	addr, _ := populator.Pubbc.DeployEventContract()

	fakeEvent.PubBCAddress = addr

	return populator.EventRepo.AddEvent(fakeEvent)
}
