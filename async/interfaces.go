package async

type IAsyncPublisher interface {
	PublishMessage(msgType string, content interface{}) error
}

type IAsyncSubscriber interface {
	ListenMessages() error
}
