package pubsub

import (
	"context"
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Acktype int

type SimpleQueueType int

const (
	SimpleQueueDurable SimpleQueueType = iota
	SimpleQueueTransient
)

func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error {
	js, err := json.Marshal(val)
	if err != nil {
		return err
	}

	ch.PublishWithContext(context.Background(), exchange, key, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        js,
	})

	return nil
}

func DeclareAndBind(conn *amqp.Connection, exchange, queueName, key string, simpleQueueType SimpleQueueType) (*amqp.Channel, amqp.Queue, error) {
	clientChan, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
		return nil, amqp.Queue{}, err
	}

	queue, err := clientChan.QueueDeclare(
		queueName,
		simpleQueueType == SimpleQueueDurable,
		simpleQueueType != SimpleQueueDurable,
		simpleQueueType != SimpleQueueDurable,
		false, nil,
	)
	if err != nil {
		log.Fatal(err)
		return nil, amqp.Queue{}, err
	}

	err = clientChan.QueueBind(queueName, key, exchange, false, nil)
	if err != nil {
		log.Fatal(err)
		return nil, amqp.Queue{}, err
	}

	return clientChan, queue, nil
}
