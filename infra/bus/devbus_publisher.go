package bus

import "context"

type TickenDevBusPublisher struct {
}

func NewTickenDevBusPublisher() *TickenDevBusPublisher {
	return new(TickenDevBusPublisher)
}

func (devbus *TickenDevBusPublisher) Connect(_ string, _ string) error {
	return nil
}

func (devbus *TickenDevBusPublisher) IsConnected() bool {
	return true
}

func (devbus *TickenDevBusPublisher) Publish(_ context.Context, _ Message) error {
	return nil
}
