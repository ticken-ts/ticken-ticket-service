package async

import (
	"encoding/json"
	"fmt"
	"ticken-ticket-service/infra"
	"ticken-ticket-service/repos"
)

type Message struct {
	Type string `json:"type"`
	Data []byte `json:"data"`
}

type Processor struct {
	busSubscriber  infra.BusSubscriber
	eventProcessor *EventProcessor
}

func NewProcessor(busSubscriber infra.BusSubscriber, repoProvider repos.IProvider) (*Processor, error) {
	if !busSubscriber.IsConnected() {
		return nil, fmt.Errorf("bus subscriber is not connected")
	}

	return &Processor{
		busSubscriber:  busSubscriber,
		eventProcessor: NewEventProcessor(repoProvider.GetEventRepository()),
	}, nil
}

func (processor *Processor) Start() error {
	err := processor.busSubscriber.Listen(processor.handler)
	if err != nil {
		return err
	}
	return nil
}

func (processor *Processor) handler(rawmsg []byte) {
	msg := new(Message)
	err := json.Unmarshal(rawmsg, msg)
	if err != nil {
		println("error processing message")
	}

	var processingError error = nil
	switch msg.Type {
	case CreateEventMessageType:
		processingError = processor.eventProcessor.CreateEvent(msg.Data)
	default:
		processingError = fmt.Errorf("message type %s not supportaed\n", msg.Type)
	}

	if processingError != nil {
		fmt.Println(processingError)
	}
}
