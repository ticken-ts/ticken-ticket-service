package bus

type TickenDevBusSubscriber struct {
}

func NewTickenDevBusSubscriber() *TickenDevBusSubscriber {
	return new(TickenDevBusSubscriber)
}

func (devbus *TickenDevBusSubscriber) Connect(_ string, _ string) error {
	return nil
}

func (devbus *TickenDevBusSubscriber) IsConnected() bool {
	return true
}

func (devbus *TickenDevBusSubscriber) Listen(_ func([]byte)) error {
	return nil
}
