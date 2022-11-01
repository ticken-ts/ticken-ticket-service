package bus

import (
	"context"
	"encoding/json"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQPublisher struct {
	conn         *amqp.Connection
	channel      *amqp.Channel
	exchangeName string
}

func NewRabbitMQPublisher() *RabbitMQPublisher {
	return new(RabbitMQPublisher)
}

func (publisher *RabbitMQPublisher) Connect(connString string, exchangeName string) error {
	conn, err := amqp.Dial(connString)
	if err != nil {
		return err
	}

	channel, err := conn.Channel()
	if err != nil {
		publisher.Shutdown()
		return err
	}

	err = channel.ExchangeDeclare(
		exchangeName,        // queue
		amqp.ExchangeFanout, // consumer
		true,                // durable
		false,               // auto-delete
		false,               // internal
		false,               // no-wait
		nil,                 // args
	)
	if err != nil {
		publisher.Shutdown()
		return err
	}

	publisher.conn = conn
	publisher.channel = channel

	err = publisher.ensureIsConnected()
	if err != nil {
		return err
	}

	return nil
}

func (publisher *RabbitMQPublisher) Publish(ctx context.Context, msg Message) error {
	err := publisher.ensureIsConnected()
	if err != nil {
		return err
	}

	serializedMsg, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	publishingMessage := amqp.Publishing{
		ContentType: "text/plain",
		Body:        serializedMsg,
	}

	err = publisher.channel.PublishWithContext(
		ctx,
		publisher.exchangeName, // exchange
		"",                     // routing key
		false,                  // mandatory
		false,                  // immediate
		publishingMessage,      // msg
	)

	if err != nil {
		return err
	}

	return nil
}

func (publisher *RabbitMQPublisher) IsConnected() bool {
	return publisher.ensureIsConnected() == nil
}

func (publisher *RabbitMQPublisher) ensureIsConnected() error {
	if publisher.conn == nil {
		return fmt.Errorf("RabbitMQ: connection is not stablished")
	}

	if publisher.conn.IsClosed() {
		return fmt.Errorf("RabbitMQ: connection is closed")
	}

	return nil
}

func (publisher *RabbitMQPublisher) Shutdown() {
	if publisher.conn != nil && !publisher.conn.IsClosed() {
		err := publisher.conn.Close()
		if err != nil {
			panic(err)
		}
	}

	if publisher.channel != nil && !publisher.channel.IsClosed() {
		err := publisher.channel.Close()
		if err != nil {
			panic(err)
		}
	}
}
