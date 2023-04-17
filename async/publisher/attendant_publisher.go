package publisher

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"ticken-ticket-service/async/asyncmsg"
	"ticken-ticket-service/infra"
	"ticken-ticket-service/infra/bus"
	"ticken-ticket-service/models"
)

type newAttendantDTO struct {
	AttendantID   uuid.UUID `json:"attendant_id"`
	WalletAddress string    `json:"wallet_address"`
	PublicKey     string    `json:"public_key"`
}

type AttendantPublisher struct {
	busPublisher infra.BusPublisher
}

func NewEventPublisher(busPublisher infra.BusPublisher) *AttendantPublisher {
	return &AttendantPublisher{busPublisher: busPublisher}
}

func (processor *AttendantPublisher) PublishNewAttendant(attendant *models.Attendant) error {
	dto := &newAttendantDTO{
		AttendantID:   attendant.UUID,
		WalletAddress: attendant.WalletAddress,
		PublicKey:     attendant.PubKey,
	}

	serializedDTO, err := json.Marshal(dto)
	if err != nil {
		return err
	}

	err = processor.busPublisher.Publish(
		context.Background(),
		bus.Message{Type: asyncmsg.NewAttendant, Data: serializedDTO},
	)
	if err != nil {
		return err
	}

	return nil
}
