package bus

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQSubscriber struct {
	conn            *amqp.Connection
	channel         *amqp.Channel
	queue           amqp.Queue
	listenQueueName string
}

func NewRabbitMQSubscriber(listenQueueName string) *RabbitMQSubscriber {
	return &RabbitMQSubscriber{
		listenQueueName: listenQueueName,
	}
}

func (subscriber *RabbitMQSubscriber) Connect(connString string, exchangeName string) error {
	conn, err := amqp.Dial(connString)
	if err != nil {
		return err
	}

	channel, err := conn.Channel()
	if err != nil {
		subscriber.Shutdown()
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
		subscriber.Shutdown()
		return err
	}

	queue, err := channel.QueueDeclare(
		subscriber.listenQueueName, // name
		false,                      // durable
		false,                      // delete when unused
		false,                      // exclusive
		false,                      // no-wait
		nil,                        // arguments
	)

	if err != nil {
		subscriber.Shutdown()
		return err
	}

	err = channel.QueueBind(
		queue.Name,   // queue name
		"",           // routing key
		exchangeName, // exchange
		false,        // no-wait
		nil,          // arguments
	)

	if err != nil {
		subscriber.Shutdown()
		return err
	}

	subscriber.conn = conn
	subscriber.channel = channel
	subscriber.queue = queue

	err = subscriber.ensureIsConnected()
	if err != nil {
		return err
	}

	return nil
}

func (subscriber *RabbitMQSubscriber) Listen(handler func([]byte)) error {
	err := subscriber.ensureIsConnected()
	if err != nil {
		return err
	}

	msgs, err := subscriber.channel.Consume(
		subscriber.queue.Name, // queue
		"",                    // consumer
		true,                  // auto-ack
		false,                 // exclusive
		false,                 // no-local
		false,                 // no-wait
		nil,                   // args
	)
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			handler(d.Body)
		}
	}()

	return nil
}

func (subscriber *RabbitMQSubscriber) IsConnected() bool {
	return subscriber.ensureIsConnected() == nil
}

func (subscriber *RabbitMQSubscriber) ensureIsConnected() error {
	if subscriber.conn == nil {
		return fmt.Errorf("RabbitMQ: connection is not stablished")
	}

	if subscriber.conn.IsClosed() {
		return fmt.Errorf("RabbitMQ: connection is closed")
	}

	return nil
}

func (subscriber *RabbitMQSubscriber) Shutdown() {
	if subscriber.conn != nil && !subscriber.conn.IsClosed() {
		err := subscriber.conn.Close()
		if err != nil {
			panic(err)
		}
	}

	if subscriber.channel != nil && !subscriber.channel.IsClosed() {
		err := subscriber.channel.Close()
		if err != nil {
			panic(err)
		}
	}
}
