package publisher

import (
	"fmt"
	"reflect"
	"ticken-ticket-service/async/asyncmsg"
	"ticken-ticket-service/infra"
	"ticken-ticket-service/log"
	"ticken-ticket-service/models"
)

type Publisher struct {
	busPublisher infra.BusPublisher
	*AttendantPublisher
}

func New(busPublisher infra.BusPublisher) (*Publisher, error) {
	if !busPublisher.IsConnected() {
		return nil, fmt.Errorf("bus publisher is not connected")
	}

	publisher := &Publisher{
		busPublisher:       busPublisher,
		AttendantPublisher: NewEventPublisher(busPublisher),
	}

	return publisher, nil
}

func (publisher *Publisher) PublishMessage(msgType string, content interface{}) error {
	var err error
	switch msgType {
	case asyncmsg.NewAttendant:
		ensureCasteability(content, &models.Attendant{})
		err = publisher.PublishNewAttendant(content.(*models.Attendant))
		break
	default:
		err = fmt.Errorf("publiser - message type is not supported: %s", msgType)
		break
	}
	return err
}

func ensureCasteability(from interface{}, to interface{}) {
	fromType, toType := reflect.TypeOf(from), reflect.TypeOf(to)
	if !fromType.ConvertibleTo(toType) {
		log.TickenLogger.Panic().Msg(fmt.Sprintf("%s is not casteable to %s", fromType.Name(), toType.Name()))
	}
}
