package main

import (
	"context"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Exchange int8

const (
	EXCHANGE_DIRECT Exchange = iota
	EXCHANGE_TOPICS
	EXCHANGE_HEADERS
	EXCHANGE_FANOUT
)

func (e Exchange) String() string {
	return [...]string{"direct", "topic", "headers", "fanout"}[e]
}

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

func (r RabbitMQ) Publish(ctx context.Context, body []byte, exName, routingKey string, kind Exchange, durable, autoDelete, internal, noWait bool, args amqp.Table) error {
	err := r.ExchangeDeclare(exName, kind, durable, autoDelete, internal, noWait, args)
	if err != nil {
		return err
	}

	return r.channel.PublishWithContext(ctx, exName, routingKey, false, false, amqp.Publishing{
		ContentType:  "application/json",
		DeliveryMode: amqp.Persistent,
		Body:         body,
	})
}

func (r RabbitMQ) Consume(ctx context.Context, body []byte, consumer, qName, exName, routingKey string, kind Exchange, durable, autoDelete, internal, noWait bool, args amqp.Table, handler func([]byte) error) error {
	err := r.ExchangeDeclare(exName, kind, durable, autoDelete, internal, noWait, args)
	if err != nil {
		return err
	}

	q, err := r.QueueDeclare(qName, durable, autoDelete, true, noWait, args)
	if err != nil {
		return err
	}

	err = r.QueueBind(q, routingKey, exName, noWait, args)
	if err != nil {
		return err
	}

	msgs, err := r.channel.Consume(q.Name, consumer, false, false, false, false, nil)
	if err != nil {
		return err
	}

	go func() {
		for msg := range msgs {
			if err := handler(msg.Body); err != nil {
				log.Printf("Failed to handle message: %v\n", err)
			}
		}
	}()

	return nil
}

func (r RabbitMQ) ExchangeDeclare(name string, kind Exchange, durable, autoDelete, internal, noWait bool, args amqp.Table) error {
	return r.channel.ExchangeDeclare(name, kind.String(), durable, autoDelete, internal, noWait, args)
}

func (r RabbitMQ) QueueDeclare(name string, durable, deleteWhenUsed, exclusive, noWait bool, args amqp.Table) (amqp.Queue, error) {
	q, err := r.channel.QueueDeclare(name, durable, deleteWhenUsed, exclusive, noWait, args)

	if err != nil {
		return amqp.Queue{}, err
	}

	return q, nil
}

func (r RabbitMQ) QueueBind(queue amqp.Queue, routingKey, exchange string, noWait bool, args amqp.Table) error {
	return r.channel.QueueBind(queue.Name, routingKey, exchange, noWait, args)
}
