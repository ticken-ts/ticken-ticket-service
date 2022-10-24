package async

import (
	"encoding/json"
	"fmt"
	"ticken-ticket-service/infra"
	"ticken-ticket-service/infra/bus"
	"ticken-ticket-service/repos"
)

type Subscriber struct {
	busSubscriber  infra.BusSubscriber
	eventProcessor *EventReceiver
}

func NewSubscriber(busSubscriber infra.BusSubscriber, repoProvider repos.IProvider) (*Subscriber, error) {
	if !busSubscriber.IsConnected() {
		return nil, fmt.Errorf("bus subscriber is not connected")
	}

	subscriber := &Subscriber{
		busSubscriber:  busSubscriber,
		eventProcessor: NewEventReceiver(repoProvider.GetEventRepository()),
	}

	return subscriber, nil
}

func (processor *Subscriber) Start() error {
	err := processor.busSubscriber.Listen(processor.handler)
	if err != nil {
		return err
	}
	return nil
}

func (processor *Subscriber) handler(rawmsg []byte) {
	msg := new(bus.Message)
	err := json.Unmarshal(rawmsg, msg)
	if err != nil {
		println("error processing message")
	}

	var processingError error = nil
	switch msg.Type {
	case NewEventMessageType:
		processingError = processor.eventProcessor.NewEventHandler(msg.Data)
	default:
		processingError = fmt.Errorf("message type %s not supportaed\n", msg.Type)
	}

	if processingError != nil {
		fmt.Println(processingError)
	}
}
